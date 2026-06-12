package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Items     []WishlistItem   `gorm:"foreignKey:WishlistID" json:"items"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}

type WishlistItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WishlistID    uuid.UUID      `gorm:"type:uuid;not null" json:"wishlistId"`
	ProductID     uuid.UUID      `gorm:"type:uuid;not null" json:"productId"`
	ProductName   string         `gorm:"not null" json:"productName"`
	ProductPrice  float64        `gorm:"not null" json:"productPrice"`
	ProductImage  string         `json:"productImage"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hooks
func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (wi *WishlistItem) BeforeCreate(tx *gorm.DB) error {
	if wi.ID == uuid.Nil {
		wi.ID = uuid.New()
	}
	return nil
}
