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

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateProductRequest true "Product data"
// @Success 201 {object} utils.Response{data=models.ProductResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrProductExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), "PRODUCT_EXISTS", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", models.ToProductResponse(product))
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response{data=models.ProductResponse}
// @Failure 404 {object} utils.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", "INVALID_ID", err.Error())
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "PRODUCT_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get product", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", models.ToProductResponse(product))
}

// GetProducts godoc
// @Summary Get products list
// @Description Get filtered and paginated products list
// @Tags products
// @Produce json
// @Param categoryId query string false "Category ID filter"
// @Param minPrice query number false "Minimum price filter"
// @Param maxPrice query number false "Maximum price filter"
// @Param search query string false "Search term"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param sortBy query string false "Sort by field" default(created_at)
// @Param sortOrder query string false "Sort order (ASC/DESC)" default(DESC)
// @Success 200 {object} utils.Response{data=models.PaginatedProductResponse}
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter models.ProductFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", "INVALID_REQUEST", err.Error())
		return
	}

	// Set defaults
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	products, total, err := h.productService.GetProducts(c.Request.Context(), &filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get products", "INTERNAL_ERROR", err.Error())
		return
	}

	// Convert to response
	productResponses := make([]models.ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = *models.ToProductResponse(&p)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	response := models.PaginatedProductResponse{
		Products:   productResponses,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", response)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product details (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body models.UpdateProductRequest true "Product update data"
// @Success 200 {object} utils.Response{data=models.ProductResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", "INVALID_ID", err.Error())
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "PRODUCT_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", models.ToProductResponse(product))
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product (Admin only)
// @Tags products
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", "INVALID_ID", err.Error())
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "PRODUCT_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}
