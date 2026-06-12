package dlq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"dlq-service/pkg/kafka"

	"gorm.io/gorm"
)

// DeadLetterMessage represents a failed message
type DeadLetterMessage struct {
	ID               uint      `gorm:"primaryKey"`
	OriginalTopic    string    `gorm:"type:varchar(255);not null"`
	OriginalPartition int      `gorm:"not null"`
	OriginalOffset   int64     `gorm:"not null"`
	OriginalMessage  string    `gorm:"type:text;not null"`
	ErrorMessage     string    `gorm:"type:text;not null"`
	Service          string    `gorm:"type:varchar(100);not null"`
	RetryCount       int       `gorm:"default:0"`
	Status           string    `gorm:"type:varchar(50);default:'pending'"` // pending, retrying, failed, resolved
	FailedAt         time.Time `gorm:"not null"`
	LastRetryAt      *time.Time
	ResolvedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// DLQHandler handles dead letter queue messages
type DLQHandler struct {
	db            *gorm.DB
	consumer      *kafka.Consumer
	retryProducer *kafka.Producer
}

// NewDLQHandler creates a new DLQ handler
func NewDLQHandler(brokers []string, db *gorm.DB) *DLQHandler {
	// Auto-migrate DLQ table
	db.AutoMigrate(&DeadLetterMessage{})

	config := kafka.ConsumerConfig{
		Brokers: brokers,
		Topic:   "dlq.events",
		GroupID: "dlq-handler",
		Service: "dlq-service",
		Handler: nil,
	}

	handler := &DLQHandler{
		db:            db,
		retryProducer: kafka.NewProducer(brokers, "dlq-service"),
	}

	config.Handler = handler.handleDLQMessage
	handler.consumer = kafka.NewConsumer(config)

	return handler
}

// Start starts the DLQ handler
func (h *DLQHandler) Start(ctx context.Context) error {
	log.Println("Starting DLQ handler...")
	
	// Start automatic retry process
	go h.autoRetryLoop(ctx)
	
	// Start consuming DLQ messages
	return h.consumer.Start(ctx)
}

// handleDLQMessage processes messages from the DLQ
func (h *DLQHandler) handleDLQMessage(ctx context.Context, event map[string]interface{}) error {
	log.Printf("Received DLQ message: %v", event)

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid DLQ event format")
	}

	// Parse DLQ message
	dlqMsg := DeadLetterMessage{
		OriginalTopic:     getString(data, "original_topic"),
		OriginalPartition: getInt(data, "original_partition"),
		OriginalOffset:    getInt64(data, "original_offset"),
		OriginalMessage:   getString(data, "original_message"),
		ErrorMessage:      getString(data, "error_message"),
		Service:           getString(data, "service"),
		RetryCount:        0,
		Status:            "pending",
		FailedAt:          time.Now(),
	}

	// Store in database
	if err := h.db.Create(&dlqMsg).Error; err != nil {
		return fmt.Errorf("failed to store DLQ message: %w", err)
	}

	log.Printf("DLQ message stored with ID: %d", dlqMsg.ID)
	return nil
}

// autoRetryLoop automatically retries failed messages
func (h *DLQHandler) autoRetryLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.retryPendingMessages(ctx)
		}
	}
}

// retryPendingMessages retries messages that are eligible for retry
func (h *DLQHandler) retryPendingMessages(ctx context.Context) {
	var messages []DeadLetterMessage

	// Find messages eligible for retry (max 3 retries, not retried in last 10 minutes)
	err := h.db.Where("status = ? AND retry_count < ?", "pending", 3).
		Where("last_retry_at IS NULL OR last_retry_at < ?", time.Now().Add(-10*time.Minute)).
		Limit(100).
		Find(&messages).Error

	if err != nil {
		log.Printf("Error fetching pending DLQ messages: %v", err)
		return
	}

	log.Printf("Found %d messages to retry", len(messages))

	for _, msg := range messages {
		if err := h.retryMessage(ctx, &msg); err != nil {
			log.Printf("Failed to retry message %d: %v", msg.ID, err)
		}
	}
}

// retryMessage retries a single message
func (h *DLQHandler) retryMessage(ctx context.Context, msg *DeadLetterMessage) error {
	log.Printf("Retrying message %d (attempt %d)", msg.ID, msg.RetryCount+1)

	// Parse original message
	var originalEvent map[string]interface{}
	if err := json.Unmarshal([]byte(msg.OriginalMessage), &originalEvent); err != nil {
		return fmt.Errorf("failed to parse original message: %w", err)
	}

	// Republish to original topic
	if err := h.retryProducer.PublishEvent(ctx, msg.OriginalTopic, "retry", originalEvent); err != nil {
		// Update retry count
		now := time.Now()
		msg.RetryCount++
		msg.LastRetryAt = &now

		if msg.RetryCount >= 3 {
			msg.Status = "failed"
		}

		h.db.Save(msg)
		return fmt.Errorf("failed to republish message: %w", err)
	}

	// Mark as resolved
	now := time.Now()
	msg.Status = "resolved"
	msg.ResolvedAt = &now
	h.db.Save(msg)

	log.Printf("Message %d successfully retried", msg.ID)
	return nil
}

// GetDLQStats returns statistics about DLQ messages
func (h *DLQHandler) GetDLQStats() map[string]interface{} {
	var stats struct {
		Pending  int64
		Retrying int64
		Failed   int64
		Resolved int64
		Total    int64
	}

	h.db.Model(&DeadLetterMessage{}).Where("status = ?", "pending").Count(&stats.Pending)
	h.db.Model(&DeadLetterMessage{}).Where("status = ?", "retrying").Count(&stats.Retrying)
	h.db.Model(&DeadLetterMessage{}).Where("status = ?", "failed").Count(&stats.Failed)
	h.db.Model(&DeadLetterMessage{}).Where("status = ?", "resolved").Count(&stats.Resolved)
	h.db.Model(&DeadLetterMessage{}).Count(&stats.Total)

	return map[string]interface{}{
		"pending":  stats.Pending,
		"retrying": stats.Retrying,
		"failed":   stats.Failed,
		"resolved": stats.Resolved,
		"total":    stats.Total,
	}
}

// ManualRetry manually retries a specific message by ID
func (h *DLQHandler) ManualRetry(ctx context.Context, messageID uint) error {
	var msg DeadLetterMessage
	if err := h.db.First(&msg, messageID).Error; err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	return h.retryMessage(ctx, &msg)
}

// MarkAsResolved manually marks a message as resolved
func (h *DLQHandler) MarkAsResolved(messageID uint, reason string) error {
	now := time.Now()
	return h.db.Model(&DeadLetterMessage{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"status":      "resolved",
			"resolved_at": &now,
			"error_message": reason,
		}).Error
}

// Helper functions
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(float64); ok {
		return int64(val)
	}
	return 0
}
