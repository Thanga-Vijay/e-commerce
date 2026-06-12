# API Contracts

Complete REST API specifications for all microservices in the E-Commerce Platform.

## Table of Contents
1. [API Gateway](#api-gateway)
2. [Auth Service](#auth-service)
3. [Product Service](#product-service)
4. [Cart Service](#cart-service)
5. [Wishlist Service](#wishlist-service)
6. [Order Service](#order-service)
7. [Payment Service](#payment-service)
8. [Inventory Service](#inventory-service)
9. [Reporting Service](#reporting-service)
10. [Standard Response Format](#standard-response-format)
11. [Error Codes](#error-codes)
12. [Authentication](#authentication)

---

## API Base URL

All services are accessed through the API Gateway:

**Base URL:** `http://localhost:8080/api/v1`

---

## Auth Service

**Base Path:** `/auth`

### 1. Register User

**Endpoint:** `POST /auth/register`

**Description:** Create a new user account

**Authentication:** No

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response:** `201 Created`
```json
{
  "status": 201,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "customer",
      "isVerified": false,
      "createdAt": "2026-06-12T10:30:00Z"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
  }
}
```

**Validation Rules:**
- Email: Valid email format, unique
- Password: Min 8 characters, at least 1 uppercase, 1 lowercase, 1 number, 1 special character
- FirstName: Min 2 characters, max 100
- LastName: Min 2 characters, max 100

---

### 2. Login

**Endpoint:** `POST /auth/login`

**Description:** Authenticate user and receive JWT tokens

**Authentication:** No

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "customer",
      "isVerified": true
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
  }
}
```

**Rate Limiting:** 5 attempts per 15 minutes per IP

---

### 3. Refresh Token

**Endpoint:** `POST /auth/refresh`

**Description:** Refresh expired access token

**Authentication:** No

**Request Body:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Token refreshed",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
  }
}
```

---

### 4. Logout

**Endpoint:** `POST /auth/logout`

**Description:** Invalidate user tokens

**Authentication:** Required

**Request Body:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Logout successful",
  "data": null
}
```

---

### 5. Get Current User

**Endpoint:** `GET /auth/me`

**Description:** Get authenticated user profile

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "User retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "customer",
    "isVerified": true,
    "createdAt": "2026-06-12T10:30:00Z",
    "updatedAt": "2026-06-12T10:30:00Z"
  }
}
```

---

### 6. Update Profile

**Endpoint:** `PUT /auth/profile`

**Description:** Update user profile

**Authentication:** Required

**Request Body:**
```json
{
  "firstName": "Jane",
  "lastName": "Smith"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Profile updated",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "firstName": "Jane",
    "lastName": "Smith",
    "role": "customer",
    "isVerified": true,
    "updatedAt": "2026-06-12T11:00:00Z"
  }
}
```

---

### 7. Forgot Password

**Endpoint:** `POST /auth/forgot-password`

**Description:** Request password reset email

**Authentication:** No

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Password reset email sent",
  "data": null
}
```

---

### 8. Reset Password

**Endpoint:** `POST /auth/reset-password`

**Description:** Reset password using token from email

**Authentication:** No

**Request Body:**
```json
{
  "token": "reset-token-from-email",
  "newPassword": "NewSecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Password reset successful",
  "data": null
}
```

---

### 9. Verify Email

**Endpoint:** `POST /auth/verify-email`

**Description:** Verify user email address

**Authentication:** No

**Request Body:**
```json
{
  "token": "verification-token-from-email"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Email verified successfully",
  "data": null
}
```

---

## Product Service

**Base Path:** `/products`

### 1. List Products

**Endpoint:** `GET /products`

**Description:** Get paginated list of products with filtering and sorting

**Authentication:** No

**Query Parameters:**
- `page` (int, default: 1) - Page number
- `limit` (int, default: 20, max: 100) - Items per page
- `category` (string) - Filter by category ID or slug
- `minPrice` (decimal) - Minimum price
- `maxPrice` (decimal) - Maximum price
- `minRating` (int, 1-5) - Minimum rating
- `brand` (string) - Filter by brand
- `sort` (string) - Sort by: `price_asc`, `price_desc`, `rating`, `newest`, `popular`
- `search` (string) - Search query

**Example:** `GET /products?page=1&limit=20&category=electronics&minPrice=100&maxPrice=1000&sort=price_asc`

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Products retrieved",
  "data": {
    "products": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Laptop",
        "description": "High-performance laptop",
        "price": 999.99,
        "category": {
          "id": "550e8400-e29b-41d4-a716-446655440100",
          "name": "Electronics",
          "slug": "electronics"
        },
        "brand": "TechBrand",
        "sku": "LAPTOP-001",
        "images": [
          "https://cdn.example.com/laptop-1.jpg",
          "https://cdn.example.com/laptop-2.jpg"
        ],
        "averageRating": 4.5,
        "totalReviews": 127,
        "createdAt": "2026-06-01T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "totalPages": 8
    }
  }
}
```

---

### 2. Get Product by ID

**Endpoint:** `GET /products/:id`

**Description:** Get detailed product information

**Authentication:** No

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Product retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Laptop",
    "description": "High-performance laptop with latest specs",
    "price": 999.99,
    "category": {
      "id": "550e8400-e29b-41d4-a716-446655440100",
      "name": "Electronics",
      "slug": "electronics",
      "parent": null
    },
    "brand": "TechBrand",
    "sku": "LAPTOP-001",
    "images": [
      "https://cdn.example.com/laptop-1.jpg",
      "https://cdn.example.com/laptop-2.jpg"
    ],
    "averageRating": 4.5,
    "totalReviews": 127,
    "stockStatus": "in_stock",
    "createdAt": "2026-06-01T10:00:00Z",
    "updatedAt": "2026-06-10T15:30:00Z"
  }
}
```

---

### 3. Create Product (Admin)

**Endpoint:** `POST /products`

**Description:** Create a new product

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "categoryId": "550e8400-e29b-41d4-a716-446655440100",
  "brand": "TechBrand",
  "sku": "LAPTOP-001",
  "images": [
    "https://cdn.example.com/laptop-1.jpg"
  ]
}
```

**Response:** `201 Created`
```json
{
  "status": 201,
  "message": "Product created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "categoryId": "550e8400-e29b-41d4-a716-446655440100",
    "brand": "TechBrand",
    "sku": "LAPTOP-001",
    "images": [
      "https://cdn.example.com/laptop-1.jpg"
    ],
    "averageRating": 0,
    "totalReviews": 0,
    "createdAt": "2026-06-12T10:30:00Z"
  }
}
```

---

### 4. Update Product (Admin)

**Endpoint:** `PUT /products/:id`

**Description:** Update product details

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "name": "Updated Laptop",
  "price": 899.99,
  "description": "Updated description"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Product updated",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Updated Laptop",
    "price": 899.99,
    "description": "Updated description",
    "updatedAt": "2026-06-12T11:00:00Z"
  }
}
```

---

### 5. Delete Product (Admin)

**Endpoint:** `DELETE /products/:id`

**Description:** Soft delete product

**Authentication:** Required (Admin)

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Product deleted",
  "data": null
}
```

---

### 6. List Categories

**Endpoint:** `GET /categories`

**Description:** Get all product categories (hierarchical)

**Authentication:** No

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Categories retrieved",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440100",
      "name": "Electronics",
      "slug": "electronics",
      "description": "Electronic devices",
      "parentId": null,
      "children": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440101",
          "name": "Laptops",
          "slug": "laptops",
          "parentId": "550e8400-e29b-41d4-a716-446655440100",
          "children": []
        }
      ]
    }
  ]
}
```

---

### 7. Get Product Reviews

**Endpoint:** `GET /products/:id/reviews`

**Description:** Get reviews for a product

**Authentication:** No

**Query Parameters:**
- `page` (int, default: 1)
- `limit` (int, default: 10)
- `sort` (string) - Sort by: `newest`, `oldest`, `rating_high`, `rating_low`

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Reviews retrieved",
  "data": {
    "reviews": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440200",
        "productId": "550e8400-e29b-41d4-a716-446655440001",
        "userId": "550e8400-e29b-41d4-a716-446655440050",
        "rating": 5,
        "title": "Excellent laptop!",
        "comment": "Very fast and reliable",
        "verifiedPurchase": true,
        "createdAt": "2026-06-10T14:20:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 127,
      "totalPages": 13
    }
  }
}
```

---

### 8. Create Review

**Endpoint:** `POST /products/:id/reviews`

**Description:** Add a review for a product

**Authentication:** Required

**Request Body:**
```json
{
  "rating": 5,
  "title": "Excellent product!",
  "comment": "This product exceeded my expectations"
}
```

**Response:** `201 Created`
```json
{
  "status": 201,
  "message": "Review created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440200",
    "productId": "550e8400-e29b-41d4-a716-446655440001",
    "userId": "550e8400-e29b-41d4-a716-446655440050",
    "rating": 5,
    "title": "Excellent product!",
    "comment": "This product exceeded my expectations",
    "verifiedPurchase": false,
    "createdAt": "2026-06-12T10:30:00Z"
  }
}
```

---

## Cart Service

**Base Path:** `/cart`

### 1. Get Cart

**Endpoint:** `GET /cart`

**Description:** Get user's shopping cart

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Cart retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440300",
    "userId": "550e8400-e29b-41d4-a716-446655440050",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440301",
        "productId": "550e8400-e29b-41d4-a716-446655440001",
        "product": {
          "name": "Laptop",
          "images": ["https://cdn.example.com/laptop-1.jpg"],
          "currentPrice": 999.99
        },
        "quantity": 2,
        "price": 999.99,
        "subtotal": 1999.98
      }
    ],
    "subtotal": 1999.98,
    "tax": 199.99,
    "total": 2199.97,
    "itemCount": 2,
    "updatedAt": "2026-06-12T10:30:00Z"
  }
}
```

---

### 2. Add Item to Cart

**Endpoint:** `POST /cart/items`

**Description:** Add product to cart

**Authentication:** Required

**Request Body:**
```json
{
  "productId": "550e8400-e29b-41d4-a716-446655440001",
  "quantity": 2
}
```

**Response:** `201 Created`
```json
{
  "status": 201,
  "message": "Item added to cart",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440301",
    "cartId": "550e8400-e29b-41d4-a716-446655440300",
    "productId": "550e8400-e29b-41d4-a716-446655440001",
    "quantity": 2,
    "price": 999.99,
    "subtotal": 1999.98,
    "createdAt": "2026-06-12T10:30:00Z"
  }
}
```

---

### 3. Update Cart Item

**Endpoint:** `PUT /cart/items/:id`

**Description:** Update item quantity

**Authentication:** Required

**Request Body:**
```json
{
  "quantity": 3
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Cart item updated",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440301",
    "quantity": 3,
    "subtotal": 2999.97,
    "updatedAt": "2026-06-12T10:35:00Z"
  }
}
```

---

### 4. Remove Cart Item

**Endpoint:** `DELETE /cart/items/:id`

**Description:** Remove item from cart

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Item removed from cart",
  "data": null
}
```

---

### 5. Clear Cart

**Endpoint:** `DELETE /cart`

**Description:** Remove all items from cart

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Cart cleared",
  "data": null
}
```

---

## Order Service

**Base Path:** `/orders`

### 1. Create Order

**Endpoint:** `POST /orders`

**Description:** Create order from cart

**Authentication:** Required

**Request Body:**
```json
{
  "shippingAddress": {
    "firstName": "John",
    "lastName": "Doe",
    "addressLine1": "123 Main St",
    "addressLine2": "Apt 4B",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "US",
    "phone": "+1234567890"
  },
  "billingAddress": {
    "firstName": "John",
    "lastName": "Doe",
    "addressLine1": "123 Main St",
    "addressLine2": "Apt 4B",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "US",
    "phone": "+1234567890"
  }
}
```

**Response:** `201 Created`
```json
{
  "status": 201,
  "message": "Order created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "orderNumber": "ORD-2026-001234",
    "userId": "550e8400-e29b-41d4-a716-446655440050",
    "status": "pending",
    "items": [
      {
        "productId": "550e8400-e29b-41d4-a716-446655440001",
        "productName": "Laptop",
        "productPrice": 999.99,
        "quantity": 2,
        "subtotal": 1999.98
      }
    ],
    "subtotal": 1999.98,
    "tax": 199.99,
    "shippingCost": 10.00,
    "total": 2209.97,
    "shippingAddress": {
      "firstName": "John",
      "lastName": "Doe",
      "addressLine1": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "US"
    },
    "createdAt": "2026-06-12T10:30:00Z"
  }
}
```

---

### 2. Get User Orders

**Endpoint:** `GET /orders`

**Description:** Get authenticated user's orders

**Authentication:** Required

**Query Parameters:**
- `page` (int, default: 1)
- `limit` (int, default: 10)
- `status` (string) - Filter by status

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Orders retrieved",
  "data": {
    "orders": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440400",
        "orderNumber": "ORD-2026-001234",
        "status": "confirmed",
        "total": 2209.97,
        "itemCount": 2,
        "createdAt": "2026-06-12T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 5,
      "totalPages": 1
    }
  }
}
```

---

### 3. Get Order by ID

**Endpoint:** `GET /orders/:id`

**Description:** Get order details

**Authentication:** Required (Own orders only, admin can view all)

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Order retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "orderNumber": "ORD-2026-001234",
    "userId": "550e8400-e29b-41d4-a716-446655440050",
    "status": "confirmed",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440401",
        "productId": "550e8400-e29b-41d4-a716-446655440001",
        "productName": "Laptop",
        "productPrice": 999.99,
        "quantity": 2,
        "subtotal": 1999.98
      }
    ],
    "subtotal": 1999.98,
    "tax": 199.99,
    "shippingCost": 10.00,
    "total": 2209.97,
    "shippingAddress": {
      "firstName": "John",
      "lastName": "Doe",
      "addressLine1": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "US"
    },
    "billingAddress": {
      "firstName": "John",
      "lastName": "Doe",
      "addressLine1": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "US"
    },
    "statusHistory": [
      {
        "status": "pending",
        "timestamp": "2026-06-12T10:30:00Z"
      },
      {
        "status": "confirmed",
        "timestamp": "2026-06-12T10:35:00Z"
      }
    ],
    "createdAt": "2026-06-12T10:30:00Z",
    "updatedAt": "2026-06-12T10:35:00Z"
  }
}
```

---

### 4. Cancel Order

**Endpoint:** `PUT /orders/:id/cancel`

**Description:** Cancel pending or confirmed order

**Authentication:** Required

**Request Body:**
```json
{
  "reason": "Changed my mind"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Order cancelled",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "status": "cancelled",
    "updatedAt": "2026-06-12T11:00:00Z"
  }
}
```

---

### 5. Track Order

**Endpoint:** `GET /orders/:id/track`

**Description:** Get order tracking information

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Tracking info retrieved",
  "data": {
    "orderId": "550e8400-e29b-41d4-a716-446655440400",
    "orderNumber": "ORD-2026-001234",
    "currentStatus": "shipped",
    "trackingNumber": "1Z999AA1012345678",
    "carrier": "UPS",
    "estimatedDelivery": "2026-06-15T17:00:00Z",
    "timeline": [
      {
        "status": "pending",
        "timestamp": "2026-06-12T10:30:00Z",
        "description": "Order placed"
      },
      {
        "status": "confirmed",
        "timestamp": "2026-06-12T10:35:00Z",
        "description": "Payment confirmed"
      },
      {
        "status": "processing",
        "timestamp": "2026-06-12T12:00:00Z",
        "description": "Order being processed"
      },
      {
        "status": "shipped",
        "timestamp": "2026-06-13T09:00:00Z",
        "description": "Order shipped"
      }
    ]
  }
}
```

---

## Payment Service

**Base Path:** `/payments`

### 1. Create Payment Intent

**Endpoint:** `POST /payments/intent`

**Description:** Create Stripe payment intent for order

**Authentication:** Required

**Request Body:**
```json
{
  "orderId": "550e8400-e29b-41d4-a716-446655440400"
}
```

**Response:** `201 Created`
```json
{
  "status": 201,
  "message": "Payment intent created",
  "data": {
    "paymentId": "550e8400-e29b-41d4-a716-446655440500",
    "clientSecret": "pi_xxx_secret_yyy",
    "amount": 2209.97,
    "currency": "usd"
  }
}
```

---

### 2. Get Payment Status

**Endpoint:** `GET /payments/:id`

**Description:** Get payment details

**Authentication:** Required

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Payment retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440500",
    "orderId": "550e8400-e29b-41d4-a716-446655440400",
    "amount": 2209.97,
    "currency": "usd",
    "status": "succeeded",
    "paymentMethod": "card",
    "createdAt": "2026-06-12T10:30:00Z",
    "updatedAt": "2026-06-12T10:35:00Z"
  }
}
```

---

### 3. Process Refund (Admin)

**Endpoint:** `POST /payments/:id/refund`

**Description:** Refund a payment

**Authentication:** Required (Admin)

**Request Body:**
```json
{
  "amount": 2209.97,
  "reason": "Customer request"
}
```

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Refund processed",
  "data": {
    "refundId": "550e8400-e29b-41d4-a716-446655440501",
    "paymentId": "550e8400-e29b-41d4-a716-446655440500",
    "amount": 2209.97,
    "status": "succeeded",
    "createdAt": "2026-06-12T11:00:00Z"
  }
}
```

---

## Reporting Service

**Base Path:** `/reports`

### 1. Dashboard Metrics (Admin)

**Endpoint:** `GET /reports/dashboard`

**Description:** Get overall dashboard metrics

**Authentication:** Required (Admin)

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Dashboard metrics retrieved",
  "data": {
    "revenue": {
      "total": 1250000.00,
      "today": 15000.00,
      "weekly": 105000.00,
      "monthly": 450000.00
    },
    "orders": {
      "total": 5234,
      "pending": 23,
      "confirmed": 45,
      "shipped": 67,
      "delivered": 5099
    },
    "products": {
      "total": 450,
      "lowStock": 12,
      "outOfStock": 3
    },
    "customers": {
      "total": 12340,
      "new": 45,
      "active": 8900
    },
    "charts": {
      "revenueChart": [
        {"date": "2026-06-05", "value": 14500.00},
        {"date": "2026-06-06", "value": 15200.00},
        {"date": "2026-06-07", "value": 13800.00},
        {"date": "2026-06-08", "value": 16100.00},
        {"date": "2026-06-09", "value": 14900.00},
        {"date": "2026-06-10", "value": 15600.00},
        {"date": "2026-06-11", "value": 14900.00},
        {"date": "2026-06-12", "value": 15000.00}
      ]
    }
  }
}
```

---

### 2. Revenue Analytics (Admin)

**Endpoint:** `GET /reports/revenue`

**Description:** Get revenue analytics

**Authentication:** Required (Admin)

**Query Parameters:**
- `from` (date, ISO 8601) - Start date
- `to` (date, ISO 8601) - End date
- `groupBy` (string) - Group by: `day`, `week`, `month`

**Response:** `200 OK`
```json
{
  "status": 200,
  "message": "Revenue analytics retrieved",
  "data": {
    "totalRevenue": 450000.00,
    "averageOrderValue": 420.50,
    "totalOrders": 1070,
    "breakdown": [
      {
        "period": "2026-06-01",
        "revenue": 15000.00,
        "orders": 35
      },
      {
        "period": "2026-06-02",
        "revenue": 14500.00,
        "orders": 33
      }
    ]
  }
}
```

---

### 3. Export Report (Admin)

**Endpoint:** `GET /reports/export`

**Description:** Export report as CSV or PDF

**Authentication:** Required (Admin)

**Query Parameters:**
- `type` (string) - Report type: `sales`, `orders`, `customers`, `products`
- `format` (string) - Format: `csv`, `pdf`
- `from` (date) - Start date
- `to` (date) - End date

**Response:** `200 OK` (File download)

---

## Standard Response Format

All API responses follow this structure:

### Success Response
```json
{
  "status": 200,
  "message": "Operation successful",
  "data": {}
}
```

### Error Response
```json
{
  "status": 400,
  "message": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "email": "Invalid email format",
    "password": "Password too short"
  }
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `UNAUTHORIZED` | 401 | Authentication required |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

### Specific Error Codes

| Code | Description |
|------|-------------|
| `EMAIL_EXISTS` | Email already registered |
| `INVALID_CREDENTIALS` | Email or password incorrect |
| `TOKEN_EXPIRED` | JWT token expired |
| `TOKEN_INVALID` | Invalid token |
| `PRODUCT_NOT_FOUND` | Product does not exist |
| `INSUFFICIENT_STOCK` | Not enough stock available |
| `CART_EMPTY` | Cart is empty |
| `ORDER_NOT_FOUND` | Order does not exist |
| `PAYMENT_FAILED` | Payment processing failed |
| `REFUND_FAILED` | Refund processing failed |

---

## Authentication

### JWT Token

All authenticated endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Claims

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "role": "customer",
  "iat": 1718185800,
  "exp": 1718186700
}
```

### Token Expiry
- **Access Token:** 15 minutes
- **Refresh Token:** 7 days

### Rate Limiting

Rate limits per endpoint:

| Endpoint | Rate Limit |
|----------|-----------|
| `/auth/login` | 5 per 15 minutes per IP |
| `/auth/register` | 3 per hour per IP |
| Public endpoints | 100 per minute per IP |
| Authenticated endpoints | 1000 per minute per user |
| Admin endpoints | 10000 per minute per user |

---

This API contract provides a complete specification for all microservices with consistent request/response formats, proper error handling, and security considerations.
