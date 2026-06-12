package service

import (
	"errors"
	"fmt"
	"order-service/internal/client"
	"order-service/internal/models"
	"order-service/internal/repository"

	"github.com/google/uuid"
)

const (
	TaxRate      = 0.08  // 8% tax rate
	ShippingCost = 10.00 // Flat shipping cost
)

type CreateOrderRequest struct {
	ShippingAddress models.Address `json:"shippingAddress"`
	BillingAddress  models.Address `json:"billingAddress"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

type OrderService interface {
	CreateOrder(userID uuid.UUID, token string, req CreateOrderRequest) (*models.Order, error)
	GetOrderByID(orderID, userID uuid.UUID, isAdmin bool) (*models.Order, error)
	GetUserOrders(userID uuid.UUID, page, limit int, status string) ([]models.Order, int64, error)
	CancelOrder(orderID, userID uuid.UUID, reason string, isAdmin bool) (*models.Order, error)
	UpdateOrderStatus(orderID uuid.UUID, newStatus, comment string) (*models.Order, error)
}

type orderService struct {
	repo       repository.OrderRepository
	cartClient client.CartClient
}

func NewOrderService(repo repository.OrderRepository, cartClient client.CartClient) OrderService {
	return &orderService{
		repo:       repo,
		cartClient: cartClient,
	}
}

func (s *orderService) CreateOrder(userID uuid.UUID, token string, req CreateOrderRequest) (*models.Order, error) {
	// Validate addresses
	if err := s.validateAddress(req.ShippingAddress); err != nil {
		return nil, fmt.Errorf("invalid shipping address: %w", err)
	}
	if err := s.validateAddress(req.BillingAddress); err != nil {
		return nil, fmt.Errorf("invalid billing address: %w", err)
	}

	// Get cart from Cart Service
	cart, err := s.cartClient.GetCart(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Validate cart has items
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Generate order number
	orderNumber, err := s.repo.GetOrderNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %w", err)
	}

	// Create order items
	orderItems := make([]models.OrderItem, len(cart.Items))
	for i, item := range cart.Items {
		orderItems[i] = models.OrderItem{
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity:     item.Quantity,
			Subtotal:     item.Subtotal,
		}
	}

	// Calculate totals
	subtotal := cart.Subtotal
	tax := subtotal * TaxRate
	total := subtotal + tax + ShippingCost

	// Create order
	order := &models.Order{
		UserID:          userID,
		OrderNumber:     orderNumber,
		Status:          models.OrderStatusPending,
		Subtotal:        subtotal,
		Tax:             tax,
		ShippingCost:    ShippingCost,
		Total:           total,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Items:           orderItems,
	}

	// Save order to database
	if err := s.repo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Clear cart after successful order creation
	if err := s.cartClient.ClearCart(token); err != nil {
		// Log error but don't fail order creation
		fmt.Printf("Warning: failed to clear cart after order creation: %v\n", err)
	}

	// Load order with relationships
	createdOrder, err := s.repo.FindByID(order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load created order: %w", err)
	}

	return createdOrder, nil
}

func (s *orderService) GetOrderByID(orderID, userID uuid.UUID, isAdmin bool) (*models.Order, error) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to view this order
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("unauthorized to view this order")
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID uuid.UUID, page, limit int, status string) ([]models.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.FindByUserID(userID, page, limit, status)
}

func (s *orderService) CancelOrder(orderID, userID uuid.UUID, reason string, isAdmin bool) (*models.Order, error) {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to cancel this order
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("unauthorized to cancel this order")
	}

	// Check if order can be cancelled
	if !order.CanTransitionTo(models.OrderStatusCancelled) {
		return nil, fmt.Errorf("cannot cancel order in %s status", order.Status)
	}

	// Update order status
	order.Status = models.OrderStatusCancelled
	if err := s.repo.Update(order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Add status history
	statusHistory := &models.OrderStatus{
		OrderID: order.ID,
		Status:  models.OrderStatusCancelled,
		Comment: reason,
	}
	if err := s.repo.AddStatusHistory(statusHistory); err != nil {
		return nil, fmt.Errorf("failed to add status history: %w", err)
	}

	// Reload order with updated history
	return s.repo.FindByID(order.ID)
}

func (s *orderService) UpdateOrderStatus(orderID uuid.UUID, newStatus, comment string) (*models.Order, error) {
	// Validate status
	if !models.IsValidStatus(newStatus) {
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}

	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if transition is valid
	if !order.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("cannot transition from %s to %s", order.Status, newStatus)
	}

	// Update order status
	order.Status = newStatus
	if err := s.repo.Update(order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Add status history
	statusHistory := &models.OrderStatus{
		OrderID: order.ID,
		Status:  newStatus,
		Comment: comment,
	}
	if err := s.repo.AddStatusHistory(statusHistory); err != nil {
		return nil, fmt.Errorf("failed to add status history: %w", err)
	}

	// Reload order with updated history
	return s.repo.FindByID(order.ID)
}

func (s *orderService) validateAddress(addr models.Address) error {
	if addr.FirstName == "" || addr.LastName == "" {
		return errors.New("first name and last name are required")
	}
	if addr.AddressLine1 == "" {
		return errors.New("address line 1 is required")
	}
	if addr.City == "" || addr.State == "" || addr.ZipCode == "" {
		return errors.New("city, state, and zip code are required")
	}
	if addr.Country == "" {
		return errors.New("country is required")
	}
	if addr.Phone == "" {
		return errors.New("phone number is required")
	}
	return nil
}
