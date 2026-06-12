package handlers

import (
	"context"
	"net/http"
	"time"

	"order-service/internal/models"
	"order-service/pkg/events"
	"order-service/pkg/kafka"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB            *gorm.DB
	KafkaProducer *kafka.Producer
}

// CreateOrder creates an order and publishes event
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
		Items           []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate total
	var totalAmount float64
	for _, item := range input.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Create order
	order := models.Order{
		UserID:          userID,
		Status:          "pending",
		TotalAmount:     totalAmount,
		ShippingAddress: input.ShippingAddress,
	}

	tx := h.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items
	for _, item := range input.Items {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}
	}

	tx.Commit()

	// Publish order.created event
	if h.KafkaProducer != nil {
		event := events.OrderCreatedEvent{
			BaseEvent: events.BaseEvent{
				ID:        uuid.New().String(),
				Type:      events.EventOrderCreated,
				Source:    "order-service",
				Timestamp: time.Now().UTC(),
				Version:   "1.0",
			},
			OrderID:      order.ID,
			UserID:       userID,
			TotalAmount:  totalAmount,
			Status:       order.Status,
			ItemCount:    len(input.Items),
			ShippingAddr: input.ShippingAddress,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.KafkaProducer.PublishEvent(ctx, "order.created", events.EventOrderCreated, event); err != nil {
			c.Writer.Header().Add("X-Event-Warning", "Failed to publish event")
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"order_id":     order.ID,
		"status":       order.Status,
		"total_amount": totalAmount,
		"item_count":   len(input.Items),
	})
}

// UpdateOrderStatus updates order status and publishes events
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := h.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	oldStatus := order.Status
	order.Status = input.Status

	if err := h.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	// Publish appropriate event based on status
	if h.KafkaProducer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var eventType string
		var topic string

		switch input.Status {
		case "confirmed":
			eventType = events.EventOrderConfirmed
			topic = "order.confirmed"
			event := events.OrderConfirmedEvent{
				BaseEvent: events.BaseEvent{
					ID:        uuid.New().String(),
					Type:      eventType,
					Source:    "order-service",
					Timestamp: time.Now().UTC(),
					Version:   "1.0",
				},
				OrderID:     order.ID,
				UserID:      order.UserID,
				TotalAmount: order.TotalAmount,
			}
			h.KafkaProducer.PublishEvent(ctx, topic, eventType, event)

		case "shipped":
			eventType = events.EventOrderShipped
			topic = "order.shipped"
			event := events.OrderShippedEvent{
				BaseEvent: events.BaseEvent{
					ID:        uuid.New().String(),
					Type:      eventType,
					Source:    "order-service",
					Timestamp: time.Now().UTC(),
					Version:   "1.0",
				},
				OrderID: order.ID,
				UserID:  order.UserID,
			}
			h.KafkaProducer.PublishEvent(ctx, topic, eventType, event)

		case "delivered":
			eventType = events.EventOrderDelivered
			topic = "order.delivered"
			event := events.OrderDeliveredEvent{
				BaseEvent: events.BaseEvent{
					ID:        uuid.New().String(),
					Type:      eventType,
					Source:    "order-service",
					Timestamp: time.Now().UTC(),
					Version:   "1.0",
				},
				OrderID:     order.ID,
				UserID:      order.UserID,
				DeliveredAt: time.Now().UTC(),
			}
			h.KafkaProducer.PublishEvent(ctx, topic, eventType, event)

		case "cancelled":
			eventType = events.EventOrderCancelled
			topic = "order.cancelled"
			event := events.OrderCancelledEvent{
				BaseEvent: events.BaseEvent{
					ID:        uuid.New().String(),
					Type:      eventType,
					Source:    "order-service",
					Timestamp: time.Now().UTC(),
					Version:   "1.0",
				},
				OrderID:      order.ID,
				UserID:       order.UserID,
				RefundAmount: order.TotalAmount,
			}
			h.KafkaProducer.PublishEvent(ctx, topic, eventType, event)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":   order.ID,
		"old_status": oldStatus,
		"new_status": order.Status,
	})
}
