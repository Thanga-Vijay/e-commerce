# E-Commerce Platform - Phase 2 Implementation

This document provides instructions for running the Auth Service and Product Service locally.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15
- Redis 7
- Docker and Docker Compose (optional)

## Quick Start with Docker Compose

The easiest way to run all services is using Docker Compose:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

Services will be available at:
- **Auth Service**: http://localhost:8081
- **Product Service**: http://localhost:8082
- **PostgreSQL (Auth DB)**: localhost:5433
- **PostgreSQL (Product DB)**: localhost:5434
- **Redis**: localhost:6379

## Manual Setup

### 1. Set Up Databases

```bash
# Create databases
createdb auth_db
createdb product_db

# Or using psql
psql -U postgres -c "CREATE DATABASE auth_db;"
psql -U postgres -c "CREATE DATABASE product_db;"
```

### 2. Run Migrations

**Auth Service:**
```bash
cd auth-service
# Using golang-migrate
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable" up
```

**Product Service:**
```bash
cd product-service
# Using golang-migrate
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/product_db?sslmode=disable" up
```

Or let GORM auto-migrate when you run the services.

### 3. Configure Environment

**Auth Service:**
```bash
cd auth-service
cp .env.example .env
# Edit .env with your configuration
```

**Product Service:**
```bash
cd product-service
cp .env.example .env
# Edit .env with your configuration
```

### 4. Install Dependencies

**Auth Service:**
```bash
cd auth-service
go mod download
```

**Product Service:**
```bash
cd product-service
go mod download
```

### 5. Run Services

**Terminal 1 - Auth Service:**
```bash
cd auth-service
go run cmd/main.go
```

**Terminal 2 - Product Service:**
```bash
cd product-service
go run cmd/main.go
```

## API Testing

### Auth Service Endpoints

**Register a new user:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!"
  }'
```

**Get Profile:**
```bash
curl -X GET http://localhost:8081/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Product Service Endpoints

**Create Category (Admin):**
```bash
curl -X POST http://localhost:8082/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and accessories"
  }'
```

**Get Categories:**
```bash
curl -X GET http://localhost:8082/api/v1/categories
```

**Create Product (Admin):**
```bash
curl -X POST http://localhost:8082/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "sku": "LAPTOP-001",
    "name": "MacBook Pro 16",
    "description": "Apple MacBook Pro 16-inch",
    "price": 2499.99,
    "categoryId": "CATEGORY_UUID",
    "stock": 10,
    "images": ["image1.jpg", "image2.jpg"]
  }'
```

**Get Products:**
```bash
curl -X GET "http://localhost:8082/api/v1/products?page=1&pageSize=20"
```

**Create Review:**
```bash
curl -X POST http://localhost:8082/api/v1/reviews \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "productId": "PRODUCT_UUID",
    "rating": 5,
    "title": "Great product!",
    "comment": "This product exceeded my expectations."
  }'
```

## Health Checks

```bash
# Auth Service
curl http://localhost:8081/health

# Product Service
curl http://localhost:8082/health
```

## Development

### Hot Reload

Install Air for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

Run with hot reload:
```bash
# Auth Service
cd auth-service
air

# Product Service
cd product-service
air
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Troubleshooting

### Database Connection Issues

1. Ensure PostgreSQL is running:
```bash
pg_isctl status
```

2. Check database credentials in `.env` files

3. Test database connection:
```bash
psql -U postgres -h localhost -d auth_db
```

### Redis Connection Issues

1. Ensure Redis is running:
```bash
redis-cli ping
```

2. Check Redis configuration in `.env` files

### Port Conflicts

If ports are already in use, update the `PORT` environment variable in `.env` files or `docker-compose.yml`.

## Next Steps

- Implement remaining microservices (Cart, Wishlist, Order, Payment, Inventory, Notification, Reporting)
- Set up API Gateway
- Implement Kafka event messaging
- Deploy to Kubernetes

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
