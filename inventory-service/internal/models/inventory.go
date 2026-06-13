package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction types
const (
	TransactionTypePurchase   = "purchase"
	TransactionTypeSale       = "sale"
	TransactionTypeAdjustment = "adjustment"
	TransactionTypeReturn     = "return"
	TransactionTypeReserve    = "reserve"
	TransactionTypeRelease    = "release"
)

// Address represents a warehouse address
type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
	Country      string `json:"country"`
}

// Scan implements the Scanner interface for Address
func (a *Address) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// Value implements the Valuer interface for Address
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Warehouse represents a warehouse location
type Warehouse struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Code      string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Address   Address        `gorm:"type:jsonb;not null" json:"address"`
	IsActive  bool           `gorm:"default:true" json:"isActive"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Inventory represents product inventory levels
type Inventory struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID          uuid.UUID      `gorm:"type:uuid;unique;not null;index" json:"productId"`
	SKU                string         `gorm:"type:varchar(100);not null" json:"sku"`
	QuantityAvailable  int            `gorm:"not null;default:0" json:"quantityAvailable"`
	QuantityReserved   int            `gorm:"not null;default:0" json:"quantityReserved"`
	QuantitySold       int            `gorm:"not null;default:0" json:"quantitySold"`
	WarehouseID        uuid.UUID      `gorm:"type:uuid;index" json:"warehouseId"`
	Warehouse          *Warehouse     `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	LowStockThreshold  int            `gorm:"default:10" json:"lowStockThreshold"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// InventoryTransaction represents an inventory change audit trail
type InventoryTransaction struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	InventoryID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"inventoryId"`
	TransactionType string         `gorm:"type:varchar(20);not null" json:"transactionType"`
	Quantity        int            `gorm:"not null" json:"quantity"`
	PreviousQty     int            `gorm:"not null" json:"previousQty"`
	NewQty          int            `gorm:"not null" json:"newQty"`
	ReferenceID     *uuid.UUID     `gorm:"type:uuid" json:"referenceId,omitempty"`
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Warehouse
func (Warehouse) TableName() string {
	return "warehouses"
}

// TableName specifies the table name for Inventory
func (Inventory) TableName() string {
	return "inventory"
}

// TableName specifies the table name for InventoryTransaction
func (InventoryTransaction) TableName() string {
	return "inventory_transactions"
}

// IsLowStock checks if inventory is below threshold
func (i *Inventory) IsLowStock() bool {
	return i.QuantityAvailable <= i.LowStockThreshold
}

// CanReserve checks if enough stock is available to reserve
func (i *Inventory) CanReserve(quantity int) bool {
	return i.QuantityAvailable >= quantity
}
