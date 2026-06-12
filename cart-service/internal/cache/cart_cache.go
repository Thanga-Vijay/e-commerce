package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ecommerce/cart-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CartCache interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	SetCart(ctx context.Context, cart *models.Cart, ttl time.Duration) error
	DeleteCart(ctx context.Context, userID uuid.UUID) error
}

type cartCache struct {
	client *redis.Client
}

func NewCartCache(client *redis.Client) CartCache {
	return &cartCache{client: client}
}

func (c *cartCache) GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	key := fmt.Sprintf("cart:%s", userID.String())
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *cartCache) SetCart(ctx context.Context, cart *models.Cart, ttl time.Duration) error {
	key := fmt.Sprintf("cart:%s", cart.UserID.String())
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *cartCache) DeleteCart(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("cart:%s", userID.String())
	return c.client.Del(ctx, key).Err()
}
