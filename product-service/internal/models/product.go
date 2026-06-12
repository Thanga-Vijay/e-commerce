package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SKU         string         `gorm:"uniqueIndex;not null" json:"sku"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	CategoryID  uuid.UUID      `gorm:"type:uuid" json:"categoryId"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Images      []string       `gorm:"type:text[]" json:"images"`
	Stock       int            `gorm:"default:0" json:"stock"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Reviews     []Review       `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	ParentID    *uuid.UUID     `gorm:"type:uuid" json:"parentId,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Review struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null" json:"productId"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Title     string         `json:"title"`
	Comment   string         `json:"comment"`
	IsVerified bool          `gorm:"default:false" json:"isVerified"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hooks
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
