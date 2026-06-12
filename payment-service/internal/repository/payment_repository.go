package repository

import (
	"errors"
	"payment-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uuid.UUID) (*models.Payment, error)
	FindByOrderID(orderID uuid.UUID) (*models.Payment, error)
	FindByStripePaymentIntentID(intentID string) (*models.Payment, error)
	Update(payment *models.Payment) error
	CreateRefund(refund *models.Refund) error
	FindRefundByID(id uuid.UUID) (*models.Refund, error)
	FindRefundsByPaymentID(paymentID uuid.UUID) ([]models.Refund, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Refunds").First(&payment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Refunds").First(&payment, "order_id = ?", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByStripePaymentIntentID(intentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Refunds").First(&payment, "stripe_payment_intent_id = ?", intentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) CreateRefund(refund *models.Refund) error {
	return r.db.Create(refund).Error
}

func (r *paymentRepository) FindRefundByID(id uuid.UUID) (*models.Refund, error) {
	var refund models.Refund
	err := r.db.First(&refund, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refund not found")
		}
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRepository) FindRefundsByPaymentID(paymentID uuid.UUID) ([]models.Refund, error) {
	var refunds []models.Refund
	err := r.db.Where("payment_id = ?", paymentID).Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}
