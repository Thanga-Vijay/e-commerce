package repository

import (
	"context"
	"errors"

	"github.com/ecommerce/product-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	FindAll(ctx context.Context, filter *models.ProductFilterRequest) ([]models.Product, int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	FindByName(ctx context.Context, name string) (*models.Category, error)
	FindAll(ctx context.Context) ([]models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	FindByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]models.Review, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, int, error)
}

type productRepository struct {
	db *gorm.DB
}

type categoryRepository struct {
	db *gorm.DB
}

type reviewRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// Product Repository Implementation
func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").Where("id = ?", id).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func (r *productRepository) FindAll(ctx context.Context, filter *models.ProductFilterRequest) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{}).Preload("Category")

	// Apply filters
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder != "" {
		sortOrder = filter.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Find(&products).Error
	return products, total, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

// Category Repository Implementation
func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").Where("id = ?", id).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Preload("Children").Where("parent_id IS NULL").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, id).Error
}

// Review Repository Implementation
func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&review).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &review, err
}

func (r *reviewRepository) FindByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Review{}).Where("product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error
	return reviews, total, err
}

func (r *reviewRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int
	}

	err := r.db.WithContext(ctx).Model(&models.Review{}).
		Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("product_id = ?", productID).
		Scan(&result).Error

	return result.AvgRating, result.Count, err
}
