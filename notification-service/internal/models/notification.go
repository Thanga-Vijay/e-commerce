package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification types
const (
	NotificationTypeEmail = "email"
	NotificationTypeSMS   = "sms"
	NotificationTypePush  = "push"
)

// Notification statuses
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// Notification templates
const (
	TemplateWelcome              = "welcome"
	TemplateEmailVerification    = "email_verification"
	TemplatePasswordReset        = "password_reset"
	TemplateOrderConfirmation    = "order_confirmation"
	TemplatePaymentReceipt       = "payment_receipt"
	TemplateOrderShipped         = "order_shipped"
	TemplateOrderDelivered       = "order_delivered"
	TemplateLowStockAlert        = "low_stock_alert"
)

// Notification represents a notification record
type Notification struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       *uuid.UUID     `gorm:"type:uuid;index" json:"userId,omitempty"`
	Type         string         `gorm:"type:varchar(50);not null" json:"type"`
	Template     string         `gorm:"type:varchar(100);not null" json:"template"`
	Recipient    string         `gorm:"type:varchar(255);not null" json:"recipient"`
	Subject      string         `gorm:"type:varchar(255);not null" json:"subject"`
	Body         string         `gorm:"type:text;not null" json:"body"`
	Status       string         `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	RetryCount   int            `gorm:"default:0" json:"retryCount"`
	ErrorMessage string         `gorm:"type:text" json:"errorMessage,omitempty"`
	SentAt       *time.Time     `json:"sentAt,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}

// CanRetry checks if notification can be retried
func (n *Notification) CanRetry() bool {
	return n.Status == NotificationStatusFailed && n.RetryCount < 3
}
