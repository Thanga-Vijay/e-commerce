package handlers

import (
	"inventory-service/internal/service"
	"inventory-service/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// CreateInventory creates new inventory record (admin only)
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req service.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := h.inventoryService.CreateInventory(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Inventory created successfully", inventory)
}

// GetInventoryByID retrieves inventory by ID
func (h *InventoryHandler) GetInventoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	inventory, err := h.inventoryService.GetInventoryByID(id)
	if err != nil {
		if err.Error() == "inventory not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inventory retrieved successfully", inventory)
}

// GetInventoryByProductID retrieves inventory by product ID
func (h *InventoryHandler) GetInventoryByProductID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	inventory, err := h.inventoryService.GetInventoryByProductID(productID)
	if err != nil {
		if err.Error() == "inventory not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inventory retrieved successfully", inventory)
}

// GetInventoryBySKU retrieves inventory by SKU
func (h *InventoryHandler) GetInventoryBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "SKU is required")
		return
	}

	inventory, err := h.inventoryService.GetInventoryBySKU(sku)
	if err != nil {
		if err.Error() == "inventory not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inventory retrieved successfully", inventory)
}

// GetAllInventory retrieves all inventory with pagination (admin only)
func (h *InventoryHandler) GetAllInventory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	inventory, total, err := h.inventoryService.GetAllInventory(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := gin.H{
		"inventory": inventory,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Inventory retrieved successfully", response)
}

// GetLowStock retrieves low stock items (admin only)
func (h *InventoryHandler) GetLowStock(c *gin.Context) {
	inventory, err := h.inventoryService.GetLowStock()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Low stock items retrieved successfully", inventory)
}

// UpdateInventory updates inventory record (admin only)
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	var req service.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := h.inventoryService.UpdateInventory(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inventory updated successfully", inventory)
}

// ReserveStock reserves stock for an order (admin only)
func (h *InventoryHandler) ReserveStock(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req service.ReserveStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := h.inventoryService.ReserveStock(productID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock reserved successfully", inventory)
}

// ReleaseStock releases reserved stock (admin only)
func (h *InventoryHandler) ReleaseStock(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req service.ReleaseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := h.inventoryService.ReleaseStock(productID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock released successfully", inventory)
}

// ConfirmSale confirms a sale and reduces reserved stock (admin only)
func (h *InventoryHandler) ConfirmSale(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req service.ConfirmSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := h.inventoryService.ConfirmSale(productID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Sale confirmed successfully", inventory)
}

// AdjustStock manually adjusts stock levels (admin only)
func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req service.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	inventory, err := h.inventoryService.AdjustStock(productID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock adjusted successfully", inventory)
}

// GetTransactionHistory retrieves transaction history for a product
func (h *InventoryHandler) GetTransactionHistory(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	transactions, err := h.inventoryService.GetTransactionHistory(productID, limit)
	if err != nil {
		if err.Error() == "inventory not found" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction history retrieved successfully", transactions)
}
