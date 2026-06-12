package repository

import (
	"context"
	"errors"

	"github.com/ecommerce/cart-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	Create(ctx context.Context, cart *models.Cart) error
	Update(ctx context.Context, cart *models.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddItem(ctx context.Context, item *models.CartItem) error
	UpdateItem(ctx context.Context, item *models.CartItem) error
	RemoveItem(ctx context.Context, itemID uuid.UUID) error
	FindItemByProductID(ctx context.Context, cartID, productID uuid.UUID) (*models.CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		First(&cart).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cart, err
}

func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

func (r *cartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Delete cart and all items (cascade)
	return r.db.WithContext(ctx).Delete(&models.Cart{}, id).Error
}

func (r *cartRepository) AddItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *cartRepository) UpdateItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, itemID).Error
}

func (r *cartRepository) FindItemByProductID(ctx context.Context, cartID, productID uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}
