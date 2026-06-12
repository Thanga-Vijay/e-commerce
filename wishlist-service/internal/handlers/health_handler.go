package handlers

import (
	"net/http"

	"github.com/ecommerce/wishlist-service/internal/utils"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", gin.H{
		"service": "wishlist-service",
		"status":  "up",
	})
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Service is ready", gin.H{
		"service": "wishlist-service",
		"status":  "ready",
	})
}
