package service

import (
	"context"
	"errors"

	"github.com/ecommerce/wishlist-service/internal/client"
	"github.com/ecommerce/wishlist-service/internal/models"
	"github.com/ecommerce/wishlist-service/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrWishlistNotFound = errors.New("wishlist not found")
	ErrItemNotFound     = errors.New("item not found in wishlist")
	ErrProductNotFound  = errors.New("product not found")
	ErrItemExists       = errors.New("item already in wishlist")
)

type WishlistService interface {
	GetWishlist(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error)
	AddToWishlist(ctx context.Context, userID uuid.UUID, req *models.AddToWishlistRequest) (*models.Wishlist, error)
	RemoveFromWishlist(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*models.Wishlist, error)
	ClearWishlist(ctx context.Context, userID uuid.UUID) error
	MoveToCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, token string, req *models.MoveToCartRequest) error
}

type wishlistService struct {
	wishlistRepo  repository.WishlistRepository
	productClient client.ProductClient
	cartClient    client.CartClient
}

func NewWishlistService(
	wishlistRepo repository.WishlistRepository,
	productClient client.ProductClient,
	cartClient client.CartClient,
) WishlistService {
	return &wishlistService{
		wishlistRepo:  wishlistRepo,
		productClient: productClient,
		cartClient:    cartClient,
	}
}

func (s *wishlistService) GetWishlist(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error) {
	wishlist, err := s.wishlistRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If no wishlist exists, create a new one
	if wishlist == nil {
		wishlist = &models.Wishlist{
			UserID: userID,
			Items:  []models.WishlistItem{},
		}
		if err := s.wishlistRepo.Create(ctx, wishlist); err != nil {
			return nil, err
		}
	}

	return wishlist, nil
}

func (s *wishlistService) AddToWishlist(ctx context.Context, userID uuid.UUID, req *models.AddToWishlistRequest) (*models.Wishlist, error) {
	// Parse product ID
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	// Get product details from Product Service
	product, err := s.productClient.GetProduct(ctx, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Validate product is active
	if !product.IsActive {
		return nil, errors.New("product is not available")
	}

	// Get or create wishlist
	wishlist, err := s.GetWishlist(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists in wishlist
	existingItem, err := s.wishlistRepo.FindItemByProductID(ctx, wishlist.ID, productID)
	if err != nil {
		return nil, err
	}

	if existingItem != nil {
		return nil, ErrItemExists
	}

	// Add new item
	productImage := ""
	if len(product.Images) > 0 {
		productImage = product.Images[0]
	}

	item := &models.WishlistItem{
		WishlistID:   wishlist.ID,
		ProductID:    productID,
		ProductName:  product.Name,
		ProductPrice: product.Price,
		ProductImage: productImage,
	}

	if err := s.wishlistRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	// Fetch updated wishlist
	wishlist, err = s.wishlistRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s *wishlistService) RemoveFromWishlist(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*models.Wishlist, error) {
	// Get wishlist
	wishlist, err := s.GetWishlist(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if item exists in wishlist
	found := false
	for _, item := range wishlist.Items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		return nil, ErrItemNotFound
	}

	// Remove item
	if err := s.wishlistRepo.RemoveItem(ctx, itemID); err != nil {
		return nil, err
	}

	// Fetch updated wishlist
	wishlist, err = s.wishlistRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s *wishlistService) ClearWishlist(ctx context.Context, userID uuid.UUID) error {
	// Get wishlist
	wishlist, err := s.GetWishlist(ctx, userID)
	if err != nil {
		return err
	}

	// Delete wishlist
	if err := s.wishlistRepo.Delete(ctx, wishlist.ID); err != nil {
		return err
	}

	return nil
}

func (s *wishlistService) MoveToCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, token string, req *models.MoveToCartRequest) error {
	// Get wishlist item
	item, err := s.wishlistRepo.FindItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}

	// Verify item belongs to user's wishlist
	wishlist, err := s.GetWishlist(ctx, userID)
	if err != nil {
		return err
	}

	if item.WishlistID != wishlist.ID {
		return ErrItemNotFound
	}

	// Add to cart using Cart Service
	if err := s.cartClient.AddToCart(ctx, token, item.ProductID, req.Quantity); err != nil {
		return err
	}

	// Remove from wishlist
	if err := s.wishlistRepo.RemoveItem(ctx, itemID); err != nil {
		return err
	}

	return nil
}
