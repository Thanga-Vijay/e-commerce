# E-Commerce Platform - Phase 6 Implementation

This document provides instructions for running the Reporting Service (Phase 6).

## Service Overview

**Phase 6** adds the Reporting Service:
- **Reporting Service** (Port 8089): Dashboard metrics, analytics, and report generation with CSV/PDF export

## Quick Start with Docker Compose

Start all services including Phases 2-6:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs for reporting service
docker-compose logs -f reporting-service

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

Services available:
- **Auth Service**: http://localhost:8081
- **Product Service**: http://localhost:8082
- **Cart Service**: http://localhost:8083
- **Wishlist Service**: http://localhost:8084
- **Order Service**: http://localhost:8085
- **Payment Service**: http://localhost:8086
- **Inventory Service**: http://localhost:8087
- **Notification Service**: http://localhost:8088
- **Reporting Service**: http://localhost:8089
- **PostgreSQL Reporting DB**: localhost:5441
- **Redis**: localhost:6379

## Manual Setup

### 1. Set Up Database

```bash
# Create database
createdb reporting_db

# Or using psql
psql -U postgres -c "CREATE DATABASE reporting_db;"
```

### 2. Run Migrations

```bash
cd reporting-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/reporting_db?sslmode=disable" up
```

### 3. Configure Environment

```bash
cd reporting-service
cp .env.example .env
# Edit .env with your configuration
```

### 4. Install Dependencies

```bash
cd reporting-service
go mod download
```

### 5. Run Service

```bash
cd reporting-service
go run cmd/main.go
```

## Reporting Service API

### Dashboard Metrics

**Get Dashboard:**
```bash
GET /api/v1/dashboard
Authorization: Bearer YOUR_ACCESS_TOKEN

Response:
{
  "status": 200,
  "message": "Dashboard metrics retrieved successfully",
  "data": {
    "totalRevenue": 125478.50,
    "totalOrders": 487,
    "totalCustomers": 234,
    "averageOrderValue": 257.65,
    "todayRevenue": 8456.25,
    "todayOrders": 32,
    "monthRevenue": 89234.75,
    "monthOrders": 342
  }
}
```

**Export Dashboard:**
```bash
# Export as CSV
GET /api/v1/dashboard/export?format=csv
Authorization: Bearer YOUR_ACCESS_TOKEN

# Export as PDF
GET /api/v1/dashboard/export?format=pdf
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Revenue Reports

**Get Revenue Report:**
```bash
GET /api/v1/reports/revenue?startDate=2026-01-01&endDate=2026-12-31&period=monthly
Authorization: Bearer YOUR_ACCESS_TOKEN

Query Parameters:
- startDate: YYYY-MM-DD (required)
- endDate: YYYY-MM-DD (required)
- period: daily, weekly, monthly, yearly (default: daily)

Response:
{
  "status": 200,
  "message": "Revenue report retrieved successfully",
  "data": [
    {
      "date": "2026-01-01T00:00:00Z",
      "revenue": 15234.50,
      "orders": 62
    },
    {
      "date": "2026-02-01T00:00:00Z",
      "revenue": 18456.75,
      "orders": 73
    }
  ]
}
```

**Export Revenue Report:**
```bash
# Export as CSV
GET /api/v1/reports/revenue/export?startDate=2026-01-01&endDate=2026-12-31&period=monthly&format=csv
Authorization: Bearer YOUR_ACCESS_TOKEN

# Export as PDF
GET /api/v1/reports/revenue/export?startDate=2026-01-01&endDate=2026-12-31&period=monthly&format=pdf
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Product Reports

**Get Top Products:**
```bash
GET /api/v1/reports/products?limit=10
Authorization: Bearer YOUR_ACCESS_TOKEN

Query Parameters:
- limit: Number of products to return (default: 10)

Response:
{
  "status": 200,
  "message": "Top products retrieved successfully",
  "data": [
    {
      "productId": "550e8400-e29b-41d4-a716-446655440000",
      "productName": "MacBook Pro 16",
      "totalSold": 145,
      "totalRevenue": 362485.55
    },
    {
      "productId": "550e8400-e29b-41d4-a716-446655440001",
      "productName": "iPhone 15 Pro",
      "totalSold": 234,
      "totalRevenue": 234000.66
    }
  ]
}
```

**Export Top Products:**
```bash
# Export as CSV
GET /api/v1/reports/products/export?limit=10&format=csv
Authorization: Bearer YOUR_ACCESS_TOKEN

# Export as PDF
GET /api/v1/reports/products/export?limit=10&format=pdf
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Customer Reports

**Get Customer Report:**
```bash
GET /api/v1/reports/customers?limit=10
Authorization: Bearer YOUR_ACCESS_TOKEN

Query Parameters:
- limit: Number of customers to return (default: 10)

Response:
{
  "status": 200,
  "message": "Customer report retrieved successfully",
  "data": [
    {
      "customerId": "650e8400-e29b-41d4-a716-446655440000",
      "customerEmail": "john@example.com",
      "totalOrders": 28,
      "totalSpent": 14562.45,
      "lastOrderDate": "2026-06-10T14:30:00Z"
    },
    {
      "customerId": "650e8400-e29b-41d4-a716-446655440001",
      "customerEmail": "jane@example.com",
      "totalOrders": 22,
      "totalSpent": 12345.67,
      "lastOrderDate": "2026-06-12T10:15:00Z"
    }
  ]
}
```

**Export Customer Report:**
```bash
# Export as CSV
GET /api/v1/reports/customers/export?limit=10&format=csv
Authorization: Bearer YOUR_ACCESS_TOKEN

# Export as PDF
GET /api/v1/reports/customers/export?limit=10&format=pdf
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Saved Reports

**Save Report:**
```bash
POST /api/v1/reports
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "name": "Q2 2026 Revenue Report",
  "type": "revenue",
  "period": "monthly",
  "startDate": "2026-04-01T00:00:00Z",
  "endDate": "2026-06-30T00:00:00Z",
  "data": {
    "revenue": 54321.98,
    "orders": 245
  }
}
```

**Get All Saved Reports:**
```bash
GET /api/v1/reports?page=1&limit=10
Authorization: Bearer YOUR_ACCESS_TOKEN

Response:
{
  "status": 200,
  "message": "Reports retrieved successfully",
  "data": {
    "reports": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 45,
      "totalPages": 5
    }
  }
}
```

**Get Report by ID:**
```bash
GET /api/v1/reports/:id
Authorization: Bearer YOUR_ACCESS_TOKEN
```

**Delete Report:**
```bash
DELETE /api/v1/reports/:id
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Admin Endpoints

**Invalidate Cache:**
```bash
POST /api/v1/cache/invalidate
Authorization: Bearer ADMIN_TOKEN

Response:
{
  "status": 200,
  "message": "Cache invalidated successfully"
}
```

## Features

### Dashboard Metrics
- **Total Revenue**: Sum of all completed orders
- **Total Orders**: Count of all orders (excluding cancelled)
- **Total Customers**: Unique customer count
- **Average Order Value**: Total revenue / Total orders
- **Today Metrics**: Revenue and orders from current day
- **Month Metrics**: Revenue and orders from current month

### Revenue Reports
- **Time-based Analysis**: Daily, weekly, monthly, or yearly reports
- **Date Range Filtering**: Custom date ranges for analysis
- **Export Options**: CSV and PDF formats
- **Cached Results**: 5-minute cache for performance

### Product Analytics
- **Top Selling Products**: By revenue or quantity sold
- **Product Performance**: Total sales and revenue per product
- **Export Options**: CSV and PDF formats

### Customer Analytics
- **Customer Lifetime Value**: Total spent per customer
- **Order Frequency**: Number of orders per customer
- **Last Order Date**: Most recent purchase tracking
- **Top Customers**: By total spending

### Report Management
- **Save Reports**: Persist generated reports for later access
- **Report History**: Access previously generated reports
- **Report Types**: Dashboard, revenue, products, customers
- **JSON Storage**: Flexible JSONB storage for report data

### Performance Optimization
- **Redis Caching**: 5-minute cache for expensive queries
- **Service Communication**: Fetches data from Order Service
- **Efficient Aggregation**: In-memory data processing
- **Pagination Support**: Handle large datasets efficiently

## Architecture

### Data Flow

```
Reporting Service 
    ↓ (fetch orders)
Order Service
    ↓ (aggregation & calculation)
Dashboard Metrics / Reports
    ↓ (export)
CSV / PDF Files
```

### Caching Strategy

- **Dashboard Metrics**: Cached for 5 minutes
- **Revenue Reports**: Cached per date range and period (5 minutes)
- **Cache Invalidation**: Manual invalidation via admin endpoint
- **Cache Keys**: 
  - `dashboard:metrics`
  - `revenue:report:{startDate}_{endDate}_{period}`

### Report Types

1. **Dashboard**: Overview metrics with totals and trends
2. **Revenue**: Time-series revenue and order data
3. **Products**: Top-selling products by revenue
4. **Customers**: Customer lifetime value and order frequency

## Database Schema

### Reporting Service (reporting_db)

**reports table:**
- id (UUID), name, type, period
- start_date, end_date
- data (JSONB) - flexible report data storage
- generated_by (UUID FK to user)
- created_at, updated_at, deleted_at

## Export Formats

### CSV Export
- Simple comma-separated values
- Headers included
- Compatible with Excel, Google Sheets
- Suitable for further data analysis

### PDF Export
- Professional formatted reports
- Tables with headers
- Generated timestamp
- Portrait (dashboard, revenue, customers)
- Landscape (products with wider tables)

## Use Cases

### Business Intelligence
1. Admin views dashboard for daily overview
2. Admin generates monthly revenue report
3. Admin exports report as PDF for stakeholders
4. Admin identifies top-selling products
5. Admin saves report for historical comparison

### Performance Monitoring
1. Monitor today's revenue vs. historical average
2. Track order volume trends
3. Identify revenue fluctuations
4. Analyze customer purchase patterns

### Customer Insights
1. Identify top customers by spending
2. Track customer lifetime value
3. Monitor customer order frequency
4. Analyze last order dates for re-engagement

### Product Analysis
1. Identify best-selling products
2. Track product revenue contribution
3. Monitor inventory needs based on sales
4. Plan promotions for underperforming products

## Performance Considerations

### Caching
- Dashboard metrics cached for 5 minutes
- Revenue reports cached per query parameters
- Cache invalidation on demand
- Redis used for distributed caching

### Data Aggregation
- Fetches all orders in batches (100 per page)
- In-memory aggregation for speed
- Filtering by date range and status
- Efficient map-based grouping

### Optimization Tips
1. Use date ranges to limit data scope
2. Leverage cached results when possible
3. Export large datasets as CSV for processing
4. Invalidate cache after significant order changes

## Health Checks

```bash
# Reporting Service
curl http://localhost:8089/health

# Check if ready (DB + Redis connected)
curl http://localhost:8089/ready
```

## Troubleshooting

### Service Issues

1. **Check database connectivity:**
```bash
psql -U postgres -h localhost -p 5441 -d reporting_db
```

2. **View saved reports:**
```sql
SELECT id, name, type, period, created_at FROM reports ORDER BY created_at DESC LIMIT 10;
```

3. **Check Redis connection:**
```bash
redis-cli -h localhost -p 6379 ping
```

4. **View cached keys:**
```bash
redis-cli -h localhost -p 6379 KEYS "*"
```

### Data Issues

1. **No orders showing in reports:**
   - Verify Order Service is running
   - Check JWT token has correct permissions
   - Ensure orders exist in order_db

2. **Cache not updating:**
   - Use admin endpoint to invalidate cache
   - Check Redis is running
   - Restart reporting service

3. **Export failing:**
   - Check sufficient memory for large datasets
   - Verify PDF library dependencies installed
   - Limit dataset size with date ranges

### Common Errors

**"Failed to fetch orders":**
- Order Service is down or unreachable
- Invalid JWT token
- Network connectivity issues

**"Failed to connect to Redis":**
- Redis service not running
- Incorrect Redis host/port
- Redis password mismatch

**"Report not found":**
- Report ID is invalid
- Report was deleted
- User lacks permissions

## Integration Examples

### Generate and Save Monthly Report

```bash
# 1. Generate revenue report
curl -X GET "http://localhost:8089/api/v1/reports/revenue?startDate=2026-06-01&endDate=2026-06-30&period=daily" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 2. Save the report
curl -X POST "http://localhost:8089/api/v1/reports" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "June 2026 Daily Revenue",
    "type": "revenue",
    "period": "daily",
    "startDate": "2026-06-01T00:00:00Z",
    "endDate": "2026-06-30T00:00:00Z",
    "data": {...}
  }'

# 3. Export as PDF
curl -X GET "http://localhost:8089/api/v1/reports/revenue/export?startDate=2026-06-01&endDate=2026-06-30&period=daily&format=pdf" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  --output june_report.pdf
```

### Scheduled Report Generation

Create a cron job to generate daily reports:

```bash
# Crontab entry (runs at 11:59 PM daily)
59 23 * * * curl -X GET "http://localhost:8089/api/v1/dashboard" \
  -H "Authorization: Bearer ADMIN_TOKEN" >> /var/log/daily_report.log
```

## Next Steps

Phase 7 will add:
- **React Frontend** (2 weeks)
  - Product catalog and search
  - Shopping cart and checkout
  - User authentication and profile
  - Order history and tracking
  - Admin dashboard with reports
  - Responsive design

Future phases:
- Phase 8: Dockerization & Container Orchestration
- Phase 9: Kafka Event Streaming
- Phase 10: Kubernetes Deployment
- Phase 11: Monitoring & Logging (Prometheus, Grafana, ELK)
- Phase 12: CI/CD Pipeline
- Phase 13: API Gateway & Load Balancing

## Resources

- [Reporting Service Documentation](reporting-service/README.md)
- [gofpdf Documentation](https://github.com/jung-kurt/gofpdf)
- [Redis Caching Guide](https://redis.io/docs/manual/client-side-caching/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
