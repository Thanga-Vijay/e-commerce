package events

import "time"

// Event Types
const (
	// User Events
	EventUserRegistered = "user.registered"
	EventUserUpdated    = "user.updated"
	EventUserDeleted    = "user.deleted"

	// Product Events
	EventProductCreated = "product.created"
	EventProductUpdated = "product.updated"
	EventProductDeleted = "product.deleted"

	// Cart Events
	EventCartItemAdded   = "cart.item.added"
	EventCartItemRemoved = "cart.item.removed"
	EventCartCleared     = "cart.cleared"

	// Order Events
	EventOrderCreated   = "order.created"
	EventOrderConfirmed = "order.confirmed"
	EventOrderShipped   = "order.shipped"
	EventOrderDelivered = "order.delivered"
	EventOrderCancelled = "order.cancelled"

	// Payment Events
	EventPaymentInitiated = "payment.initiated"
	EventPaymentCompleted = "payment.completed"
	EventPaymentFailed    = "payment.failed"
	EventPaymentRefunded  = "payment.refunded"

	// Inventory Events
	EventInventoryUpdated   = "inventory.updated"
	EventInventoryLowStock  = "inventory.low-stock"
	EventInventoryOutOfStock = "inventory.out-of-stock"

	// Notification Events
	EventNotificationEmailSend   = "notification.email.send"
	EventNotificationEmailSent   = "notification.email.sent"
	EventNotificationEmailFailed = "notification.email.failed"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// User Events

type UserRegisteredEvent struct {
	BaseEvent
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserUpdatedEvent struct {
	BaseEvent
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserDeletedEvent struct {
	BaseEvent
	UserID uint `json:"user_id"`
}

// Product Events

type ProductCreatedEvent struct {
	BaseEvent
	ProductID   uint    `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  uint    `json:"category_id"`
	Stock       int     `json:"stock"`
}

type ProductUpdatedEvent struct {
	BaseEvent
	ProductID   uint    `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  uint    `json:"category_id"`
	Stock       int     `json:"stock"`
}

type ProductDeletedEvent struct {
	BaseEvent
	ProductID uint `json:"product_id"`
}

// Cart Events

type CartItemAddedEvent struct {
	BaseEvent
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CartItemRemovedEvent struct {
	BaseEvent
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
}

type CartClearedEvent struct {
	BaseEvent
	UserID uint `json:"user_id"`
}

// Order Events

type OrderCreatedEvent struct {
	BaseEvent
	OrderID      uint    `json:"order_id"`
	UserID       uint    `json:"user_id"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
	ItemCount    int     `json:"item_count"`
	ShippingAddr string  `json:"shipping_address"`
}

type OrderConfirmedEvent struct {
	BaseEvent
	OrderID     uint   `json:"order_id"`
	UserID      uint   `json:"user_id"`
	PaymentID   string `json:"payment_id"`
	TotalAmount float64 `json:"total_amount"`
}

type OrderShippedEvent struct {
	BaseEvent
	OrderID        uint   `json:"order_id"`
	UserID         uint   `json:"user_id"`
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"`
}

type OrderDeliveredEvent struct {
	BaseEvent
	OrderID     uint      `json:"order_id"`
	UserID      uint      `json:"user_id"`
	DeliveredAt time.Time `json:"delivered_at"`
}

type OrderCancelledEvent struct {
	BaseEvent
	OrderID      uint   `json:"order_id"`
	UserID       uint   `json:"user_id"`
	Reason       string `json:"reason"`
	RefundAmount float64 `json:"refund_amount"`
}

// Payment Events

type PaymentInitiatedEvent struct {
	BaseEvent
	PaymentID string  `json:"payment_id"`
	OrderID   uint    `json:"order_id"`
	UserID    uint    `json:"user_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Method    string  `json:"method"`
}

type PaymentCompletedEvent struct {
	BaseEvent
	PaymentID       string  `json:"payment_id"`
	OrderID         uint    `json:"order_id"`
	UserID          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	TransactionID   string  `json:"transaction_id"`
	PaymentProvider string  `json:"payment_provider"`
}

type PaymentFailedEvent struct {
	BaseEvent
	PaymentID string `json:"payment_id"`
	OrderID   uint   `json:"order_id"`
	UserID    uint   `json:"user_id"`
	Reason    string `json:"reason"`
	ErrorCode string `json:"error_code"`
}

type PaymentRefundedEvent struct {
	BaseEvent
	PaymentID     string  `json:"payment_id"`
	OrderID       uint    `json:"order_id"`
	UserID        uint    `json:"user_id"`
	RefundAmount  float64 `json:"refund_amount"`
	RefundReason  string  `json:"refund_reason"`
	TransactionID string  `json:"transaction_id"`
}

// Inventory Events

type InventoryUpdatedEvent struct {
	BaseEvent
	ProductID     uint   `json:"product_id"`
	WarehouseID   uint   `json:"warehouse_id"`
	PreviousStock int    `json:"previous_stock"`
	NewStock      int    `json:"new_stock"`
	ChangeType    string `json:"change_type"` // addition, deduction, adjustment
}

type InventoryLowStockEvent struct {
	BaseEvent
	ProductID   uint `json:"product_id"`
	WarehouseID uint `json:"warehouse_id"`
	CurrentStock int  `json:"current_stock"`
	Threshold   int  `json:"threshold"`
}

type InventoryOutOfStockEvent struct {
	BaseEvent
	ProductID   uint `json:"product_id"`
	WarehouseID uint `json:"warehouse_id"`
}

// Notification Events

type NotificationEmailSendEvent struct {
	BaseEvent
	To       string                 `json:"to"`
	Subject  string                 `json:"subject"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

type NotificationEmailSentEvent struct {
	BaseEvent
	To        string `json:"to"`
	Subject   string `json:"subject"`
	MessageID string `json:"message_id"`
}

type NotificationEmailFailedEvent struct {
	BaseEvent
	To      string `json:"to"`
	Subject string `json:"subject"`
	Error   string `json:"error"`
}

// Dead Letter Queue Event
type DeadLetterEvent struct {
	BaseEvent
	OriginalTopic   string `json:"original_topic"`
	OriginalMessage string `json:"original_message"`
	ErrorMessage    string `json:"error_message"`
	RetryCount      int    `json:"retry_count"`
}
