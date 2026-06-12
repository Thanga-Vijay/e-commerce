package handlers

import (
	"errors"
	"net/http"

	"github.com/ecommerce/wishlist-service/internal/models"
	"github.com/ecommerce/wishlist-service/internal/service"
	"github.com/ecommerce/wishlist-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	wishlistService service.WishlistService
}

func NewWishlistHandler(wishlistService service.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// GetWishlist godoc
// @Summary Get user wishlist
// @Description Get current user's wishlist
// @Tags wishlist
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.WishlistResponse}
// @Failure 401 {object} utils.Response
// @Router /wishlist [get]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	wishlist, err := h.wishlistService.GetWishlist(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get wishlist", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Wishlist retrieved successfully", models.ToWishlistResponse(wishlist))
}

// AddToWishlist godoc
// @Summary Add item to wishlist
// @Description Add a product to the wishlist
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddToWishlistRequest true "Add to wishlist request"
// @Success 200 {object} utils.Response{data=models.WishlistResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /wishlist/items [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	var req models.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	wishlist, err := h.wishlistService.AddToWishlist(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "PRODUCT_NOT_FOUND", nil)
			return
		}
		if errors.Is(err, service.ErrItemExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), "ITEM_EXISTS", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add item to wishlist", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item added to wishlist successfully", models.ToWishlistResponse(wishlist))
}

// RemoveFromWishlist godoc
// @Summary Remove item from wishlist
// @Description Remove an item from the wishlist
// @Tags wishlist
// @Security BearerAuth
// @Param id path string true "Wishlist Item ID"
// @Success 200 {object} utils.Response{data=models.WishlistResponse}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /wishlist/items/{id} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
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

	wishlist, err := h.wishlistService.RemoveFromWishlist(c.Request.Context(), userID, itemID)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "ITEM_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove item from wishlist", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item removed from wishlist successfully", models.ToWishlistResponse(wishlist))
}

// ClearWishlist godoc
// @Summary Clear wishlist
// @Description Remove all items from the wishlist
// @Tags wishlist
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /wishlist [delete]
func (h *WishlistHandler) ClearWishlist(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	if err := h.wishlistService.ClearWishlist(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear wishlist", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Wishlist cleared successfully", nil)
}

// MoveToCart godoc
// @Summary Move item from wishlist to cart
// @Description Move a wishlist item to the shopping cart
// @Tags wishlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wishlist Item ID"
// @Param request body models.MoveToCartRequest true "Move to cart request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /wishlist/items/{id}/move-to-cart [post]
func (h *WishlistHandler) MoveToCart(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID", err.Error())
		return
	}

	token, _ := c.Get("token")
	tokenStr := token.(string)

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", "INVALID_ID", err.Error())
		return
	}

	var req models.MoveToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.wishlistService.MoveToCart(c.Request.Context(), userID, itemID, tokenStr, &req); err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), "ITEM_NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to move item to cart", "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Item moved to cart successfully", nil)
}
