package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecommerce/product-service/internal/models"
	"github.com/ecommerce/product-service/internal/service"
	"github.com/ecommerce/product-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	review, err := h.reviewService.CreateReview(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "PRODUCT_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create review", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Review created successfully", models.ToReviewResponse(review))
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productIDParam := c.Param("productId")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", "INVALID_ID", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	reviews, total, err := h.reviewService.GetProductReviews(c.Request.Context(), productID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get reviews", "INTERNAL_ERROR", err.Error())
		return
	}

	// Convert to response
	reviewResponses := make([]models.ReviewResponse, len(reviews))
	for i, r := range reviews {
		reviewResponses[i] = *models.ToReviewResponse(&r)
	}

	utils.SuccessResponse(c, http.StatusOK, "Reviews retrieved successfully", gin.H{
		"reviews":  reviewResponses,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", "INVALID_ID", err.Error())
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	review, err := h.reviewService.UpdateReview(c.Request.Context(), id, userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "REVIEW_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update review", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Review updated successfully", models.ToReviewResponse(review))
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid review ID", "INVALID_ID", err.Error())
		return
	}

	if err := h.reviewService.DeleteReview(c.Request.Context(), id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete review", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Review deleted successfully", nil)
}
