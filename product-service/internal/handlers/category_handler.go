package handlers

import (
	"errors"
	"net/http"

	"github.com/ecommerce/product-service/internal/models"
	"github.com/ecommerce/product-service/internal/service"
	"github.com/ecommerce/product-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), "CATEGORY_EXISTS", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", models.ToCategoryResponse(category))
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", "INVALID_ID", err.Error())
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "CATEGORY_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get category", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", models.ToCategoryResponse(category))
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetCategories(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get categories", "INTERNAL_ERROR", err.Error())
		return
	}

	// Convert to response
	categoryResponses := make([]models.CategoryResponse, len(categories))
	for i, cat := range categories {
		categoryResponses[i] = *models.ToCategoryResponse(&cat)
	}

	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categoryResponses)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", "INVALID_ID", err.Error())
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "CATEGORY_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", models.ToCategoryResponse(category))
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", "INVALID_ID", err.Error())
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
