# Cart Service

**Port:** 8083  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**Cache:** Redis  
**ORM:** GORM

## Responsibilities

- **Cart Management** - Add, update, remove items
- **Cart Persistence** - Save cart across sessions
- **Price Calculation** - Calculate cart totals
- **Stock Validation** - Verify product availability
- **Cart Expiry** - Clean up abandoned carts

## Database Schema

### carts
- id (UUID, PK)
- user_id (unique)
- created_at
- updated_at
- expires_at

### cart_items
- id (UUID, PK)
- cart_id (FK)
- product_id
- quantity
- price (snapshot at time of adding)
- created_at
- updated_at

## API Endpoints

- `GET /api/v1/cart` - Get user's cart
- `POST /api/v1/cart/items` - Add item to cart
- `PUT /api/v1/cart/items/:id` - Update item quantity
- `DELETE /api/v1/cart/items/:id` - Remove item
- `DELETE /api/v1/cart` - Clear cart
- `GET /api/v1/cart/summary` - Get cart summary (total, tax, etc.)

## Redis Caching

- Cart data cached with 24-hour TTL
- Cache invalidated on cart updates
- Use user_id as cache key

## Integration

- Calls Product Service to validate products and prices
- Calls Inventory Service to check stock availability

## Directory Structure

```
cart-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── cache/
│   └── clients/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```

Test