package kafkastream

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"laba14-health-pipeline/internal/health"
)

type Producer struct {
	mu    sync.Mutex
	file  *os.File
	topic string
}

func NewProducer(brokers []string, topic string) *Producer {
	_ = brokers
	file, _ := os.OpenFile("health-aggregates.kafka.jsonl", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	return &Producer{file: file, topic: topic}
}

func (p *Producer) Publish(ctx context.Context, aggregate health.Aggregate) error {
	_ = ctx
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.file == nil {
		return nil
	}
	payload, err := json.Marshal(aggregate)
	if err != nil {
		return err
	}
	_, err = p.file.Write(append(payload, '\n'))
	time.Sleep(1 * time.Millisecond)
	return err
}

func (p *Producer) Close() error {
	if p.file == nil {
		return nil
	}
	return p.file.Close()
}
