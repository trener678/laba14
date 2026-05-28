package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"laba14-health-pipeline/internal/arrowipc"
	"laba14-health-pipeline/internal/health"
	"laba14-health-pipeline/internal/kafkastream"
	"laba14-health-pipeline/internal/sharding"
	"laba14-health-pipeline/internal/validator"
	"laba14-health-pipeline/internal/window"
)

func main() {
	var (
		collectorID = flag.String("collector-id", env("COLLECTOR_ID", "collector-local"), "collector instance id")
		httpAddr    = flag.String("http", env("HTTP_ADDR", ":8080"), "Arrow HTTP server address")
		windowSize  = flag.Duration("window", 10*time.Second, "tumbling window duration")
		interval    = flag.Duration("interval", time.Second, "sensor polling interval")
		etcdCSV     = flag.String("etcd", env("ETCD_ENDPOINTS", ""), "comma-separated etcd endpoints")
		kafkaCSV    = flag.String("kafka", env("KAFKA_BROKERS", ""), "comma-separated Kafka brokers")
		kafkaTopic  = flag.String("topic", env("KAFKA_TOPIC", "health-aggregates"), "Kafka topic")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collectors := []string{*collectorID}
	if *etcdCSV != "" {
		coordinator, err := sharding.NewEtcdCoordinator(strings.Split(*etcdCSV, ","), "/laba14/collectors", 10)
		if err != nil {
			log.Fatalf("connect etcd: %v", err)
		}
		defer coordinator.Close()
		if err := coordinator.Register(ctx, *collectorID); err != nil {
			log.Fatalf("register collector: %v", err)
		}
		if active, err := coordinator.Collectors(ctx); err == nil {
			collectors = active
		}
	}

	assignments := sharding.AssignSources(collectors, health.DemoSources())
	mySources := assignments[*collectorID].Sources
	if len(mySources) == 0 {
		log.Printf("collector %s has no sources in current assignment", *collectorID)
	}

	store := arrowipc.NewStore()
	go func() {
		log.Printf("Arrow IPC server listening on %s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, store.Handler()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	var producer *kafkastream.Producer
	if *kafkaCSV != "" {
		producer = kafkastream.NewProducer(strings.Split(*kafkaCSV, ","), *kafkaTopic)
		defer producer.Close()
	}

	sim := health.NewSimulator(time.Now().UnixNano())
	aggregator := window.NewTumblingAggregator(*windowSize)
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for at := range ticker.C {
		for _, source := range mySources {
			reading := sim.Reading(source, at)
			if err := validator.Validate(reading); err != nil {
				log.Printf("skip invalid reading: %v", err)
				continue
			}
			for _, aggregate := range aggregator.Add(reading) {
				store.Append(aggregate)
				if producer != nil {
					if err := producer.Publish(ctx, aggregate); err != nil {
						log.Printf("publish aggregate to kafka: %v", err)
					}
				}
			}
		}
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
