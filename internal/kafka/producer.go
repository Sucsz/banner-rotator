package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer — интерфейс Kafka‑продюсера.
type Producer interface {
	Send(ctx context.Context, event BannerEvent) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

// NewProducer создаёт Kafka‑продюсер с заданными брокерами и топиком.
func NewProducer(brokers []string, topic string) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

// Send сериализует BannerEvent и отправляет его в Kafka.
func (p *producer) Send(ctx context.Context, event BannerEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("kafka.Producer.Send: marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event.Type),
		Value: data,
		Time:  time.Now(),
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka.Producer.Send: write message: %w", err)
	}
	return nil
}

// Close закрывает внутренний kafka.Writer.
func (p *producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("kafka.Producer.Close: %w", err)
	}
	return nil
}
