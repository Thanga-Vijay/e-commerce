package service

import (
	"context"
	"fmt"
	"reporting-service/internal/cache"
	"reporting-service/internal/client"
	"reporting-service/internal/models"
	"sort"
	"time"

	"github.com/google/uuid"
)

const (
	dashboardCacheKey = "dashboard:metrics"
	revenueCacheKey   = "revenue:report:%s" // format: start_end
	cacheTTL          = 5 * time.Minute
)

type MetricsService interface {
	GetDashboardMetrics(ctx context.Context, token string) (*models.DashboardMetrics, error)
	GetRevenueReport(ctx context.Context, token string, startDate, endDate time.Time, period string) ([]models.RevenueReport, error)
	GetTopProducts(ctx context.Context, token string, limit int) ([]models.TopProduct, error)
	GetCustomerReports(ctx context.Context, token string, limit int) ([]models.CustomerReport, error)
	InvalidateCache(ctx context.Context) error
}

type metricsService struct {
	client       *client.ServiceClient
	cacheService cache.CacheService
}

func NewMetricsService(client *client.ServiceClient, cacheService cache.CacheService) MetricsService {
	return &metricsService{
		client:       client,
		cacheService: cacheService,
	}
}

func (s *metricsService) GetDashboardMetrics(ctx context.Context, token string) (*models.DashboardMetrics, error) {
	// Try to get from cache
	var metrics models.DashboardMetrics
	if err := s.cacheService.Get(ctx, dashboardCacheKey, &metrics); err == nil {
		return &metrics, nil
	}

	// Calculate metrics from orders
	orders, err := s.fetchAllOrders(token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	metrics = models.DashboardMetrics{}
	todayStart := time.Now().Truncate(24 * time.Hour)
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())

	customerMap := make(map[uuid.UUID]bool)
	var totalRevenue float64
	var todayRevenue float64
	var monthRevenue float64
	var todayOrders int64
	var monthOrders int64

	for _, order := range orders {
		if order.Status != "cancelled" {
			totalRevenue += order.Total
			customerMap[order.UserID] = true

			if order.CreatedAt.After(todayStart) {
				todayRevenue += order.Total
				todayOrders++
			}

			if order.CreatedAt.After(monthStart) {
				monthRevenue += order.Total
				monthOrders++
			}
		}
	}

	metrics.TotalRevenue = totalRevenue
	metrics.TotalOrders = int64(len(orders))
	metrics.TotalCustomers = int64(len(customerMap))
	if len(orders) > 0 {
		metrics.AverageOrderValue = totalRevenue / float64(len(orders))
	}
	metrics.TodayRevenue = todayRevenue
	metrics.TodayOrders = todayOrders
	metrics.MonthRevenue = monthRevenue
	metrics.MonthOrders = monthOrders

	// Cache the result
	s.cacheService.Set(ctx, dashboardCacheKey, metrics, cacheTTL)

	return &metrics, nil
}

func (s *metricsService) GetRevenueReport(ctx context.Context, token string, startDate, endDate time.Time, period string) ([]models.RevenueReport, error) {
	cacheKey := fmt.Sprintf(revenueCacheKey, fmt.Sprintf("%s_%s_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), period))
	
	// Try to get from cache
	var report []models.RevenueReport
	if err := s.cacheService.Get(ctx, cacheKey, &report); err == nil {
		return report, nil
	}

	// Fetch orders in date range
	orders, err := s.client.GetOrdersByDateRange(token, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	// Group orders by period
	revenueMap := make(map[string]*models.RevenueReport)

	for _, order := range orders {
		if order.Status == "cancelled" {
			continue
		}

		var key string
		switch period {
		case "daily":
			key = order.CreatedAt.Format("2006-01-02")
		case "weekly":
			year, week := order.CreatedAt.ISOWeek()
			key = fmt.Sprintf("%d-W%02d", year, week)
		case "monthly":
			key = order.CreatedAt.Format("2006-01")
		case "yearly":
			key = order.CreatedAt.Format("2006")
		default:
			key = order.CreatedAt.Format("2006-01-02")
		}

		if _, exists := revenueMap[key]; !exists {
			date, _ := time.Parse("2006-01-02", key)
			revenueMap[key] = &models.RevenueReport{
				Date:    date,
				Revenue: 0,
				Orders:  0,
			}
		}

		revenueMap[key].Revenue += order.Total
		revenueMap[key].Orders++
	}

	// Convert map to slice
	for _, rev := range revenueMap {
		report = append(report, *rev)
	}

	// Sort by date
	sort.Slice(report, func(i, j int) bool {
		return report[i].Date.Before(report[j].Date)
	})

	// Cache the result
	s.cacheService.Set(ctx, cacheKey, report, cacheTTL)

	return report, nil
}

func (s *metricsService) GetTopProducts(ctx context.Context, token string, limit int) ([]models.TopProduct, error) {
	orders, err := s.fetchAllOrders(token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	// Aggregate product sales
	productMap := make(map[uuid.UUID]*models.TopProduct)

	for _, order := range orders {
		if order.Status == "cancelled" {
			continue
		}

		for _, item := range order.Items {
			if _, exists := productMap[item.ProductID]; !exists {
				productMap[item.ProductID] = &models.TopProduct{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
					TotalSold:   0,
					TotalRevenue: 0,
				}
			}

			productMap[item.ProductID].TotalSold += int64(item.Quantity)
			productMap[item.ProductID].TotalRevenue += item.Total
		}
	}

	// Convert map to slice
	var products []models.TopProduct
	for _, product := range productMap {
		products = append(products, *product)
	}

	// Sort by total revenue
	sort.Slice(products, func(i, j int) bool {
		return products[i].TotalRevenue > products[j].TotalRevenue
	})

	// Limit results
	if limit > 0 && len(products) > limit {
		products = products[:limit]
	}

	return products, nil
}

func (s *metricsService) GetCustomerReports(ctx context.Context, token string, limit int) ([]models.CustomerReport, error) {
	orders, err := s.fetchAllOrders(token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	// Aggregate customer data
	customerMap := make(map[uuid.UUID]*models.CustomerReport)

	for _, order := range orders {
		if order.Status == "cancelled" {
			continue
		}

		if _, exists := customerMap[order.UserID]; !exists {
			customerMap[order.UserID] = &models.CustomerReport{
				CustomerID:    order.UserID,
				CustomerEmail: "", // Will need to fetch from auth service
				TotalOrders:   0,
				TotalSpent:    0,
				LastOrderDate: order.CreatedAt,
			}
		}

		customerMap[order.UserID].TotalOrders++
		customerMap[order.UserID].TotalSpent += order.Total

		if order.CreatedAt.After(customerMap[order.UserID].LastOrderDate) {
			customerMap[order.UserID].LastOrderDate = order.CreatedAt
		}
	}

	// Convert map to slice
	var customers []models.CustomerReport
	for _, customer := range customerMap {
		customers = append(customers, *customer)
	}

	// Sort by total spent
	sort.Slice(customers, func(i, j int) bool {
		return customers[i].TotalSpent > customers[j].TotalSpent
	})

	// Limit results
	if limit > 0 && len(customers) > limit {
		customers = customers[:limit]
	}

	return customers, nil
}

func (s *metricsService) InvalidateCache(ctx context.Context) error {
	if err := s.cacheService.Delete(ctx, dashboardCacheKey); err != nil {
		return err
	}
	return s.cacheService.DeleteByPattern(ctx, "revenue:report:*")
}

func (s *metricsService) fetchAllOrders(token string) ([]models.Order, error) {
	var allOrders []models.Order
	page := 1
	limit := 100

	for {
		orders, total, err := s.client.GetOrders(token, page, limit)
		if err != nil {
			return nil, err
		}

		allOrders = append(allOrders, orders...)

		if int64(page*limit) >= total {
			break
		}
		page++
	}

	return allOrders, nil
}
