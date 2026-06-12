package repository

import (
	"context"
	"errors"

	"github.com/ecommerce/wishlist-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error)
	Create(ctx context.Context, wishlist *models.Wishlist) error
	Update(ctx context.Context, wishlist *models.Wishlist) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddItem(ctx context.Context, item *models.WishlistItem) error
	RemoveItem(ctx context.Context, itemID uuid.UUID) error
	FindItemByID(ctx context.Context, itemID uuid.UUID) (*models.WishlistItem, error)
	FindItemByProductID(ctx context.Context, wishlistID, productID uuid.UUID) (*models.WishlistItem, error)
}

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		First(&wishlist).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wishlist, err
}

func (r *wishlistRepository) Create(ctx context.Context, wishlist *models.Wishlist) error {
	return r.db.WithContext(ctx).Create(wishlist).Error
}

func (r *wishlistRepository) Update(ctx context.Context, wishlist *models.Wishlist) error {
	return r.db.WithContext(ctx).Save(wishlist).Error
}

func (r *wishlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Wishlist{}, id).Error
}

func (r *wishlistRepository) AddItem(ctx context.Context, item *models.WishlistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *wishlistRepository) RemoveItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.WishlistItem{}, itemID).Error
}

func (r *wishlistRepository) FindItemByID(ctx context.Context, itemID uuid.UUID) (*models.WishlistItem, error) {
	var item models.WishlistItem
	err := r.db.WithContext(ctx).Where("id = ?", itemID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *wishlistRepository) FindItemByProductID(ctx context.Context, wishlistID, productID uuid.UUID) (*models.WishlistItem, error) {
	var item models.WishlistItem
	err := r.db.WithContext(ctx).
		Where("wishlist_id = ? AND product_id = ?", wishlistID, productID).
		First(&item).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}
