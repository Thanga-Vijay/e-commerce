package service

import (
	"errors"
	"fmt"
	"payment-service/config"
	"payment-service/internal/models"
	"payment-service/internal/repository"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

type CreatePaymentIntentRequest struct {
	OrderID uuid.UUID `json:"orderId" binding:"required"`
	Amount  float64   `json:"amount" binding:"required,gt=0"`
}

type CreateRefundRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Reason string  `json:"reason"`
}

type PaymentService interface {
	CreatePaymentIntent(userID uuid.UUID, req CreatePaymentIntentRequest) (*models.Payment, error)
	GetPaymentByID(paymentID, userID uuid.UUID, isAdmin bool) (*models.Payment, error)
	GetPaymentByOrderID(orderID, userID uuid.UUID, isAdmin bool) (*models.Payment, error)
	ProcessRefund(paymentID uuid.UUID, req CreateRefundRequest) (*models.Refund, error)
	HandleWebhook(payload []byte, signature string) error
}

type paymentService struct {
	repo   repository.PaymentRepository
	config *config.Config
}

func NewPaymentService(repo repository.PaymentRepository, cfg *config.Config) PaymentService {
	// Initialize Stripe
	stripe.Key = cfg.Stripe.SecretKey
	return &paymentService{
		repo:   repo,
		config: cfg,
	}
}

func (s *paymentService) CreatePaymentIntent(userID uuid.UUID, req CreatePaymentIntentRequest) (*models.Payment, error) {
	// Check if payment already exists for this order
	existingPayment, err := s.repo.FindByOrderID(req.OrderID)
	if err == nil && existingPayment != nil {
		// Payment already exists, return it if it's pending or succeeded
		if existingPayment.Status == models.PaymentStatusPending || existingPayment.Status == models.PaymentStatusSucceeded {
			return existingPayment, nil
		}
	}

	// Convert amount to cents (Stripe uses smallest currency unit)
	amountInCents := int64(req.Amount * 100)

	// Create Stripe payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(s.config.Stripe.Currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"order_id": req.OrderID.String(),
			"user_id":  userID.String(),
		},
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create payment record
	payment := &models.Payment{
		OrderID:               req.OrderID,
		UserID:                userID,
		StripePaymentIntentID: intent.ID,
		Amount:                req.Amount,
		Currency:              s.config.Stripe.Currency,
		Status:                models.PaymentStatusPending,
		ClientSecret:          intent.ClientSecret,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	return payment, nil
}

func (s *paymentService) GetPaymentByID(paymentID, userID uuid.UUID, isAdmin bool) (*models.Payment, error) {
	payment, err := s.repo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to view this payment
	if !isAdmin && payment.UserID != userID {
		return nil, errors.New("unauthorized to view this payment")
	}

	return payment, nil
}

func (s *paymentService) GetPaymentByOrderID(orderID, userID uuid.UUID, isAdmin bool) (*models.Payment, error) {
	payment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to view this payment
	if !isAdmin && payment.UserID != userID {
		return nil, errors.New("unauthorized to view this payment")
	}

	return payment, nil
}

func (s *paymentService) ProcessRefund(paymentID uuid.UUID, req CreateRefundRequest) (*models.Refund, error) {
	// Get payment
	payment, err := s.repo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}

	// Check if payment is succeeded
	if payment.Status != models.PaymentStatusSucceeded {
		return nil, errors.New("can only refund succeeded payments")
	}

	// Check if refund amount is valid
	if req.Amount > payment.Amount {
		return nil, errors.New("refund amount cannot exceed payment amount")
	}

	// Check total refunded amount
	existingRefunds, err := s.repo.FindRefundsByPaymentID(payment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing refunds: %w", err)
	}

	totalRefunded := 0.0
	for _, r := range existingRefunds {
		if r.Status == models.RefundStatusSucceeded {
			totalRefunded += r.Amount
		}
	}

	if totalRefunded+req.Amount > payment.Amount {
		return nil, fmt.Errorf("total refund amount would exceed payment amount")
	}

	// Create Stripe refund
	amountInCents := int64(req.Amount * 100)
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(payment.StripePaymentIntentID),
		Amount:        stripe.Int64(amountInCents),
	}

	if req.Reason != "" {
		params.Reason = stripe.String(stripe.RefundReasonRequestedByCustomer)
	}

	stripeRefund, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe refund: %w", err)
	}

	// Create refund record
	refundRecord := &models.Refund{
		PaymentID:      payment.ID,
		StripeRefundID: stripeRefund.ID,
		Amount:         req.Amount,
		Reason:         req.Reason,
		Status:         models.RefundStatusSucceeded,
	}

	if err := s.repo.CreateRefund(refundRecord); err != nil {
		return nil, fmt.Errorf("failed to create refund record: %w", err)
	}

	return refundRecord, nil
}

func (s *paymentService) HandleWebhook(payload []byte, signature string) error {
	// Verify webhook signature
	event, err := webhook.ConstructEvent(payload, signature, s.config.Stripe.WebhookSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		return s.handlePaymentIntentSucceeded(event)
	case "payment_intent.payment_failed":
		return s.handlePaymentIntentFailed(event)
	case "charge.refunded":
		return s.handleChargeRefunded(event)
	default:
		// Unhandled event type
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	return nil
}

func (s *paymentService) handlePaymentIntentSucceeded(event stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := event.DataObjectJSON.UnmarshalJSON([]byte(event.Data.Raw), &intent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	// Find payment by Stripe payment intent ID
	payment, err := s.repo.FindByStripePaymentIntentID(intent.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment status
	payment.Status = models.PaymentStatusSucceeded
	if intent.PaymentMethod != nil {
		payment.PaymentMethod = string(intent.PaymentMethod.Type)
	}

	if err := s.repo.Update(payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	fmt.Printf("Payment succeeded: %s\n", payment.ID)
	return nil
}

func (s *paymentService) handlePaymentIntentFailed(event stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := event.DataObjectJSON.UnmarshalJSON([]byte(event.Data.Raw), &intent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	// Find payment by Stripe payment intent ID
	payment, err := s.repo.FindByStripePaymentIntentID(intent.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment status
	payment.Status = models.PaymentStatusFailed

	if err := s.repo.Update(payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	fmt.Printf("Payment failed: %s\n", payment.ID)
	return nil
}

func (s *paymentService) handleChargeRefunded(event stripe.Event) error {
	var charge stripe.Charge
	if err := event.DataObjectJSON.UnmarshalJSON([]byte(event.Data.Raw), &charge); err != nil {
		return fmt.Errorf("failed to parse charge: %w", err)
	}

	fmt.Printf("Charge refunded: %s\n", charge.ID)
	return nil
}
