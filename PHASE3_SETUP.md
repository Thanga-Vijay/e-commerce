# E-Commerce Platform - Phase 3 Implementation

This document provides instructions for running the Cart Service and Wishlist Service (Phase 3).

## Services Overview

**Phase 3** adds two new microservices:
- **Cart Service** (Port 8083): Shopping cart management
- **Wishlist Service** (Port 8084): Wishlist management with move-to-cart functionality

## Quick Start with Docker Compose

Start all services including Phase 2 and Phase 3:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs for specific service
docker-compose logs -f cart-service
docker-compose logs -f wishlist-service

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
- **PostgreSQL Cart DB**: localhost:5435
- **PostgreSQL Wishlist DB**: localhost:5436
- **Redis**: localhost:6379

## Manual Setup

### 1. Set Up Databases

```bash
# Create databases
createdb cart_db
createdb wishlist_db

# Or using psql
psql -U postgres -c "CREATE DATABASE cart_db;"
psql -U postgres -c "CREATE DATABASE wishlist_db;"
```

### 2. Run Migrations

**Cart Service:**
```bash
cd cart-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/cart_db?sslmode=disable" up
```

**Wishlist Service:**
```bash
cd wishlist-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/wishlist_db?sslmode=disable" up
```

### 3. Configure Environment

**Cart Service:**
```bash
cd cart-service
cp .env.example .env
# Edit .env with your configuration
```

**Wishlist Service:**
```bash
cd wishlist-service
cp .env.example .env
# Edit .env with your configuration
```

### 4. Install Dependencies

```bash
# Cart Service
cd cart-service
go mod download

# Wishlist Service
cd wishlist-service
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

**Terminal 3 - Cart Service:**
```bash
cd cart-service
go run cmd/main.go
```

**Terminal 4 - Wishlist Service:**
```bash
cd wishlist-service
go run cmd/main.go
```

## API Testing

### Prerequisites

First, register and login to get an access token:

```bash
# Register
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!",
    "firstName": "John",
    "lastName": "Doe"
  }'

# Login (save the accessToken)
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!"
  }'
```

### Cart Service Endpoints

**Get Cart:**
```bash
curl -X GET http://localhost:8083/api/v1/cart \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Add Item to Cart:**
```bash
curl -X POST http://localhost:8083/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "productId": "PRODUCT_UUID",
    "quantity": 2
  }'
```

**Update Cart Item Quantity:**
```bash
curl -X PUT http://localhost:8083/api/v1/cart/items/ITEM_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "quantity": 3
  }'
```

**Remove Item from Cart:**
```bash
curl -X DELETE http://localhost:8083/api/v1/cart/items/ITEM_UUID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Clear Cart:**
```bash
curl -X DELETE http://localhost:8083/api/v1/cart \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Wishlist Service Endpoints

**Get Wishlist:**
```bash
curl -X GET http://localhost:8084/api/v1/wishlist \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Add Item to Wishlist:**
```bash
curl -X POST http://localhost:8084/api/v1/wishlist/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "productId": "PRODUCT_UUID"
  }'
```

**Remove Item from Wishlist:**
```bash
curl -X DELETE http://localhost:8084/api/v1/wishlist/items/ITEM_UUID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Move Item from Wishlist to Cart:**
```bash
curl -X POST http://localhost:8084/api/v1/wishlist/items/ITEM_UUID/move-to-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "quantity": 1
  }'
```

**Clear Wishlist:**
```bash
curl -X DELETE http://localhost:8084/api/v1/wishlist \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Features

### Cart Service

- **Get User Cart**: Retrieve current cart with items, subtotal, tax (8%), and total
- **Add to Cart**: Add products with quantity validation and stock checking
- **Update Quantity**: Change item quantities with stock validation
- **Remove Items**: Remove individual items from cart
- **Clear Cart**: Remove all items at once
- **Auto-calculation**: Automatic subtotal, tax, and total calculation
- **Price Snapshot**: Stores product price at time of adding
- **Cart Expiry**: Carts expire after 30 days of inactivity
- **Redis Caching**: 24-hour cache TTL for improved performance
- **Product Service Integration**: Real-time product validation and stock checking

### Wishlist Service

- **Get Wishlist**: Retrieve user's wishlist with all items
- **Add to Wishlist**: Add products with duplicate checking
- **Remove Items**: Remove items from wishlist
- **Clear Wishlist**: Remove all items at once
- **Move to Cart**: Move wishlist items to cart with quantity selection
- **Product Service Integration**: Real-time product validation
- **Cart Service Integration**: Seamless transfer of items to cart

## Architecture

### Service Communication

```
Wishlist Service --> Product Service (Get product details)
Wishlist Service --> Cart Service (Move item to cart)
Cart Service --> Product Service (Validate product and stock)
```

### Data Flow

1. **Add to Cart Flow**:
   - User sends request to Cart Service
   - Cart Service validates product with Product Service
   - Cart Service checks stock availability
   - Item added to cart with price snapshot
   - Cart cached in Redis
   - Response sent to user

2. **Move to Cart Flow**:
   - User requests to move wishlist item to cart
   - Wishlist Service retrieves item details
   - Wishlist Service calls Cart Service to add item
   - Cart Service validates and adds item
   - Wishlist Service removes item from wishlist
   - Response sent to user

## Database Schema

### Cart Service (cart_db)

**carts table:**
- id (UUID, PK)
- user_id (UUID, unique)
- expires_at (timestamp)
- created_at, updated_at, deleted_at

**cart_items table:**
- id (UUID, PK)
- cart_id (UUID, FK)
- product_id (UUID)
- product_name, product_price, product_image
- quantity (integer)
- subtotal (decimal)
- created_at, updated_at, deleted_at
- Unique constraint on (cart_id, product_id)

### Wishlist Service (wishlist_db)

**wishlists table:**
- id (UUID, PK)
- user_id (UUID, unique)
- created_at, updated_at, deleted_at

**wishlist_items table:**
- id (UUID, PK)
- wishlist_id (UUID, FK)
- product_id (UUID)
- product_name, product_price, product_image
- created_at, updated_at, deleted_at
- Unique constraint on (wishlist_id, product_id)

## Health Checks

```bash
# Cart Service
curl http://localhost:8083/health

# Wishlist Service
curl http://localhost:8084/health
```

## Troubleshooting

### Service Communication Issues

1. Ensure all services are running:
```bash
docker-compose ps
```

2. Check service logs:
```bash
docker-compose logs cart-service
docker-compose logs wishlist-service
```

3. Verify network connectivity:
```bash
docker-compose exec cart-service ping product-service
docker-compose exec wishlist-service ping cart-service
```

### Database Issues

1. Check database connectivity:
```bash
psql -U postgres -h localhost -p 5435 -d cart_db
psql -U postgres -h localhost -p 5436 -d wishlist_db
```

2. View database tables:
```sql
\dt
```

### Redis Cache Issues

1. Clear Redis cache:
```bash
docker-compose exec redis redis-cli FLUSHALL
```

2. Check cached keys:
```bash
docker-compose exec redis redis-cli KEYS "cart:*"
```

## Next Steps

Phase 4 will add:
- Order Service with state machine
- Payment Service with Stripe integration
- Order creation from cart
- Payment processing workflow

## Resources

- [Cart Service Documentation](cart-service/README.md)
- [Wishlist Service Documentation](wishlist-service/README.md)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
