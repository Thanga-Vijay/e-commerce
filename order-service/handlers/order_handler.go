package handlers

import (
	"net/http"
	"order-service/internal/service"
	"order-service/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder creates a new order from cart
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	token, exists := c.Get("token")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Token not found")
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(userID.(uuid.UUID), token.(string), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Order created successfully", order)
}

// GetOrders retrieves user's orders with pagination
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	orders, total, err := h.orderService.GetUserOrders(userID.(uuid.UUID), page, limit, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Orders retrieved successfully", response)
}

// GetOrderByID retrieves a specific order
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.orderService.GetOrderByID(orderID, userID.(uuid.UUID), isAdmin)
	if err != nil {
		if err.Error() == "order not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "unauthorized to view this order" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req service.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Reason = "Cancelled by user"
	}

	order, err := h.orderService.CancelOrder(orderID, userID.(uuid.UUID), req.Reason, isAdmin)
	if err != nil {
		if err.Error() == "order not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "unauthorized to cancel this order" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order cancelled successfully", order)
}

// TrackOrder retrieves order tracking information
func (h *OrderHandler) TrackOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.orderService.GetOrderByID(orderID, userID.(uuid.UUID), isAdmin)
	if err != nil {
		if err.Error() == "order not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "unauthorized to view this order" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Build tracking response
	trackingInfo := gin.H{
		"orderId":      order.ID,
		"orderNumber":  order.OrderNumber,
		"currentStatus": order.Status,
		"timeline":     order.StatusHistory,
	}

	utils.SuccessResponse(c, http.StatusOK, "Tracking info retrieved successfully", trackingInfo)
}

// UpdateOrderStatus updates order status (admin only)
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Status  string `json:"status" binding:"required"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderService.UpdateOrderStatus(orderID, req.Status, req.Comment)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Order status updated successfully", order)
}
