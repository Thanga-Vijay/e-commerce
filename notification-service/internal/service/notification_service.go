package service

import (
	"fmt"
	"notification-service/internal/mailer"
	"notification-service/internal/models"
	"notification-service/internal/repository"
	"notification-service/internal/templates"
	"time"

	"github.com/google/uuid"
)

type SendNotificationRequest struct {
	UserID    *uuid.UUID             `json:"userId"`
	Type      string                 `json:"type" binding:"required"`
	Template  string                 `json:"template" binding:"required"`
	Recipient string                 `json:"recipient" binding:"required,email"`
	Data      map[string]interface{} `json:"data" binding:"required"`
}

type NotificationService interface {
	SendNotification(req SendNotificationRequest) (*models.Notification, error)
	GetNotificationByID(id uuid.UUID) (*models.Notification, error)
	GetNotificationsByUserID(userID uuid.UUID, page, limit int) ([]models.Notification, int64, error)
	ProcessPending(limit int) error
	RetryFailed(limit int) error
}

type notificationService struct {
	repo   repository.NotificationRepository
	mailer mailer.Mailer
}

func NewNotificationService(repo repository.NotificationRepository, m mailer.Mailer) NotificationService {
	return &notificationService{
		repo:   repo,
		mailer: m,
	}
}

func (s *notificationService) SendNotification(req SendNotificationRequest) (*models.Notification, error) {
	// Validate notification type
	if req.Type != models.NotificationTypeEmail {
		return nil, fmt.Errorf("unsupported notification type: %s", req.Type)
	}

	// Generate subject based on template
	subject := s.getSubject(req.Template)

	// Render template with data
	body, err := templates.RenderTemplate(req.Template, req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Create notification record
	notification := &models.Notification{
		UserID:    req.UserID,
		Type:      req.Type,
		Template:  req.Template,
		Recipient: req.Recipient,
		Subject:   subject,
		Body:      body,
		Status:    models.NotificationStatusPending,
	}

	if err := s.repo.Create(notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Try to send immediately
	if err := s.sendEmail(notification); err != nil {
		// Mark as failed but don't return error - will retry later
		notification.Status = models.NotificationStatusFailed
		notification.ErrorMessage = err.Error()
		s.repo.Update(notification)
	}

	return s.repo.FindByID(notification.ID)
}

func (s *notificationService) GetNotificationByID(id uuid.UUID) (*models.Notification, error) {
	return s.repo.FindByID(id)
}

func (s *notificationService) GetNotificationsByUserID(userID uuid.UUID, page, limit int) ([]models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindByUserID(userID, page, limit)
}

func (s *notificationService) ProcessPending(limit int) error {
	if limit <= 0 {
		limit = 10
	}

	notifications, err := s.repo.FindPending(limit)
	if err != nil {
		return fmt.Errorf("failed to find pending notifications: %w", err)
	}

	for _, notification := range notifications {
		if err := s.sendEmail(&notification); err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.ErrorMessage = err.Error()
			notification.RetryCount++
			s.repo.Update(&notification)
		}
	}

	return nil
}

func (s *notificationService) RetryFailed(limit int) error {
	if limit <= 0 {
		limit = 10
	}

	notifications, err := s.repo.FindFailed(limit)
	if err != nil {
		return fmt.Errorf("failed to find failed notifications: %w", err)
	}

	for _, notification := range notifications {
		if !notification.CanRetry() {
			continue
		}

		// Exponential backoff: wait 2^retryCount minutes
		waitTime := time.Duration(1<<notification.RetryCount) * time.Minute
		if time.Since(notification.CreatedAt) < waitTime {
			continue
		}

		if err := s.sendEmail(&notification); err != nil {
			notification.ErrorMessage = err.Error()
			notification.RetryCount++
			s.repo.Update(&notification)
		}
	}

	return nil
}

func (s *notificationService) sendEmail(notification *models.Notification) error {
	if err := s.mailer.SendEmail(notification.Recipient, notification.Subject, notification.Body); err != nil {
		return err
	}

	// Update notification status
	now := time.Now()
	notification.Status = models.NotificationStatusSent
	notification.SentAt = &now
	notification.ErrorMessage = ""

	return s.repo.Update(notification)
}

func (s *notificationService) getSubject(template string) string {
	subjects := map[string]string{
		models.TemplateWelcome:           "Welcome to E-Commerce Platform!",
		models.TemplateEmailVerification: "Verify Your Email Address",
		models.TemplatePasswordReset:     "Reset Your Password",
		models.TemplateOrderConfirmation: "Order Confirmation",
		models.TemplatePaymentReceipt:    "Payment Receipt",
		models.TemplateOrderShipped:      "Your Order Has Shipped",
		models.TemplateOrderDelivered:    "Your Order Has Been Delivered",
		models.TemplateLowStockAlert:     "Low Stock Alert",
	}

	if subject, exists := subjects[template]; exists {
		return subject
	}

	return "Notification from E-Commerce Platform"
}
