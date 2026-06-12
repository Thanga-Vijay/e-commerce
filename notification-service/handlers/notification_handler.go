package handlers

import (
	"net/http"
	"notification-service/internal/service"
	"notification-service/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SendNotification sends a notification
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req service.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	notification, err := h.notificationService.SendNotification(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Notification queued successfully", notification)
}

// GetNotificationByID retrieves notification by ID
func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	notification, err := h.notificationService.GetNotificationByID(id)
	if err != nil {
		if err.Error() == "notification not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification retrieved successfully", notification)
}

// GetNotificationsByUserID retrieves notifications for a user
func (h *NotificationHandler) GetNotificationsByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notifications, total, err := h.notificationService.GetNotificationsByUserID(userID.(uuid.UUID), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := gin.H{
		"notifications": notifications,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifications retrieved successfully", response)
}

// ProcessPending processes pending notifications (admin only, typically called by cron)
func (h *NotificationHandler) ProcessPending(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err := h.notificationService.ProcessPending(limit); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pending notifications processed", nil)
}

// RetryFailed retries failed notifications (admin only, typically called by cron)
func (h *NotificationHandler) RetryFailed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err := h.notificationService.RetryFailed(limit); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Failed notifications retried", nil)
}
