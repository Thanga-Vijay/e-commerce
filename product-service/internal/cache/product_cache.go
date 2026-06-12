package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ecommerce/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ProductCache interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	SetProduct(ctx context.Context, product *models.Product, ttl time.Duration) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	SetCategory(ctx context.Context, category *models.Category, ttl time.Duration) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetSearchResults(ctx context.Context, key string) ([]models.Product, error)
	SetSearchResults(ctx context.Context, key string, products []models.Product, ttl time.Duration) error
	InvalidateSearchCache(ctx context.Context) error
}

type productCache struct {
	client *redis.Client
}

func NewProductCache(client *redis.Client) ProductCache {
	return &productCache{client: client}
}

func (c *productCache) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	key := fmt.Sprintf("product:%s", id.String())
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *productCache) SetProduct(ctx context.Context, product *models.Product, ttl time.Duration) error {
	key := fmt.Sprintf("product:%s", product.ID.String())
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *productCache) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("product:%s", id.String())
	return c.client.Del(ctx, key).Err()
}

func (c *productCache) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	key := fmt.Sprintf("category:%s", id.String())
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var category models.Category
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *productCache) SetCategory(ctx context.Context, category *models.Category, ttl time.Duration) error {
	key := fmt.Sprintf("category:%s", category.ID.String())
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *productCache) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("category:%s", id.String())
	return c.client.Del(ctx, key).Err()
}

func (c *productCache) GetSearchResults(ctx context.Context, key string) ([]models.Product, error) {
	cacheKey := fmt.Sprintf("search:%s", key)
	data, err := c.client.Get(ctx, cacheKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (c *productCache) SetSearchResults(ctx context.Context, key string, products []models.Product, ttl time.Duration) error {
	cacheKey := fmt.Sprintf("search:%s", key)
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, cacheKey, data, ttl).Err()
}

func (c *productCache) InvalidateSearchCache(ctx context.Context) error {
	// Delete all search cache keys
	iter := c.client.Scan(ctx, 0, "search:*", 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
