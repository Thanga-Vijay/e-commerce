package consumers

import (
	"context"
	"fmt"
	"log"

	"inventory-service/internal/models"
	"inventory-service/pkg/kafka"

	"gorm.io/gorm"
)

// OrderEventConsumer listens for order events to update inventory
type OrderEventConsumer struct {
	db            *gorm.DB
	kafkaProducer *kafka.Producer
	consumer      *kafka.Consumer
}

// NewOrderEventConsumer creates a new order event consumer for inventory
func NewOrderEventConsumer(brokers []string, db *gorm.DB, producer *kafka.Producer) *OrderEventConsumer {
	config := kafka.ConsumerConfig{
		Brokers: brokers,
		Topic:   "order.confirmed",
		GroupID: "inventory-service-orders",
		Service: "inventory-service",
		Handler: nil,
	}

	oec := &OrderEventConsumer{
		db:            db,
		kafkaProducer: producer,
	}

	config.Handler = oec.handleOrderConfirmed
	oec.consumer = kafka.NewConsumer(config)

	return oec
}

// Start starts the order event consumer
func (oec *OrderEventConsumer) Start(ctx context.Context) error {
	log.Println("Starting inventory order event consumer...")
	return oec.consumer.Start(ctx)
}

// handleOrderConfirmed processes order.confirmed events to deduct inventory
func (oec *OrderEventConsumer) handleOrderConfirmed(ctx context.Context, event map[string]interface{}) error {
	log.Printf("Received order.confirmed event for inventory update: %v", event)

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data format")
	}

	orderID := data["order_id"]
	
	// In production, fetch order items from order service or event data
	// For this example, we'll assume items are in the event
	items, ok := data["items"].([]interface{})
	if !ok {
		log.Printf("No items in order event, skipping inventory update")
		return nil
	}

	tx := oec.db.Begin()

	for _, item := range items {
		itemData := item.(map[string]interface{})
		productID := uint(itemData["product_id"].(float64))
		quantity := int(itemData["quantity"].(float64))

		// Find inventory for product (assuming warehouse_id = 1 for simplicity)
		var inventory models.Inventory
		if err := tx.Where("product_id = ? AND warehouse_id = ?", productID, 1).First(&inventory).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("inventory not found for product %d: %w", productID, err)
		}

		// Check if sufficient stock
		if inventory.Quantity < quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient stock for product %d: available %d, required %d", 
				productID, inventory.Quantity, quantity)
		}

		// Deduct inventory
		previousStock := inventory.Quantity
		inventory.Quantity -= quantity

		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update inventory: %w", err)
		}

		// Publish inventory.updated event
		if oec.kafkaProducer != nil {
			updateEvent := map[string]interface{}{
				"product_id":      productID,
				"warehouse_id":    1,
				"previous_stock":  previousStock,
				"new_stock":       inventory.Quantity,
				"change_type":     "deduction",
				"order_id":        orderID,
			}

			if err := oec.kafkaProducer.PublishEvent(ctx, "inventory.updated", "inventory.updated", updateEvent); err != nil {
				log.Printf("Failed to publish inventory.updated event: %v", err)
			}

			// Check for low stock
			if inventory.Quantity <= inventory.ReorderLevel {
				lowStockEvent := map[string]interface{}{
					"product_id":    productID,
					"warehouse_id":  1,
					"current_stock": inventory.Quantity,
					"threshold":     inventory.ReorderLevel,
				}

				if err := oec.kafkaProducer.PublishEvent(ctx, "inventory.low-stock", "inventory.low-stock", lowStockEvent); err != nil {
					log.Printf("Failed to publish inventory.low-stock event: %v", err)
				}
			}

			// Check for out of stock
			if inventory.Quantity == 0 {
				outOfStockEvent := map[string]interface{}{
					"product_id":   productID,
					"warehouse_id": 1,
				}

				if err := oec.kafkaProducer.PublishEvent(ctx, "inventory.out-of-stock", "inventory.out-of-stock", outOfStockEvent); err != nil {
					log.Printf("Failed to publish inventory.out-of-stock event: %v", err)
				}
			}
		}

		log.Printf("Inventory updated for product %d: %d -> %d", productID, previousStock, inventory.Quantity)
	}

	tx.Commit()
	log.Printf("Inventory successfully updated for order: %v", orderID)
	return nil
}

// OrderCancelledConsumer listens for order.cancelled events to restore inventory
type OrderCancelledConsumer struct {
	db            *gorm.DB
	kafkaProducer *kafka.Producer
	consumer      *kafka.Consumer
}

// NewOrderCancelledConsumer creates a consumer for order cancellations
func NewOrderCancelledConsumer(brokers []string, db *gorm.DB, producer *kafka.Producer) *OrderCancelledConsumer {
	config := kafka.ConsumerConfig{
		Brokers: brokers,
		Topic:   "order.cancelled",
		GroupID: "inventory-service-cancellations",
		Service: "inventory-service",
		Handler: nil,
	}

	occ := &OrderCancelledConsumer{
		db:            db,
		kafkaProducer: producer,
	}

	config.Handler = occ.handleOrderCancelled
	occ.consumer = kafka.NewConsumer(config)

	return occ
}

// Start starts the order cancelled consumer
func (occ *OrderCancelledConsumer) Start(ctx context.Context) error {
	log.Println("Starting inventory order cancellation consumer...")
	return occ.consumer.Start(ctx)
}

// handleOrderCancelled processes order.cancelled events to restore inventory
func (occ *OrderCancelledConsumer) handleOrderCancelled(ctx context.Context, event map[string]interface{}) error {
	log.Printf("Received order.cancelled event for inventory restoration: %v", event)

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data format")
	}

	orderID := data["order_id"]
	
	// Fetch order items and restore inventory
	items, ok := data["items"].([]interface{})
	if !ok {
		log.Printf("No items in cancellation event, skipping inventory restoration")
		return nil
	}

	tx := occ.db.Begin()

	for _, item := range items {
		itemData := item.(map[string]interface{})
		productID := uint(itemData["product_id"].(float64))
		quantity := int(itemData["quantity"].(float64))

		var inventory models.Inventory
		if err := tx.Where("product_id = ? AND warehouse_id = ?", productID, 1).First(&inventory).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("inventory not found for product %d: %w", productID, err)
		}

		previousStock := inventory.Quantity
		inventory.Quantity += quantity

		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to restore inventory: %w", err)
		}

		// Publish inventory.updated event
		if occ.kafkaProducer != nil {
			updateEvent := map[string]interface{}{
				"product_id":     productID,
				"warehouse_id":   1,
				"previous_stock": previousStock,
				"new_stock":      inventory.Quantity,
				"change_type":    "addition",
				"order_id":       orderID,
				"reason":         "order_cancelled",
			}

			occ.kafkaProducer.PublishEvent(ctx, "inventory.updated", "inventory.updated", updateEvent)
		}

		log.Printf("Inventory restored for product %d: %d -> %d", productID, previousStock, inventory.Quantity)
	}

	tx.Commit()
	log.Printf("Inventory successfully restored for cancelled order: %v", orderID)
	return nil
}
