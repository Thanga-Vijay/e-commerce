package repository

import (
	"errors"
	"inventory-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	// Inventory operations
	Create(inventory *models.Inventory) error
	FindByID(id uuid.UUID) (*models.Inventory, error)
	FindByProductID(productID uuid.UUID) (*models.Inventory, error)
	FindBySKU(sku string) (*models.Inventory, error)
	FindAll(page, limit int) ([]models.Inventory, int64, error)
	FindLowStock() ([]models.Inventory, error)
	Update(inventory *models.Inventory) error
	
	// Warehouse operations
	CreateWarehouse(warehouse *models.Warehouse) error
	FindWarehouseByID(id uuid.UUID) (*models.Warehouse, error)
	FindWarehouseByCode(code string) (*models.Warehouse, error)
	FindAllWarehouses() ([]models.Warehouse, error)
	
	// Transaction operations
	CreateTransaction(transaction *models.InventoryTransaction) error
	FindTransactionsByInventoryID(inventoryID uuid.UUID, limit int) ([]models.InventoryTransaction, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Inventory operations
func (r *inventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *inventoryRepository) FindByID(id uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Preload("Warehouse").First(&inventory, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory not found")
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) FindByProductID(productID uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Preload("Warehouse").First(&inventory, "product_id = ?", productID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory not found")
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) FindBySKU(sku string) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Preload("Warehouse").First(&inventory, "sku = ?", sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory not found")
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) FindAll(page, limit int) ([]models.Inventory, int64, error) {
	var inventory []models.Inventory
	var total int64

	query := r.db.Model(&models.Inventory{})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate and fetch
	offset := (page - 1) * limit
	err := query.
		Preload("Warehouse").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&inventory).Error

	if err != nil {
		return nil, 0, err
	}

	return inventory, total, nil
}

func (r *inventoryRepository) FindLowStock() ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.db.
		Preload("Warehouse").
		Where("quantity_available <= low_stock_threshold").
		Find(&inventory).Error
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *inventoryRepository) Update(inventory *models.Inventory) error {
	return r.db.Save(inventory).Error
}

// Warehouse operations
func (r *inventoryRepository) CreateWarehouse(warehouse *models.Warehouse) error {
	return r.db.Create(warehouse).Error
}

func (r *inventoryRepository) FindWarehouseByID(id uuid.UUID) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.First(&warehouse, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("warehouse not found")
		}
		return nil, err
	}
	return &warehouse, nil
}

func (r *inventoryRepository) FindWarehouseByCode(code string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.First(&warehouse, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("warehouse not found")
		}
		return nil, err
	}
	return &warehouse, nil
}

func (r *inventoryRepository) FindAllWarehouses() ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	err := r.db.Where("is_active = ?", true).Find(&warehouses).Error
	if err != nil {
		return nil, err
	}
	return warehouses, nil
}

// Transaction operations
func (r *inventoryRepository) CreateTransaction(transaction *models.InventoryTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *inventoryRepository) FindTransactionsByInventoryID(inventoryID uuid.UUID, limit int) ([]models.InventoryTransaction, error) {
	var transactions []models.InventoryTransaction
	query := r.db.Where("inventory_id = ?", inventoryID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
