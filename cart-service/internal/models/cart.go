package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Items     []CartItem     `gorm:"foreignKey:CartID" json:"items"`
	ExpiresAt time.Time      `gorm:"not null" json:"expiresAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CartItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CartID        uuid.UUID      `gorm:"type:uuid;not null" json:"cartId"`
	ProductID     uuid.UUID      `gorm:"type:uuid;not null" json:"productId"`
	ProductName   string         `gorm:"not null" json:"productName"`
	ProductPrice  float64        `gorm:"not null" json:"productPrice"`
	ProductImage  string         `json:"productImage"`
	Quantity      int            `gorm:"not null" json:"quantity"`
	Subtotal      float64        `gorm:"not null" json:"subtotal"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hooks
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.ExpiresAt.IsZero() {
		// Cart expires after 30 days
		c.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
	}
	return nil
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}

// CalculateSubtotal calculates the subtotal for a cart item
func (ci *CartItem) CalculateSubtotal() {
	ci.Subtotal = ci.ProductPrice * float64(ci.Quantity)
}
