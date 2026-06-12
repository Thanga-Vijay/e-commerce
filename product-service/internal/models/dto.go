package models

// Request DTOs
type CreateProductRequest struct {
	SKU         string   `json:"sku" binding:"required"`
	Name        string   `json:"name" binding:"required,min=3"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	CategoryID  string   `json:"categoryId" binding:"required,uuid"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock" binding:"gte=0"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=3"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"omitempty,gt=0"`
	CategoryID  string   `json:"categoryId" binding:"omitempty,uuid"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock" binding:"omitempty,gte=0"`
	IsActive    *bool    `json:"isActive"`
}

type ProductFilterRequest struct {
	CategoryID string  `form:"categoryId"`
	MinPrice   float64 `form:"minPrice"`
	MaxPrice   float64 `form:"maxPrice"`
	Search     string  `form:"search"`
	Page       int     `form:"page" binding:"gte=1"`
	PageSize   int     `form:"pageSize" binding:"gte=1,lte=100"`
	SortBy     string  `form:"sortBy"`
	SortOrder  string  `form:"sortOrder"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId" binding:"omitempty,uuid"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=2"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId" binding:"omitempty,uuid"`
	IsActive    *bool   `json:"isActive"`
}

type CreateReviewRequest struct {
	ProductID string `json:"productId" binding:"required,uuid"`
	Rating    int    `json:"rating" binding:"required,gte=1,lte=5"`
	Title     string `json:"title"`
	Comment   string `json:"comment" binding:"required,min=10"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"omitempty,gte=1,lte=5"`
	Title   string `json:"title"`
	Comment string `json:"comment" binding:"omitempty,min=10"`
}

// Response DTOs
type ProductResponse struct {
	ID          string           `json:"id"`
	SKU         string           `json:"sku"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	CategoryID  string           `json:"categoryId"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Images      []string         `json:"images"`
	Stock       int              `json:"stock"`
	IsActive    bool             `json:"isActive"`
	AverageRating float64        `json:"averageRating,omitempty"`
	ReviewCount int              `json:"reviewCount,omitempty"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

type CategoryResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	ParentID    string              `json:"parentId,omitempty"`
	Parent      *CategoryResponse   `json:"parent,omitempty"`
	Children    []CategoryResponse  `json:"children,omitempty"`
	IsActive    bool                `json:"isActive"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

type ReviewResponse struct {
	ID         string `json:"id"`
	ProductID  string `json:"productId"`
	UserID     string `json:"userId"`
	Rating     int    `json:"rating"`
	Title      string `json:"title"`
	Comment    string `json:"comment"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type PaginatedProductResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}

// Helper functions
func ToProductResponse(product *Product) *ProductResponse {
	resp := &ProductResponse{
		ID:          product.ID.String(),
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID.String(),
		Images:      product.Images,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if product.Category != nil {
		resp.Category = ToCategoryResponse(product.Category)
	}

	return resp
}

func ToCategoryResponse(category *Category) *CategoryResponse {
	resp := &CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if category.ParentID != nil {
		resp.ParentID = category.ParentID.String()
	}

	if category.Parent != nil {
		resp.Parent = ToCategoryResponse(category.Parent)
	}

	if len(category.Children) > 0 {
		resp.Children = make([]CategoryResponse, len(category.Children))
		for i, child := range category.Children {
			resp.Children[i] = *ToCategoryResponse(&child)
		}
	}

	return resp
}

func ToReviewResponse(review *Review) *ReviewResponse {
	return &ReviewResponse{
		ID:         review.ID.String(),
		ProductID:  review.ProductID.String(),
		UserID:     review.UserID.String(),
		Rating:     review.Rating,
		Title:      review.Title,
		Comment:    review.Comment,
		IsVerified: review.IsVerified,
		CreatedAt:  review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  review.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
