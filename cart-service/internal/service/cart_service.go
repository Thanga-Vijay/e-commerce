package service

import (
	"context"
	"errors"
	"time"

	"github.com/ecommerce/cart-service/internal/cache"
	"github.com/ecommerce/cart-service/internal/client"
	"github.com/ecommerce/cart-service/internal/models"
	"github.com/ecommerce/cart-service/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrCartNotFound     = errors.New("cart not found")
	ErrItemNotFound     = errors.New("item not found in cart")
	ErrProductNotFound  = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity  = errors.New("invalid quantity")
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req *models.AddToCartRequest) (*models.Cart, error)
	UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *models.UpdateCartItemRequest) (*models.Cart, error)
	RemoveFromCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*models.Cart, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type cartService struct {
	cartRepo      repository.CartRepository
	cache         cache.CartCache
	productClient client.ProductClient
	cacheTTL      time.Duration
}

func NewCartService(
	cartRepo repository.CartRepository,
	cache cache.CartCache,
	productClient client.ProductClient,
	cacheTTL time.Duration,
) CartService {
	return &cartService{
		cartRepo:      cartRepo,
		cache:         cache,
		productClient: productClient,
		cacheTTL:      cacheTTL,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	// Try cache first
	cart, err := s.cache.GetCart(ctx, userID)
	if err == nil && cart != nil {
		return cart, nil
	}

	// Fetch from database
	cart, err = s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If no cart exists, create a new one
	if cart == nil {
		cart = &models.Cart{
			UserID: userID,
			Items:  []models.CartItem{},
		}
		if err := s.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	// Cache the cart
	_ = s.cache.SetCart(ctx, cart, s.cacheTTL)

	return cart, nil
}

func (s *cartService) AddToCart(ctx context.Context, userID uuid.UUID, req *models.AddToCartRequest) (*models.Cart, error) {
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

	// Validate product is active and has stock
	if !product.IsActive {
		return nil, errors.New("product is not available")
	}
	if product.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	// Get or create cart
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.FindItemByProductID(ctx, cart.ID, productID)
	if err != nil {
		return nil, err
	}

	if existingItem != nil {
		// Update quantity
		existingItem.Quantity += req.Quantity
		existingItem.CalculateSubtotal()
		
		if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		// Add new item
		productImage := ""
		if len(product.Images) > 0 {
			productImage = product.Images[0]
		}

		item := &models.CartItem{
			CartID:       cart.ID,
			ProductID:    productID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			ProductImage: productImage,
			Quantity:     req.Quantity,
		}
		item.CalculateSubtotal()

		if err := s.cartRepo.AddItem(ctx, item); err != nil {
			return nil, err
		}
	}

	// Invalidate cache and fetch updated cart
	_ = s.cache.DeleteCart(ctx, userID)
	cart, err = s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache updated cart
	_ = s.cache.SetCart(ctx, cart, s.cacheTTL)

	return cart, nil
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *models.UpdateCartItemRequest) (*models.Cart, error) {
	if req.Quantity < 1 {
		return nil, ErrInvalidQuantity
	}

	// Get cart
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find item in cart
	var item *models.CartItem
	for i := range cart.Items {
		if cart.Items[i].ID == itemID {
			item = &cart.Items[i]
			break
		}
	}

	if item == nil {
		return nil, ErrItemNotFound
	}

	// Get product to check stock
	product, err := s.productClient.GetProduct(ctx, item.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	// Update quantity
	item.Quantity = req.Quantity
	item.CalculateSubtotal()

	if err := s.cartRepo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	// Invalidate cache and fetch updated cart
	_ = s.cache.DeleteCart(ctx, userID)
	cart, err = s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache updated cart
	_ = s.cache.SetCart(ctx, cart, s.cacheTTL)

	return cart, nil
}

func (s *cartService) RemoveFromCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) (*models.Cart, error) {
	// Get cart
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if item exists in cart
	found := false
	for _, item := range cart.Items {
		if item.ID == itemID {
			found = true
			break
		}
	}

	if !found {
		return nil, ErrItemNotFound
	}

	// Remove item
	if err := s.cartRepo.RemoveItem(ctx, itemID); err != nil {
		return nil, err
	}

	// Invalidate cache and fetch updated cart
	_ = s.cache.DeleteCart(ctx, userID)
	cart, err = s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache updated cart
	_ = s.cache.SetCart(ctx, cart, s.cacheTTL)

	return cart, nil
}

func (s *cartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	// Get cart
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Delete cart
	if err := s.cartRepo.Delete(ctx, cart.ID); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.DeleteCart(ctx, userID)

	return nil
}
