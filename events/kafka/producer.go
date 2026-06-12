package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer  *kafka.Writer
	service string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, service string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Compression:  kafka.Snappy,
		MaxAttempts:  3,
		RequiredAcks: kafka.RequireOne,
		Async:        true, // Non-blocking writes
		ErrorLogger:  kafka.LoggerFunc(func(msg string, args ...interface{}) {
			fmt.Printf("ERROR: "+msg+"\n", args...)
		}),
	}

	return &Producer{
		writer:  writer,
		service: service,
	}
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, topic string, eventType string, payload interface{}) error {
	// Create event envelope
	envelope := map[string]interface{}{
		"id":        uuid.New().String(),
		"type":      eventType,
		"source":    p.service,
		"timestamp": time.Now().UTC(),
		"version":   "1.0",
		"data":      payload,
	}

	// Serialize to JSON
	value, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	message := kafka.Message{
		Topic: topic,
		Key:   []byte(fmt.Sprintf("%s-%s", eventType, time.Now().Format("20060102"))),
		Value: value,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(eventType)},
			{Key: "source", Value: []byte(p.service)},
			{Key: "version", Value: []byte("1.0")},
		},
	}

	// Write message to Kafka
	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

// PublishBatchEvents publishes multiple events in a batch
func (p *Producer) PublishBatchEvents(ctx context.Context, topic string, events []interface{}) error {
	messages := make([]kafka.Message, 0, len(events))

	for _, event := range events {
		envelope := map[string]interface{}{
			"id":        uuid.New().String(),
			"type":      topic,
			"source":    p.service,
			"timestamp": time.Now().UTC(),
			"version":   "1.0",
			"data":      event,
		}

		value, err := json.Marshal(envelope)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		message := kafka.Message{
			Topic: topic,
			Value: value,
			Time:  time.Now(),
			Headers: []kafka.Header{
				{Key: "source", Value: []byte(p.service)},
			},
		}

		messages = append(messages, message)
	}

	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		return fmt.Errorf("failed to write batch messages: %w", err)
	}

	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Stats returns producer statistics
func (p *Producer) Stats() kafka.WriterStats {
	return p.writer.Stats()
}
