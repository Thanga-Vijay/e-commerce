package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order statuses
const (
	OrderStatusPending    = "pending"
	OrderStatusConfirmed  = "confirmed"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
)

// Address represents a shipping or billing address
type Address struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
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

// Order represents a customer order
type Order struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID       `gorm:"type:uuid;not null;index" json:"userId"`
	OrderNumber     string          `gorm:"type:varchar(50);unique;not null;index" json:"orderNumber"`
	Status          string          `gorm:"type:varchar(20);not null;index" json:"status"`
	Subtotal        float64         `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Tax             float64         `gorm:"type:decimal(10,2);not null" json:"tax"`
	ShippingCost    float64         `gorm:"type:decimal(10,2);not null" json:"shippingCost"`
	Total           float64         `gorm:"type:decimal(10,2);not null" json:"total"`
	ShippingAddress Address         `gorm:"type:jsonb;not null" json:"shippingAddress"`
	BillingAddress  Address         `gorm:"type:jsonb;not null" json:"billingAddress"`
	Items           []OrderItem     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	StatusHistory   []OrderStatus   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"statusHistory,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"orderId"`
	ProductID    uuid.UUID      `gorm:"type:uuid;not null" json:"productId"`
	ProductName  string         `gorm:"type:varchar(255);not null" json:"productName"`
	ProductPrice float64        `gorm:"type:decimal(10,2);not null" json:"productPrice"`
	Quantity     int            `gorm:"not null;check:quantity > 0" json:"quantity"`
	Subtotal     float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderStatus represents order status history
type OrderStatus struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"orderId"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"`
	Comment   string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time `json:"timestamp"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}

// TableName specifies the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// TableName specifies the table name for OrderStatus
func (OrderStatus) TableName() string {
	return "order_status_history"
}

// ValidStatuses returns all valid order statuses
func ValidStatuses() []string {
	return []string{
		OrderStatusPending,
		OrderStatusConfirmed,
		OrderStatusProcessing,
		OrderStatusShipped,
		OrderStatusDelivered,
		OrderStatusCancelled,
	}
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status string) bool {
	for _, s := range ValidStatuses() {
		if s == status {
			return true
		}
	}
	return false
}

// CanTransitionTo checks if an order can transition to a new status
func (o *Order) CanTransitionTo(newStatus string) bool {
	// Define allowed transitions
	transitions := map[string][]string{
		OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
		OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:    {OrderStatusDelivered},
		OrderStatusDelivered:  {},
		OrderStatusCancelled:  {},
	}

	allowedTransitions, exists := transitions[o.Status]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}
