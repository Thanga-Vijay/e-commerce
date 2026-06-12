package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DashboardMetrics represents dashboard summary metrics
type DashboardMetrics struct {
	TotalRevenue      float64 `json:"totalRevenue"`
	TotalOrders       int64   `json:"totalOrders"`
	TotalCustomers    int64   `json:"totalCustomers"`
	AverageOrderValue float64 `json:"averageOrderValue"`
	TodayRevenue      float64 `json:"todayRevenue"`
	TodayOrders       int64   `json:"todayOrders"`
	MonthRevenue      float64 `json:"monthRevenue"`
	MonthOrders       int64   `json:"monthOrders"`
}

// RevenueReport represents revenue over time
type RevenueReport struct {
	Date    time.Time `json:"date"`
	Revenue float64   `json:"revenue"`
	Orders  int64     `json:"orders"`
}

// TopProduct represents a top-selling product
type TopProduct struct {
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	TotalSold   int64     `json:"totalSold"`
	TotalRevenue float64  `json:"totalRevenue"`
}

// CustomerReport represents customer analytics
type CustomerReport struct {
	CustomerID    uuid.UUID `json:"customerId"`
	CustomerEmail string    `json:"customerEmail"`
	TotalOrders   int64     `json:"totalOrders"`
	TotalSpent    float64   `json:"totalSpent"`
	LastOrderDate time.Time `json:"lastOrderDate"`
}

// Report represents a saved report
type Report struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Type        string         `gorm:"not null" json:"type"` // dashboard, revenue, products, customers
	Period      string         `json:"period"`                // daily, weekly, monthly, yearly
	StartDate   *time.Time     `json:"startDate"`
	EndDate     *time.Time     `json:"endDate"`
	Data        string         `gorm:"type:jsonb" json:"data"` // JSON data of the report
	GeneratedBy uuid.UUID      `gorm:"type:uuid" json:"generatedBy"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Order represents order data from order service
type Order struct {
	ID          uuid.UUID `json:"id"`
	OrderNumber string    `json:"orderNumber"`
	UserID      uuid.UUID `json:"userId"`
	Status      string    `json:"status"`
	SubTotal    float64   `json:"subTotal"`
	Tax         float64   `json:"tax"`
	ShippingCost float64  `json:"shippingCost"`
	Total       float64   `json:"total"`
	CreatedAt   time.Time `json:"createdAt"`
	Items       []OrderItem `json:"items"`
}

// OrderItem represents an order item
type OrderItem struct {
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Total       float64   `json:"total"`
}

// Payment represents payment data from payment service
type Payment struct {
	ID      uuid.UUID `json:"id"`
	OrderID uuid.UUID `json:"orderId"`
	Amount  float64   `json:"amount"`
	Status  string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// User represents user data from auth service
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
