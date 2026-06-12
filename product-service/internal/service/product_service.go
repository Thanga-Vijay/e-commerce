package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ecommerce/product-service/internal/cache"
	"github.com/ecommerce/product-service/internal/models"
	"github.com/ecommerce/product-service/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrReviewNotFound   = errors.New("review not found")
	ErrProductExists    = errors.New("product with this SKU already exists")
	ErrCategoryExists   = errors.New("category with this name already exists")
	ErrInvalidData      = errors.New("invalid data provided")
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProducts(ctx context.Context, filter *models.ProductFilterRequest) ([]models.Product, int64, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

type CategoryService interface {
	CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req *models.UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type ReviewService interface {
	CreateReview(ctx context.Context, userID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error)
	GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetProductReviews(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]models.Review, int64, error)
	GetUserReviews(ctx context.Context, userID uuid.UUID) ([]models.Review, error)
	UpdateReview(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateReviewRequest) (*models.Review, error)
	DeleteReview(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type productService struct {
	productRepo repository.ProductRepository
	cache       cache.ProductCache
	cacheTTL    time.Duration
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	cache        cache.ProductCache
	cacheTTL     time.Duration
}

type reviewService struct {
	reviewRepo  repository.ReviewRepository
	productRepo repository.ProductRepository
	cache       cache.ProductCache
}

func NewProductService(productRepo repository.ProductRepository, cache cache.ProductCache, cacheTTL time.Duration) ProductService {
	return &productService{
		productRepo: productRepo,
		cache:       cache,
		cacheTTL:    cacheTTL,
	}
}

func NewCategoryService(categoryRepo repository.CategoryRepository, cache cache.ProductCache, cacheTTL time.Duration) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		cache:        cache,
		cacheTTL:     cacheTTL,
	}
}

func NewReviewService(reviewRepo repository.ReviewRepository, productRepo repository.ProductRepository, cache cache.ProductCache) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
		cache:       cache,
	}
}

// Product Service Implementation
func (s *productService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	// Check if SKU exists
	existingProduct, err := s.productRepo.FindBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if existingProduct != nil {
		return nil, ErrProductExists
	}

	// Parse category ID
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, ErrInvalidData
	}

	product := &models.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryID:  categoryID,
		Images:      req.Images,
		Stock:       req.Stock,
		IsActive:    true,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Cache the product
	_ = s.cache.SetProduct(ctx, product, s.cacheTTL)

	// Invalidate search cache
	_ = s.cache.InvalidateSearchCache(ctx)

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	// Try cache first
	product, err := s.cache.GetProduct(ctx, id)
	if err == nil && product != nil {
		return product, nil
	}

	// Fetch from database
	product, err = s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Cache the product
	_ = s.cache.SetProduct(ctx, product, s.cacheTTL)

	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, filter *models.ProductFilterRequest) ([]models.Product, int64, error) {
	// Generate cache key from filter
	cacheKey := fmt.Sprintf("filter:%v", filter)

	// Try cache first
	cachedProducts, err := s.cache.GetSearchResults(ctx, cacheKey)
	if err == nil && cachedProducts != nil {
		return cachedProducts, int64(len(cachedProducts)), nil
	}

	// Fetch from database
	products, total, err := s.productRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Cache the results
	_ = s.cache.SetSearchResults(ctx, cacheKey, products, 15*time.Minute)

	return products, total, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, ErrInvalidData
		}
		product.CategoryID = categoryID
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	// Update cache
	_ = s.cache.SetProduct(ctx, product, s.cacheTTL)

	// Invalidate search cache
	_ = s.cache.InvalidateSearchCache(ctx)

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := s.productRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete from cache
	_ = s.cache.DeleteProduct(ctx, id)

	// Invalidate search cache
	_ = s.cache.InvalidateSearchCache(ctx)

	return nil
}

func (s *productService) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	if err := s.productRepo.UpdateStock(ctx, id, quantity); err != nil {
		return err
	}

	// Delete from cache to force refresh
	_ = s.cache.DeleteProduct(ctx, id)

	return nil
}

// Category Service Implementation
func (s *categoryService) CreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error) {
	// Check if category exists
	existingCategory, err := s.categoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingCategory != nil {
		return nil, ErrCategoryExists
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	// Parse parent ID if provided
	if req.ParentID != nil && *req.ParentID != "" {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, ErrInvalidData
		}
		category.ParentID = &parentID
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	// Cache the category
	_ = s.cache.SetCategory(ctx, category, s.cacheTTL)

	return category, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	// Try cache first
	category, err := s.cache.GetCategory(ctx, id)
	if err == nil && category != nil {
		return category, nil
	}

	// Fetch from database
	category, err = s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	// Cache the category
	_ = s.cache.SetCategory(ctx, category, s.cacheTTL)

	return category, nil
}

func (s *categoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req *models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil && *req.ParentID != "" {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, ErrInvalidData
		}
		category.ParentID = &parentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	// Update cache
	_ = s.cache.SetCategory(ctx, category, s.cacheTTL)

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete from cache
	_ = s.cache.DeleteCategory(ctx, id)

	return nil
}

// Review Service Implementation
func (s *reviewService) CreateReview(ctx context.Context, userID uuid.UUID, req *models.CreateReviewRequest) (*models.Review, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, ErrInvalidData
	}

	// Check if product exists
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	review := &models.Review{
		ProductID:  productID,
		UserID:     userID,
		Rating:     req.Rating,
		Title:      req.Title,
		Comment:    req.Comment,
		IsVerified: false,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	// Invalidate product cache
	_ = s.cache.DeleteProduct(ctx, productID)

	return review, nil
}

func (s *reviewService) GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if review == nil {
		return nil, ErrReviewNotFound
	}
	return review, nil
}

func (s *reviewService) GetProductReviews(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.FindByProductID(ctx, productID, page, pageSize)
}

func (s *reviewService) GetUserReviews(ctx context.Context, userID uuid.UUID) ([]models.Review, error) {
	return s.reviewRepo.FindByUserID(ctx, userID)
}

func (s *reviewService) UpdateReview(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateReviewRequest) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if review == nil {
		return nil, ErrReviewNotFound
	}

	// Check ownership
	if review.UserID != userID {
		return nil, errors.New("unauthorized to update this review")
	}

	// Update fields
	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Title != "" {
		review.Title = req.Title
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	// Invalidate product cache
	_ = s.cache.DeleteProduct(ctx, review.ProductID)

	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if review == nil {
		return ErrReviewNotFound
	}

	// Check ownership
	if review.UserID != userID {
		return errors.New("unauthorized to delete this review")
	}

	if err := s.reviewRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate product cache
	_ = s.cache.DeleteProduct(ctx, review.ProductID)

	return nil
}
