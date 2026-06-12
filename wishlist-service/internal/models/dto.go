package models

// Request DTOs
type AddToWishlistRequest struct {
	ProductID string `json:"productId" binding:"required,uuid"`
}

type MoveToCartRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}

// Response DTOs
type WishlistResponse struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"userId"`
	Items     []WishlistItemResponse `json:"items"`
	ItemCount int                    `json:"itemCount"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

type WishlistItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"productId"`
	ProductName  string  `json:"productName"`
	ProductPrice float64 `json:"productPrice"`
	ProductImage string  `json:"productImage"`
	AddedAt      string  `json:"addedAt"`
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
func ToWishlistResponse(wishlist *Wishlist) *WishlistResponse {
	items := make([]WishlistItemResponse, len(wishlist.Items))
	
	for i, item := range wishlist.Items {
		items[i] = WishlistItemResponse{
			ID:           item.ID.String(),
			ProductID:    item.ProductID.String(),
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			ProductImage: item.ProductImage,
			AddedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &WishlistResponse{
		ID:        wishlist.ID.String(),
		UserID:    wishlist.UserID.String(),
		Items:     items,
		ItemCount: len(wishlist.Items),
		CreatedAt: wishlist.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: wishlist.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
