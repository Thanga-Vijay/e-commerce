package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"notification-service/internal/services"
	"notification-service/pkg/kafka"
)

// EmailConsumer listens for events that require email notifications
type EmailConsumer struct {
	emailService *services.EmailService
	consumer     *kafka.Consumer
}

// NewEmailConsumer creates a new email consumer
func NewEmailConsumer(brokers []string, emailService *services.EmailService) *EmailConsumer {
	config := kafka.ConsumerConfig{
		Brokers: brokers,
		Topic:   "notification.email.send",
		GroupID: "notification-service-email",
		Service: "notification-service",
		Handler: nil, // Will be set below
	}

	ec := &EmailConsumer{
		emailService: emailService,
	}

	// Set handler
	config.Handler = ec.handleEmailEvent
	ec.consumer = kafka.NewConsumer(config)

	return ec
}

// Start starts the email consumer
func (ec *EmailConsumer) Start(ctx context.Context) error {
	log.Println("Starting email consumer...")
	return ec.consumer.Start(ctx)
}

// handleEmailEvent processes email send events
func (ec *EmailConsumer) handleEmailEvent(ctx context.Context, event map[string]interface{}) error {
	log.Printf("Received email event: %v", event)

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data format")
	}

	to, _ := data["to"].(string)
	subject, _ := data["subject"].(string)
	template, _ := data["template"].(string)
	emailData, _ := data["data"].(map[string]interface{})

	if to == "" || subject == "" {
		return fmt.Errorf("missing required fields: to or subject")
	}

	// Send email
	var body string
	if template != "" {
		// Load template and render with data
		body = ec.renderTemplate(template, emailData)
	} else {
		// Use body from data
		bodyData, _ := emailData["body"].(string)
		body = bodyData
	}

	if err := ec.emailService.SendEmail(to, subject, body); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to: %s", to)
	return nil
}

// renderTemplate renders email template with data
func (ec *EmailConsumer) renderTemplate(template string, data map[string]interface{}) string {
	// Simple template rendering - in production use proper template engine
	dataJSON, _ := json.Marshal(data)
	return fmt.Sprintf("Template: %s\nData: %s", template, string(dataJSON))
}

// OrderEventConsumer listens for order events to send notifications
type OrderEventConsumer struct {
	emailService  *services.EmailService
	kafkaProducer *kafka.Producer
	consumer      *kafka.Consumer
}

// NewOrderEventConsumer creates a new order event consumer
func NewOrderEventConsumer(brokers []string, emailService *services.EmailService, producer *kafka.Producer) *OrderEventConsumer {
	config := kafka.ConsumerConfig{
		Brokers: brokers,
		Topic:   "order.created",
		GroupID: "notification-service-orders",
		Service: "notification-service",
		Handler: nil,
	}

	oec := &OrderEventConsumer{
		emailService:  emailService,
		kafkaProducer: producer,
	}

	config.Handler = oec.handleOrderCreated
	oec.consumer = kafka.NewConsumer(config)

	return oec
}

// Start starts the order event consumer
func (oec *OrderEventConsumer) Start(ctx context.Context) error {
	log.Println("Starting order event consumer...")
	return oec.consumer.Start(ctx)
}

// handleOrderCreated processes order.created events
func (oec *OrderEventConsumer) handleOrderCreated(ctx context.Context, event map[string]interface{}) error {
	log.Printf("Received order.created event: %v", event)

	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data format")
	}

	orderID := data["order_id"]
	userID := data["user_id"]
	totalAmount := data["total_amount"]

	// In production, fetch user email from user service or database
	userEmail := fmt.Sprintf("user_%v@example.com", userID)

	// Send order confirmation email
	subject := fmt.Sprintf("Order Confirmation - Order #%v", orderID)
	body := fmt.Sprintf(`
		Dear Customer,

		Thank you for your order!

		Order ID: %v
		Total Amount: $%.2f

		Your order is being processed and you will receive a shipping confirmation soon.

		Best regards,
		E-Commerce Team
	`, orderID, totalAmount)

	if err := oec.emailService.SendEmail(userEmail, subject, body); err != nil {
		return fmt.Errorf("failed to send order confirmation email: %w", err)
	}

	// Publish notification.email.sent event
	if oec.kafkaProducer != nil {
		sentEvent := map[string]interface{}{
			"to":         userEmail,
			"subject":    subject,
			"message_id": fmt.Sprintf("order-%v", orderID),
		}

		if err := oec.kafkaProducer.PublishEvent(ctx, "notification.email.sent", "notification.email.sent", sentEvent); err != nil {
			log.Printf("Failed to publish email.sent event: %v", err)
		}
	}

	log.Printf("Order confirmation email sent for order: %v", orderID)
	return nil
}
