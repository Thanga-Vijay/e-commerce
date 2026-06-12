# Wishlist Service

**Port:** 8084  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**ORM:** GORM

## Responsibilities

- **Wishlist Management** - Add, remove products from wishlist
- **Wishlist Viewing** - View user's wishlist
- **Product Details** - Fetch product information
- **Wishlist Sharing** - Share wishlist (future feature)

## Database Schema

### wishlists
- id (UUID, PK)
- user_id (unique)
- created_at
- updated_at

### wishlist_items
- id (UUID, PK)
- wishlist_id (FK)
- product_id
- added_at
- created_at

## API Endpoints

- `GET /api/v1/wishlist` - Get user's wishlist
- `POST /api/v1/wishlist/items` - Add item to wishlist
- `DELETE /api/v1/wishlist/items/:id` - Remove item
- `DELETE /api/v1/wishlist` - Clear wishlist
- `POST /api/v1/wishlist/items/:id/move-to-cart` - Move to cart

## Integration

- Calls Product Service to fetch product details
- Calls Cart Service to move items to cart

## Directory Structure

```
wishlist-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── clients/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```
