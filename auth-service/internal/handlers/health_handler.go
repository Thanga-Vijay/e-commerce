package handlers

import (
	"net/http"

	"github.com/ecommerce/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", gin.H{
		"service": "auth-service",
		"status":  "up",
	})
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check database connection, redis, etc.
	// For now, simple implementation
	utils.SuccessResponse(c, http.StatusOK, "Service is ready", gin.H{
		"service": "auth-service",
		"status":  "ready",
	})
}
