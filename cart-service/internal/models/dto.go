package models

// Request DTOs
type AddToCartRequest struct {
	ProductID string `json:"productId" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gte=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}

// Response DTOs
type CartResponse struct {
	ID        string             `json:"id"`
	UserID    string             `json:"userId"`
	Items     []CartItemResponse `json:"items"`
	Subtotal  float64            `json:"subtotal"`
	Tax       float64            `json:"tax"`
	Total     float64            `json:"total"`
	ItemCount int                `json:"itemCount"`
	ExpiresAt string             `json:"expiresAt"`
	CreatedAt string             `json:"createdAt"`
	UpdatedAt string             `json:"updatedAt"`
}

type CartItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"productId"`
	ProductName  string  `json:"productName"`
	ProductPrice float64 `json:"productPrice"`
	ProductImage string  `json:"productImage"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

// Product Service Response
type ProductResponse struct {
	ID          string   `json:"id"`
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock"`
	IsActive    bool     `json:"isActive"`
}

// Helper functions
func ToCartResponse(cart *Cart) *CartResponse {
	items := make([]CartItemResponse, len(cart.Items))
	subtotal := 0.0
	itemCount := 0

	for i, item := range cart.Items {
		items[i] = CartItemResponse{
			ID:           item.ID.String(),
			ProductID:    item.ProductID.String(),
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			ProductImage: item.ProductImage,
			Quantity:     item.Quantity,
			Subtotal:     item.Subtotal,
		}
		subtotal += item.Subtotal
		itemCount += item.Quantity
	}

	tax := subtotal * 0.08 // 8% tax rate
	total := subtotal + tax

	return &CartResponse{
		ID:        cart.ID.String(),
		UserID:    cart.UserID.String(),
		Items:     items,
		Subtotal:  subtotal,
		Tax:       tax,
		Total:     total,
		ItemCount: itemCount,
		ExpiresAt: cart.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt: cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
