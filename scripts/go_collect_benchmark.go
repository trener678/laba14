package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"laba14-health-pipeline/internal/health"
)

func main() {
	duration := flag.Duration("duration", 30*time.Second, "benchmark duration")
	rate := flag.Int("rate", 0, "records per second; 0 means max throughput")
	flag.Parse()

	sim := health.NewSimulator(42)
	sources := health.DemoSources()
	deadline := time.Now().Add(*duration)
	total := 0
	start := time.Now()

	for time.Now().Before(deadline) {
		limit := *rate
		if limit <= 0 {
			limit = 10000
		}
		tickStart := time.Now()
		for i := 0; i < limit; i++ {
			source := sources[i%len(sources)]
			reading := sim.Reading(source, time.Now())
			if err := health.ValidateReading(reading); err == nil {
				total++
			}
		}
		if *rate > 0 {
			if remaining := time.Second - time.Since(tickStart); remaining > 0 {
				time.Sleep(remaining)
			}
		} else {
			runtime.Gosched()
		}
	}

	elapsed := time.Since(start).Seconds()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("records_per_second=%.2f\n", float64(total)/elapsed)
	fmt.Printf("rss_mb=%.2f\n", float64(mem.Sys)/1024/1024)
}
