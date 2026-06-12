package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment statuses
const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
	PaymentStatusCanceled  = "canceled"
)

// Refund statuses
const (
	RefundStatusPending   = "pending"
	RefundStatusSucceeded = "succeeded"
	RefundStatusFailed    = "failed"
)

// Payment represents a payment transaction
type Payment struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID                uuid.UUID      `gorm:"type:uuid;unique;not null;index" json:"orderId"`
	UserID                 uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	StripePaymentIntentID  string         `gorm:"type:varchar(255);unique" json:"stripePaymentIntentId,omitempty"`
	Amount                 float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency               string         `gorm:"type:varchar(3);default:'usd'" json:"currency"`
	Status                 string         `gorm:"type:varchar(20);not null;index" json:"status"`
	PaymentMethod          string         `gorm:"type:varchar(50)" json:"paymentMethod,omitempty"`
	ClientSecret           string         `gorm:"type:varchar(500)" json:"clientSecret,omitempty"`
	Refunds                []Refund       `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE" json:"refunds,omitempty"`
	CreatedAt              time.Time      `json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// Refund represents a payment refund
type Refund struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PaymentID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"paymentId"`
	StripeRefundID  string         `gorm:"type:varchar(255);unique" json:"stripeRefundId,omitempty"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Reason          string         `gorm:"type:text" json:"reason,omitempty"`
	Status          string         `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt       time.Time      `json:"createdAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Payment
func (Payment) TableName() string {
	return "payments"
}

// TableName specifies the table name for Refund
func (Refund) TableName() string {
	return "refunds"
}

// IsValidPaymentStatus checks if a payment status is valid
func IsValidPaymentStatus(status string) bool {
	validStatuses := []string{
		PaymentStatusPending,
		PaymentStatusSucceeded,
		PaymentStatusFailed,
		PaymentStatusCanceled,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidRefundStatus checks if a refund status is valid
func IsValidRefundStatus(status string) bool {
	validStatuses := []string{
		RefundStatusPending,
		RefundStatusSucceeded,
		RefundStatusFailed,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
