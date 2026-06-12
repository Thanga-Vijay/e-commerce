# Reporting Service

**Port:** 8089  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL (read replicas)  
**ORM:** GORM

## Responsibilities

- **Sales Analytics** - Revenue, orders, products sold
- **Customer Analytics** - New customers, retention, lifetime value
- **Product Analytics** - Best sellers, categories, reviews
- **Dashboard Metrics** - Real-time KPIs
- **Report Generation** - CSV/PDF exports
- **Time-based Reports** - Daily, weekly, monthly reports

## Database Access

- Reads from PostgreSQL read replicas
- Aggregates data from multiple services
- Uses materialized views for performance

## API Endpoints

### Dashboard
- `GET /api/v1/reports/dashboard` - Overall dashboard metrics
- `GET /api/v1/reports/revenue` - Revenue analytics
- `GET /api/v1/reports/orders` - Order analytics
- `GET /api/v1/reports/customers` - Customer analytics
- `GET /api/v1/reports/products` - Product analytics

### Time-based Reports
- `GET /api/v1/reports/daily?date=YYYY-MM-DD` - Daily report
- `GET /api/v1/reports/weekly?week=YYYY-WW` - Weekly report
- `GET /api/v1/reports/monthly?month=YYYY-MM` - Monthly report
- `GET /api/v1/reports/yearly?year=YYYY` - Yearly report

### Export
- `GET /api/v1/reports/export/csv?type=sales&from=&to=` - Export CSV
- `GET /api/v1/reports/export/pdf?type=sales&from=&to=` - Export PDF

## Dashboard Metrics

### Revenue Metrics
- Total revenue (all-time)
- Today's revenue
- Weekly revenue
- Monthly revenue
- Yearly revenue
- Revenue by category
- Revenue trend (chart data)

### Order Metrics
- Total orders
- Pending orders
- Shipped orders
- Delivered orders
- Cancelled orders
- Average order value
- Orders by status (chart)

### Product Metrics
- Total products
- Products sold
- Best sellers (top 10)
- Worst sellers
- Average rating
- Total reviews
- Low stock products

### Customer Metrics
- Total customers
- New customers (today, week, month)
- Active customers
- Customer retention rate
- Customer lifetime value
- Top customers

## Aggregation Jobs

Background jobs that run periodically:
- `UpdateDailyStats()` - Runs daily at midnight
- `UpdateWeeklyStats()` - Runs weekly on Sunday
- `UpdateMonthlyStats()` - Runs monthly on 1st
- `RefreshMaterializedViews()` - Runs hourly

## Events Consumed

- `order.created`
- `order.completed`
- `payment.completed`
- `product.created`
- `user.created`

Updates real-time metrics in cache.

## Caching Strategy

- Dashboard metrics: 5 minutes TTL
- Daily reports: 1 hour TTL (if today, 5 minutes)
- Historical reports: 24 hours TTL

## Directory Structure

```
reporting-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── aggregator/
│   ├── exporter/
│   └── consumers/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```
