# Product Service

**Port:** 8082  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**Cache:** Redis  
**ORM:** GORM

## Responsibilities

- **Product Management** - CRUD operations for products
- **Category Management** - Manage product categories
- **Product Search** - Search and filter products
- **Reviews & Ratings** - Customer reviews and ratings
- **Image Management** - Product image uploads
- **Inventory Integration** - Check stock availability

## Database Schema

### products
- id (UUID, PK)
- name
- description
- price
- category_id (FK)
- brand
- sku (unique)
- images (JSON array)
- average_rating
- total_reviews
- created_at
- updated_at

### categories
- id (UUID, PK)
- name (unique)
- slug (unique)
- parent_id (FK, nullable)
- description
- created_at
- updated_at

### reviews
- id (UUID, PK)
- product_id (FK)
- user_id
- rating (1-5)
- title
- comment
- verified_purchase
- created_at
- updated_at

## API Endpoints

### Products
- `GET /api/v1/products` - List products (paginated, filterable)
- `GET /api/v1/products/:id` - Get product details
- `POST /api/v1/products` - Create product (admin)
- `PUT /api/v1/products/:id` - Update product (admin)
- `DELETE /api/v1/products/:id` - Delete product (admin)
- `POST /api/v1/products/:id/images` - Upload product images (admin)

### Categories
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/:id` - Get category
- `POST /api/v1/categories` - Create category (admin)
- `PUT /api/v1/categories/:id` - Update category (admin)
- `DELETE /api/v1/categories/:id` - Delete category (admin)

### Reviews
- `GET /api/v1/products/:id/reviews` - Get product reviews
- `POST /api/v1/products/:id/reviews` - Add review
- `PUT /api/v1/reviews/:id` - Update review
- `DELETE /api/v1/reviews/:id` - Delete review

### Search
- `GET /api/v1/products/search?q=keyword` - Search products
- `GET /api/v1/products?category=&minPrice=&maxPrice=&sort=` - Filter products

## Redis Caching

Cache Strategy:
- Product details: 1 hour TTL
- Category list: 24 hours TTL
- Search results: 15 minutes TTL
- Product list by category: 30 minutes TTL

## Events Published

- `product.created`
- `product.updated`
- `product.deleted`
- `review.created`

## Directory Structure

```
product-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── cache/
│   └── events/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```

Test