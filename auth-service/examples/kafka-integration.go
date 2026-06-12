package handlers

import (
	"context"
	"net/http"
	"time"

	"auth-service/internal/models"
	"auth-service/pkg/events"
	"auth-service/pkg/kafka"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB            *gorm.DB
	KafkaProducer *kafka.Producer
}

// Register handles user registration with Kafka event
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Username string `json:"username" binding:"required,min=3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Username: input.Username,
		Role:     "customer",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Publish user.registered event
	if h.KafkaProducer != nil {
		event := events.UserRegisteredEvent{
			BaseEvent: events.BaseEvent{
				ID:        uuid.New().String(),
				Type:      events.EventUserRegistered,
				Source:    "auth-service",
				Timestamp: time.Now().UTC(),
				Version:   "1.0",
			},
			UserID:   user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.KafkaProducer.PublishEvent(ctx, "user.registered", events.EventUserRegistered, event); err != nil {
			// Log error but don't fail the request
			c.Writer.Header().Add("X-Event-Warning", "Failed to publish event")
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	})
}

// UpdateUser handles user updates with Kafka event
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Username string `json:"username" binding:"omitempty,min=3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Username != "" {
		user.Username = input.Username
	}

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Publish user.updated event
	if h.KafkaProducer != nil {
		event := events.UserUpdatedEvent{
			BaseEvent: events.BaseEvent{
				ID:        uuid.New().String(),
				Type:      events.EventUserUpdated,
				Source:    "auth-service",
				Timestamp: time.Now().UTC(),
				Version:   "1.0",
			},
			UserID:   user.ID,
			Email:    user.Email,
			Username: user.Username,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.KafkaProducer.PublishEvent(ctx, "user.updated", events.EventUserUpdated, event); err != nil {
			c.Writer.Header().Add("X-Event-Warning", "Failed to publish event")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}
