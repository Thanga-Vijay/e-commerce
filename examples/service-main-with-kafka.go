package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example-service/pkg/kafka"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Kafka configuration
	kafkaConfig := kafka.LoadConfig()

	var producer *kafka.Producer
	var consumers []*kafka.Consumer

	if kafkaConfig.Enabled {
		// Initialize Kafka producer
		producer = kafka.NewProducer(kafkaConfig.Brokers, "example-service")
		log.Println("Kafka producer initialized")

		// Initialize Kafka consumers
		// Example: consuming order events
		orderConsumerConfig := kafka.ConsumerConfig{
			Brokers: kafkaConfig.Brokers,
			Topic:   "order.created",
			GroupID: "example-service-orders",
			Service: "example-service",
			Handler: handleOrderEvent,
		}
		orderConsumer := kafka.NewConsumer(orderConsumerConfig)
		consumers = append(consumers, orderConsumer)

		// Start consumers in goroutines
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, consumer := range consumers {
			go func(c *kafka.Consumer) {
				if err := c.Start(ctx); err != nil {
					log.Printf("Consumer error: %v", err)
				}
			}(consumer)
		}
	}

	// Initialize Gin router
	r := gin.Default()

	// Pass Kafka producer to handlers
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Example route that publishes an event
	r.POST("/example", func(c *gin.Context) {
		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Publish event
		if producer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := producer.PublishEvent(ctx, "example.topic", "example.created", input); err != nil {
				log.Printf("Failed to publish event: %v", err)
				c.JSON(500, gin.H{"error": "Failed to publish event"})
				return
			}
		}

		c.JSON(200, gin.H{"status": "event published"})
	})

	// Start server
	srv := &gin.Engine{}
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Close Kafka connections
	if producer != nil {
		producer.Close()
	}

	log.Println("Server exited")
}

// handleOrderEvent is an example event handler
func handleOrderEvent(ctx context.Context, event map[string]interface{}) error {
	log.Printf("Processing order event: %v", event)
	// Add your business logic here
	return nil
}
