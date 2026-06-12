package handlers

import (
	"errors"
	"net/http"

	"github.com/ecommerce/cart-service/internal/models"
	"github.com/ecommerce/cart-service/internal/service"
	"github.com/ecommerce/cart-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart godoc
// @Summary Get user cart
// @Description Get current user's shopping cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.CartResponse}
// @Failure 401 {object} utils.Response
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get cart", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart retrieved successfully", models.ToCartResponse(cart))
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product to the shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddToCartRequest true "Add to cart request"
// @Success 200 {object} utils.Response{data=models.CartResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	cart, err := h.cartService.AddToCart(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "PRODUCT_NOT_FOUND", nil)
			return
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "INSUFFICIENT_STOCK", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add item to cart", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item added to cart successfully", models.ToCartResponse(cart))
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Description Update the quantity of an item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Param request body models.UpdateCartItemRequest true "Update cart item request"
// @Success 200 {object} utils.Response{data=models.CartResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /cart/items/{id} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", "INVALID_ID", err.Error())
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	cart, err := h.cartService.UpdateCartItem(c.Request.Context(), userID, itemID, &req)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "ITEM_NOT_FOUND", nil)
			return
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), "INSUFFICIENT_STOCK", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update cart item", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart item updated successfully", models.ToCartResponse(cart))
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Description Remove an item from the shopping cart
// @Tags cart
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Success 200 {object} utils.Response{data=models.CartResponse}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /cart/items/{id} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", "INVALID_ID", err.Error())
		return
	}

	cart, err := h.cartService.RemoveFromCart(c.Request.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "ITEM_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove item from cart", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item removed from cart successfully", models.ToCartResponse(cart))
}

// ClearCart godoc
// @Summary Clear cart
// @Description Remove all items from the shopping cart
// @Tags cart
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear cart", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cart cleared successfully", nil)
}
