package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"productId"`
	ProductName  string    `json:"productName"`
	ProductPrice float64   `json:"productPrice"`
	Quantity     int       `json:"quantity"`
	Subtotal     float64   `json:"subtotal"`
}

type Cart struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Items     []CartItem `json:"items"`
	Subtotal  float64    `json:"subtotal"`
	Tax       float64    `json:"tax"`
	Total     float64    `json:"total"`
	ItemCount int        `json:"itemCount"`
}

type CartResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Cart   `json:"data"`
}

type CartClient interface {
	GetCart(token string) (*Cart, error)
	ClearCart(token string) error
}

type cartClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCartClient(baseURL string) CartClient {
	return &cartClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *cartClient) GetCart(token string) (*Cart, error) {
	url := fmt.Sprintf("%s/api/v1/cart", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cart service returned status %d: %s", resp.StatusCode, string(body))
	}

	var cartResp CartResponse
	if err := json.Unmarshal(body, &cartResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &cartResp.Data, nil
}

func (c *cartClient) ClearCart(token string) error {
	url := fmt.Sprintf("%s/api/v1/cart", c.baseURL)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cart service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
