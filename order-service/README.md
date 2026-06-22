# Order Service

**Port:** 8085  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**ORM:** GORM

## Responsibilities

- **Order Creation** - Create orders from cart
- **Order Management** - Update order status
- **Order Tracking** - Track order progress
- **Order History** - View past orders
- **Order Cancellation** - Cancel pending orders

## Database Schema

### orders
- id (UUID, PK)
- user_id
- order_number (unique)
- status (pending, confirmed, shipped, delivered, cancelled)
- subtotal
- tax
- shipping_cost
- total
- shipping_address (JSON)
- billing_address (JSON)
- created_at
- updated_at

### order_items
- id (UUID, PK)
- order_id (FK)
- product_id
- product_name (snapshot)
- product_price (snapshot)
- quantity
- subtotal
- created_at

### order_status_history
- id (UUID, PK)
- order_id (FK)
- status
- comment
- created_at

## API Endpoints

- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - Get user's orders
- `GET /api/v1/orders/:id` - Get order details
- `PUT /api/v1/orders/:id/cancel` - Cancel order
- `GET /api/v1/orders/:id/track` - Track order
- `PUT /api/v1/orders/:id/status` - Update status (admin)

## State Machine

Order Status Flow:
1. `pending` - Order created, awaiting payment
2. `confirmed` - Payment confirmed
3. `processing` - Order being prepared
4. `shipped` - Order shipped
5. `delivered` - Order delivered
6. `cancelled` - Order cancelled

## Events Published

- `order.created`
- `order.confirmed`
- `order.shipped`
- `order.delivered`
- `order.cancelled`

## Events Consumed

- `payment.completed` - Confirm order
- `payment.failed` - Cancel order

## Directory Structure

```
order-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── events/
│   └── statemachine/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```

Test