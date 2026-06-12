package handlers

import (
	"errors"
	"net/http"

	"github.com/ecommerce/auth-service/internal/models"
	"github.com/ecommerce/auth-service/internal/service"
	"github.com/ecommerce/auth-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} utils.Response{data=models.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), "USER_EXISTS", nil)
			return
		}
		if errors.Is(err, service.ErrWeakPassword) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "WEAK_PASSWORD", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=models.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "INVALID_CREDENTIALS", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} utils.Response{data=models.TokenResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "INVALID_TOKEN", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID, req.RefreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to logout", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "USER_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get profile", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", models.ToUserResponse(user))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "USER_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", models.ToUserResponse(user))
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset token to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	token, err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process request", "INTERNAL_ERROR", err.Error())
		return
	}

	// In production, send email with token
	// For development, return token in response
	utils.SuccessResponse(c, http.StatusOK, "Password reset instructions sent to email", gin.H{"token": token})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), "INVALID_TOKEN", nil)
			return
		}
		if errors.Is(err, service.ErrWeakPassword) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "WEAK_PASSWORD", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reset password", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}
