package service

import (
	"errors"
	"fmt"
	"inventory-service/internal/models"
	"inventory-service/internal/repository"

	"github.com/google/uuid"
)

type CreateInventoryRequest struct {
	ProductID         uuid.UUID `json:"productId" binding:"required"`
	SKU               string    `json:"sku" binding:"required"`
	QuantityAvailable int       `json:"quantityAvailable" binding:"min=0"`
	WarehouseID       uuid.UUID `json:"warehouseId"`
	LowStockThreshold int       `json:"lowStockThreshold" binding:"min=0"`
}

type UpdateInventoryRequest struct {
	QuantityAvailable *int       `json:"quantityAvailable" binding:"omitempty,min=0"`
	WarehouseID       *uuid.UUID `json:"warehouseId"`
	LowStockThreshold *int       `json:"lowStockThreshold" binding:"omitempty,min=0"`
}

type ReserveStockRequest struct {
	Quantity    int        `json:"quantity" binding:"required,gt=0"`
	ReferenceID *uuid.UUID `json:"referenceId"`
	Notes       string     `json:"notes"`
}

type ReleaseStockRequest struct {
	Quantity    int        `json:"quantity" binding:"required,gt=0"`
	ReferenceID *uuid.UUID `json:"referenceId"`
	Notes       string     `json:"notes"`
}

type ConfirmSaleRequest struct {
	Quantity    int        `json:"quantity" binding:"required,gt=0"`
	ReferenceID *uuid.UUID `json:"referenceId"`
	Notes       string     `json:"notes"`
}

type AdjustStockRequest struct {
	Quantity int    `json:"quantity" binding:"required"`
	Notes    string `json:"notes" binding:"required"`
}

type InventoryService interface {
	CreateInventory(req CreateInventoryRequest) (*models.Inventory, error)
	GetInventoryByID(id uuid.UUID) (*models.Inventory, error)
	GetInventoryByProductID(productID uuid.UUID) (*models.Inventory, error)
	GetInventoryBySKU(sku string) (*models.Inventory, error)
	GetAllInventory(page, limit int) ([]models.Inventory, int64, error)
	GetLowStock() ([]models.Inventory, error)
	UpdateInventory(id uuid.UUID, req UpdateInventoryRequest) (*models.Inventory, error)
	ReserveStock(productID uuid.UUID, req ReserveStockRequest) (*models.Inventory, error)
	ReleaseStock(productID uuid.UUID, req ReleaseStockRequest) (*models.Inventory, error)
	ConfirmSale(productID uuid.UUID, req ConfirmSaleRequest) (*models.Inventory, error)
	AdjustStock(productID uuid.UUID, req AdjustStockRequest) (*models.Inventory, error)
	GetTransactionHistory(productID uuid.UUID, limit int) ([]models.InventoryTransaction, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CreateInventory(req CreateInventoryRequest) (*models.Inventory, error) {
	// Check if inventory already exists for this product
	existing, _ := s.repo.FindByProductID(req.ProductID)
	if existing != nil {
		return nil, errors.New("inventory already exists for this product")
	}

	// Validate warehouse if provided
	if req.WarehouseID != uuid.Nil {
		_, err := s.repo.FindWarehouseByID(req.WarehouseID)
		if err != nil {
			return nil, fmt.Errorf("invalid warehouse: %w", err)
		}
	}

	inventory := &models.Inventory{
		ProductID:         req.ProductID,
		SKU:               req.SKU,
		QuantityAvailable: req.QuantityAvailable,
		QuantityReserved:  0,
		QuantitySold:      0,
		WarehouseID:       req.WarehouseID,
		LowStockThreshold: req.LowStockThreshold,
	}

	if inventory.LowStockThreshold == 0 {
		inventory.LowStockThreshold = 10 // Default threshold
	}

	if err := s.repo.Create(inventory); err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	// Create initial transaction
	transaction := &models.InventoryTransaction{
		InventoryID:     inventory.ID,
		TransactionType: models.TransactionTypePurchase,
		Quantity:        req.QuantityAvailable,
		PreviousQty:     0,
		NewQty:          req.QuantityAvailable,
		Notes:           "Initial inventory",
	}
	s.repo.CreateTransaction(transaction)

	return s.repo.FindByID(inventory.ID)
}

func (s *inventoryService) GetInventoryByID(id uuid.UUID) (*models.Inventory, error) {
	return s.repo.FindByID(id)
}

func (s *inventoryService) GetInventoryByProductID(productID uuid.UUID) (*models.Inventory, error) {
	return s.repo.FindByProductID(productID)
}

func (s *inventoryService) GetInventoryBySKU(sku string) (*models.Inventory, error) {
	return s.repo.FindBySKU(sku)
}

func (s *inventoryService) GetAllInventory(page, limit int) ([]models.Inventory, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit)
}

func (s *inventoryService) GetLowStock() ([]models.Inventory, error) {
	return s.repo.FindLowStock()
}

func (s *inventoryService) UpdateInventory(id uuid.UUID, req UpdateInventoryRequest) (*models.Inventory, error) {
	inventory, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.QuantityAvailable != nil {
		inventory.QuantityAvailable = *req.QuantityAvailable
	}

	if req.WarehouseID != nil && *req.WarehouseID != uuid.Nil {
		_, err := s.repo.FindWarehouseByID(*req.WarehouseID)
		if err != nil {
			return nil, fmt.Errorf("invalid warehouse: %w", err)
		}
		inventory.WarehouseID = *req.WarehouseID
	}

	if req.LowStockThreshold != nil {
		inventory.LowStockThreshold = *req.LowStockThreshold
	}

	if err := s.repo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	return s.repo.FindByID(id)
}

func (s *inventoryService) ReserveStock(productID uuid.UUID, req ReserveStockRequest) (*models.Inventory, error) {
	inventory, err := s.repo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	// Check if enough stock is available
	if !inventory.CanReserve(req.Quantity) {
		return nil, fmt.Errorf("insufficient stock: available=%d, requested=%d", inventory.QuantityAvailable, req.Quantity)
	}

	previousQty := inventory.QuantityAvailable
	inventory.QuantityAvailable -= req.Quantity
	inventory.QuantityReserved += req.Quantity

	if err := s.repo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		InventoryID:     inventory.ID,
		TransactionType: models.TransactionTypeReserve,
		Quantity:        req.Quantity,
		PreviousQty:     previousQty,
		NewQty:          inventory.QuantityAvailable,
		ReferenceID:     req.ReferenceID,
		Notes:           req.Notes,
	}
	s.repo.CreateTransaction(transaction)

	return s.repo.FindByID(inventory.ID)
}

func (s *inventoryService) ReleaseStock(productID uuid.UUID, req ReleaseStockRequest) (*models.Inventory, error) {
	inventory, err := s.repo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	// Check if enough reserved stock
	if inventory.QuantityReserved < req.Quantity {
		return nil, fmt.Errorf("insufficient reserved stock: reserved=%d, requested=%d", inventory.QuantityReserved, req.Quantity)
	}

	previousQty := inventory.QuantityAvailable
	inventory.QuantityAvailable += req.Quantity
	inventory.QuantityReserved -= req.Quantity

	if err := s.repo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to release stock: %w", err)
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		InventoryID:     inventory.ID,
		TransactionType: models.TransactionTypeRelease,
		Quantity:        req.Quantity,
		PreviousQty:     previousQty,
		NewQty:          inventory.QuantityAvailable,
		ReferenceID:     req.ReferenceID,
		Notes:           req.Notes,
	}
	s.repo.CreateTransaction(transaction)

	return s.repo.FindByID(inventory.ID)
}

func (s *inventoryService) ConfirmSale(productID uuid.UUID, req ConfirmSaleRequest) (*models.Inventory, error) {
	inventory, err := s.repo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	// Check if enough reserved stock
	if inventory.QuantityReserved < req.Quantity {
		return nil, fmt.Errorf("insufficient reserved stock: reserved=%d, requested=%d", inventory.QuantityReserved, req.Quantity)
	}

	previousQty := inventory.QuantityReserved
	inventory.QuantityReserved -= req.Quantity
	inventory.QuantitySold += req.Quantity

	if err := s.repo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to confirm sale: %w", err)
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		InventoryID:     inventory.ID,
		TransactionType: models.TransactionTypeSale,
		Quantity:        req.Quantity,
		PreviousQty:     previousQty,
		NewQty:          inventory.QuantityReserved,
		ReferenceID:     req.ReferenceID,
		Notes:           req.Notes,
	}
	s.repo.CreateTransaction(transaction)

	return s.repo.FindByID(inventory.ID)
}

func (s *inventoryService) AdjustStock(productID uuid.UUID, req AdjustStockRequest) (*models.Inventory, error) {
	inventory, err := s.repo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	previousQty := inventory.QuantityAvailable
	inventory.QuantityAvailable += req.Quantity

	// Prevent negative inventory
	if inventory.QuantityAvailable < 0 {
		return nil, errors.New("adjustment would result in negative inventory")
	}

	if err := s.repo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to adjust stock: %w", err)
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		InventoryID:     inventory.ID,
		TransactionType: models.TransactionTypeAdjustment,
		Quantity:        req.Quantity,
		PreviousQty:     previousQty,
		NewQty:          inventory.QuantityAvailable,
		Notes:           req.Notes,
	}
	s.repo.CreateTransaction(transaction)

	return s.repo.FindByID(inventory.ID)
}

func (s *inventoryService) GetTransactionHistory(productID uuid.UUID, limit int) ([]models.InventoryTransaction, error) {
	inventory, err := s.repo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	return s.repo.FindTransactionsByInventoryID(inventory.ID, limit)
}
