package kafkastream

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"

	"laba14-health-pipeline/internal/health"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireOne,
			AllowAutoTopicCreation: true,
			BatchSize:              32,
			BatchTimeout:           100 * time.Millisecond,
			WriteTimeout:           10 * time.Second,
			ReadTimeout:            10 * time.Second,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, aggregate health.Aggregate) error {
	payload, err := json.Marshal(aggregate)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Time:  time.Now().UTC(),
		Key:   []byte(aggregate.Metric),
		Value: payload,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
