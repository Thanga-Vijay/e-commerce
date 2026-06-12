package handlers

import (
	"io"
	"net/http"
	"payment-service/internal/service"
	"payment-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePaymentIntent creates a Stripe payment intent
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req service.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := h.paymentService.CreatePaymentIntent(userID.(uuid.UUID), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Don't send full payment details to client, only what's needed
	response := gin.H{
		"paymentId":    payment.ID,
		"clientSecret": payment.ClientSecret,
		"amount":       payment.Amount,
		"currency":     payment.Currency,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Payment intent created successfully", response)
}

// GetPaymentByID retrieves payment details
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.paymentService.GetPaymentByID(paymentID, userID.(uuid.UUID), isAdmin)
	if err != nil {
		if err.Error() == "payment not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "unauthorized to view this payment" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Remove sensitive data before sending to client
	payment.ClientSecret = ""

	utils.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// GetPaymentByOrderID retrieves payment for an order
func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(orderID, userID.(uuid.UUID), isAdmin)
	if err != nil {
		if err.Error() == "payment not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "unauthorized to view this payment" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Remove sensitive data before sending to client
	payment.ClientSecret = ""

	utils.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// ProcessRefund processes a refund (admin only)
func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var req service.CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	refund, err := h.paymentService.ProcessRefund(paymentID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Refund processed successfully", refund)
}

// HandleWebhook handles Stripe webhooks
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to read request body")
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Missing Stripe signature")
		return
	}

	if err := h.paymentService.HandleWebhook(payload, signature); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
