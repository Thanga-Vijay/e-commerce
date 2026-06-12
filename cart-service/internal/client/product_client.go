package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecommerce/cart-service/internal/models"
	"github.com/google/uuid"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (*models.ProductResponse, error)
}

type productClient struct {
	baseURL string
	client  *http.Client
}

func NewProductClient(baseURL string) ProductClient {
	return &productClient{
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
