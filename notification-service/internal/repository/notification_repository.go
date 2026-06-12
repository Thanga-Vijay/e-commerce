package repository

import (
	"errors"
	"notification-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByID(id uuid.UUID) (*models.Notification, error)
	FindByUserID(userID uuid.UUID, page, limit int) ([]models.Notification, int64, error)
	FindPending(limit int) ([]models.Notification, error)
	FindFailed(limit int) ([]models.Notification, error)
	Update(notification *models.Notification) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) FindByUserID(userID uuid.UUID, page, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate and fetch
	offset := (page - 1) * limit
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepository) FindPending(limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.
		Where("status = ?", models.NotificationStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) FindFailed(limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.
		Where("status = ? AND retry_count < ?", models.NotificationStatusFailed, 3).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}
