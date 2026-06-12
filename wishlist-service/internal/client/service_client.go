package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecommerce/wishlist-service/internal/models"
	"github.com/google/uuid"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (*models.ProductResponse, error)
}

type CartClient interface {
	AddToCart(ctx context.Context, token string, productID uuid.UUID, quantity int) error
}

type productClient struct {
	baseURL string
	client  *http.Client
}

type cartClient struct {
	baseURL string
	client  *http.Client
}

func NewProductClient(baseURL string) ProductClient {
	return &productClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func NewCartClient(baseURL string) CartClient {
	return &cartClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *productClient) GetProduct(ctx context.Context, productID uuid.UUID) (*models.ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, productID.String())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get product: status code %d", resp.StatusCode)
	}

	var response struct {
		Data models.ProductResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *cartClient) AddToCart(ctx context.Context, token string, productID uuid.UUID, quantity int) error {
	url := fmt.Sprintf("%s/api/v1/cart/items", c.baseURL)
	
	requestBody := map[string]interface{}{
		"productId": productID.String(),
		"quantity":  quantity,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add to cart: status code %d", resp.StatusCode)
	}

	return nil
}
