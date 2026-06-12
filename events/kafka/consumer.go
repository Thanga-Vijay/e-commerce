package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventHandler func(ctx context.Context, event map[string]interface{}) error

type Consumer struct {
	reader  *kafka.Reader
	service string
	handler EventHandler
	dlqProducer *Producer
}

type ConsumerConfig struct {
	Brokers       []string
	Topic         string
	GroupID       string
	Service       string
	Handler       EventHandler
	MinBytes      int
	MaxBytes      int
	CommitInterval time.Duration
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config ConsumerConfig) *Consumer {
	if config.MinBytes == 0 {
		config.MinBytes = 10e3 // 10KB
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = 10e6 // 10MB
	}
	if config.CommitInterval == 0 {
		config.CommitInterval = time.Second
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topic,
		GroupID:        config.GroupID,
		MinBytes:       config.MinBytes,
		MaxBytes:       config.MaxBytes,
		CommitInterval: config.CommitInterval,
		StartOffset:    kafka.LastOffset,
		MaxAttempts:    3,
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("ERROR: "+msg+"\n", args...)
		}),
	})

	// Create DLQ producer for failed messages
	dlqProducer := NewProducer(config.Brokers, config.Service)

	return &Consumer{
		reader:      reader,
		service:     config.Service,
		handler:     config.Handler,
		dlqProducer: dlqProducer,
	}
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("Starting consumer for topic: %s, service: %s", c.reader.Config().Topic, c.service)

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutting down...")
			return c.Close()
		default:
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				log.Printf("Error fetching message: %v", err)
				continue
			}

			if err := c.processMessage(ctx, message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Send to DLQ
				c.sendToDLQ(ctx, message, err)
			} else {
				// Commit message offset
				if err := c.reader.CommitMessages(ctx, message); err != nil {
					log.Printf("Error committing message: %v", err)
				}
			}
		}
	}
}

// processMessage processes a single Kafka message
func (c *Consumer) processMessage(ctx context.Context, message kafka.Message) error {
	log.Printf("Processing message from topic: %s, partition: %d, offset: %d",
		message.Topic, message.Partition, message.Offset)

	// Deserialize message
	var event map[string]interface{}
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Add message metadata to event
	event["_kafka_topic"] = message.Topic
	event["_kafka_partition"] = message.Partition
	event["_kafka_offset"] = message.Offset
	event["_kafka_timestamp"] = message.Time

	// Call handler
	if err := c.handler(ctx, event); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}

	log.Printf("Successfully processed message from offset: %d", message.Offset)
	return nil
}

// sendToDLQ sends failed messages to Dead Letter Queue
func (c *Consumer) sendToDLQ(ctx context.Context, message kafka.Message, processingError error) {
	dlqEvent := map[string]interface{}{
		"original_topic":     message.Topic,
		"original_partition": message.Partition,
		"original_offset":    message.Offset,
		"original_message":   string(message.Value),
		"error_message":      processingError.Error(),
		"service":            c.service,
		"failed_at":          time.Now().UTC(),
	}

	if err := c.dlqProducer.PublishEvent(ctx, "dlq.events", "message.failed", dlqEvent); err != nil {
		log.Printf("Failed to send message to DLQ: %v", err)
	} else {
		log.Printf("Message sent to DLQ from offset: %d", message.Offset)
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	if err := c.dlqProducer.Close(); err != nil {
		log.Printf("Error closing DLQ producer: %v", err)
	}
	return c.reader.Close()
}

// Stats returns consumer statistics
func (c *Consumer) Stats() kafka.ReaderStats {
	return c.reader.Stats()
}
