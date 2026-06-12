package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reporting-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type ServiceClient struct {
	authServiceURL    string
	productServiceURL string
	orderServiceURL   string
	paymentServiceURL string
	httpClient        *http.Client
}

func NewServiceClient(authURL, productURL, orderURL, paymentURL string) *ServiceClient {
	return &ServiceClient{
		authServiceURL:    authURL,
		productServiceURL: productURL,
		orderServiceURL:   orderURL,
		paymentServiceURL: paymentURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type OrdersResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    OrdersData       `json:"data"`
}

type OrdersData struct {
	Orders     []models.Order   `json:"orders"`
	Pagination PaginationInfo   `json:"pagination"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type PaymentsResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []models.Payment `json:"data"`
}

type UsersResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    UsersData    `json:"data"`
}

type UsersData struct {
	Users      []models.User  `json:"users"`
	Pagination PaginationInfo `json:"pagination"`
}

// GetOrders fetches all orders from order service
func (c *ServiceClient) GetOrders(token string, page, limit int) ([]models.Order, int64, error) {
	url := fmt.Sprintf("%s/api/v1/orders?page=%d&limit=%d", c.orderServiceURL, page, limit)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("failed to fetch orders: status %d, body: %s", resp.StatusCode, string(body))
	}
	
	var response OrdersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return response.Data.Orders, response.Data.Pagination.Total, nil
}

// GetOrdersByDateRange fetches orders within a date range
func (c *ServiceClient) GetOrdersByDateRange(token string, startDate, endDate time.Time) ([]models.Order, error) {
	// Note: This assumes the order service has a date range filter endpoint
	// If not, we fetch all orders and filter in memory
	var allOrders []models.Order
	page := 1
	limit := 100
	
	for {
		orders, total, err := c.GetOrders(token, page, limit)
		if err != nil {
			return nil, err
		}
		
		// Filter orders by date range
		for _, order := range orders {
			if order.CreatedAt.After(startDate) && order.CreatedAt.Before(endDate) {
				allOrders = append(allOrders, order)
			}
		}
		
		if int64(page*limit) >= total {
			break
		}
		page++
	}
	
	return allOrders, nil
}

// GetPaymentsByOrderIDs fetches payments for specific orders
func (c *ServiceClient) GetPaymentsByOrderIDs(token string, orderIDs []uuid.UUID) (map[uuid.UUID]models.Payment, error) {
	payments := make(map[uuid.UUID]models.Payment)
	
	for _, orderID := range orderIDs {
		url := fmt.Sprintf("%s/api/v1/payments/order/%s", c.paymentServiceURL, orderID.String())
		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		
		resp, err := c.httpClient.Do(req)
		if err != nil {
			continue
		}
		
		if resp.StatusCode == http.StatusOK {
			var payment models.Payment
			if err := json.NewDecoder(resp.Body).Decode(&payment); err == nil {
				payments[orderID] = payment
			}
		}
		resp.Body.Close()
	}
	
	return payments, nil
}

// GetUsers fetches all users from auth service
func (c *ServiceClient) GetUsers(token string, page, limit int) ([]models.User, int64, error) {
	url := fmt.Sprintf("%s/api/v1/users?page=%d&limit=%d", c.authServiceURL, page, limit)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// If the auth service doesn't have a users list endpoint, return empty
		return []models.User{}, 0, nil
	}
	
	var response UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return response.Data.Users, response.Data.Pagination.Total, nil
}
