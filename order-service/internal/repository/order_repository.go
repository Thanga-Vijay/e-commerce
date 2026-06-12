package repository

import (
	"errors"
	"fmt"
	"order-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uuid.UUID) (*models.Order, error)
	FindByUserID(userID uuid.UUID, page, limit int, status string) ([]models.Order, int64, error)
	Update(order *models.Order) error
	AddStatusHistory(status *models.OrderStatus) error
	GetOrderNumber() (string, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Create initial status history
		statusHistory := &models.OrderStatus{
			OrderID: order.ID,
			Status:  order.Status,
			Comment: "Order created",
		}
		if err := tx.Create(statusHistory).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").Preload("StatusHistory").First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uuid.UUID, page, limit int, status string) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).Where("user_id = ?", userID)

	// Filter by status if provided
	if status != "" {
		if !models.IsValidStatus(status) {
			return nil, 0, errors.New("invalid status")
		}
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate and fetch orders
	offset := (page - 1) * limit
	err := query.
		Preload("Items").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) AddStatusHistory(status *models.OrderStatus) error {
	return r.db.Create(status).Error
}

func (r *orderRepository) GetOrderNumber() (string, error) {
	var count int64
	if err := r.db.Model(&models.Order{}).Count(&count).Error; err != nil {
		return "", err
	}

	// Generate order number: ORD-YYYY-NNNNNN
	orderNumber := fmt.Sprintf("ORD-2026-%06d", count+1)
	return orderNumber, nil
}
