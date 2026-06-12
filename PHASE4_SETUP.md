# E-Commerce Platform - Phase 4 Implementation

This document provides instructions for running the Order Service and Payment Service (Phase 4).

## Services Overview

**Phase 4** adds two new microservices:
- **Order Service** (Port 8085): Order management with state machine
- **Payment Service** (Port 8086): Payment processing with Stripe integration

## Quick Start with Docker Compose

Start all services including Phases 2, 3, and 4:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs for specific service
docker-compose logs -f order-service
docker-compose logs -f payment-service

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
- **PostgreSQL Order DB**: localhost:5437
- **PostgreSQL Payment DB**: localhost:5438
- **Redis**: localhost:6379

## Manual Setup

### 1. Set Up Databases

```bash
# Create databases
createdb order_db
createdb payment_db

# Or using psql
psql -U postgres -c "CREATE DATABASE order_db;"
psql -U postgres -c "CREATE DATABASE payment_db;"
```

### 2. Run Migrations

**Order Service:**
```bash
cd order-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/order_db?sslmode=disable" up
```

**Payment Service:**
```bash
cd payment-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/payment_db?sslmode=disable" up
```

### 3. Configure Stripe

1. Sign up at [stripe.com](https://stripe.com)
2. Get test API keys from Dashboard → Developers → API keys
3. Set up webhook endpoint:
   - URL: `http://your-domain.com/api/v1/webhooks/stripe`
   - Events: `payment_intent.succeeded`, `payment_intent.payment_failed`, `charge.refunded`
4. Copy webhook signing secret

### 4. Configure Environment

**Order Service:**
```bash
cd order-service
cp .env.example .env
# Edit .env with your configuration
```

**Payment Service:**
```bash
cd payment-service
cp .env.example .env
# Edit .env with your Stripe keys
```

### 5. Install Dependencies

```bash
# Order Service
cd order-service
go mod download

# Payment Service
cd payment-service
go mod download
```

### 6. Run Services

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

**Terminal 5 - Order Service:**
```bash
cd order-service
go run cmd/main.go
```

**Terminal 6 - Payment Service:**
```bash
cd payment-service
go run cmd/main.go
```

## Complete Order Flow

### Step 1: User Registration and Login

```bash
# Register
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "Password123!",
    "firstName": "John",
    "lastName": "Doe"
  }'

# Login (save the accessToken)
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "Password123!"
  }'
```

### Step 2: Add Products to Cart

```bash
# Add product to cart
curl -X POST http://localhost:8083/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "productId": "PRODUCT_UUID",
    "quantity": 2
  }'

# View cart
curl -X GET http://localhost:8083/api/v1/cart \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Step 3: Create Order

```bash
# Create order from cart
curl -X POST http://localhost:8085/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
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
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "US",
      "phone": "+1234567890"
    }
  }'
```

**Response includes order ID - save it for payment:**
```json
{
  "status": 201,
  "message": "Order created successfully",
  "data": {
    "id": "ORDER_UUID",
    "orderNumber": "ORD-2026-000001",
    "total": 2709.98,
    ...
  }
}
```

### Step 4: Create Payment Intent

```bash
# Create payment intent
curl -X POST http://localhost:8086/api/v1/payments/intent \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "orderId": "ORDER_UUID",
    "amount": 2709.98
  }'
```

**Response includes clientSecret for Stripe Elements:**
```json
{
  "status": 201,
  "message": "Payment intent created successfully",
  "data": {
    "paymentId": "PAYMENT_UUID",
    "clientSecret": "pi_xxx_secret_yyy",
    "amount": 2709.98,
    "currency": "usd"
  }
}
```

### Step 5: Process Payment (Frontend)

Use the clientSecret with Stripe Elements in your frontend to collect payment:

```javascript
const stripe = Stripe('YOUR_PUBLISHABLE_KEY');
const {error} = await stripe.confirmPayment({
  elements,
  clientSecret: 'pi_xxx_secret_yyy',
  confirmParams: {
    return_url: 'https://your-site.com/order-complete'
  }
});
```

### Step 6: Webhook Updates Payment Status

Stripe sends webhook events to update payment status automatically.

### Step 7: Check Order Status

```bash
# Get order details
curl -X GET http://localhost:8085/api/v1/orders/ORDER_UUID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Track order
curl -X GET http://localhost:8085/api/v1/orders/ORDER_UUID/track \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get payment status
curl -X GET http://localhost:8086/api/v1/payments/PAYMENT_UUID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Order Service Endpoints

### Customer Endpoints

**Create Order:**
```bash
POST /api/v1/orders
```

**Get Orders (with pagination):**
```bash
GET /api/v1/orders?page=1&limit=10&status=confirmed
```

**Get Order by ID:**
```bash
GET /api/v1/orders/:id
```

**Cancel Order:**
```bash
PUT /api/v1/orders/:id/cancel
```

**Track Order:**
```bash
GET /api/v1/orders/:id/track
```

### Admin Endpoints

**Update Order Status:**
```bash
PUT /api/v1/orders/:id/status
Content-Type: application/json

{
  "status": "confirmed",
  "comment": "Payment confirmed"
}
```

## Payment Service Endpoints

### Customer Endpoints

**Create Payment Intent:**
```bash
POST /api/v1/payments/intent
```

**Get Payment by ID:**
```bash
GET /api/v1/payments/:id
```

**Get Payment by Order ID:**
```bash
GET /api/v1/payments/order/:orderId
```

### Admin Endpoints

**Process Refund:**
```bash
POST /api/v1/payments/:id/refund
Content-Type: application/json

{
  "amount": 2709.98,
  "reason": "Customer request"
}
```

### Webhook Endpoint

**Stripe Webhook:**
```bash
POST /api/v1/webhooks/stripe
```

## Order State Machine

```
pending → confirmed → processing → shipped → delivered
  ↓           ↓            ↓
cancelled  cancelled  cancelled
```

**Valid Transitions:**
- `pending` → `confirmed`, `cancelled`
- `confirmed` → `processing`, `cancelled`
- `processing` → `shipped`, `cancelled`
- `shipped` → `delivered`
- `delivered` → (terminal state)
- `cancelled` → (terminal state)

## Payment Flow

1. **Create Payment Intent**: Order Service creates order with pending status
2. **Client Processing**: Frontend collects payment via Stripe Elements
3. **Webhook Events**: Stripe sends events to Payment Service
4. **Status Update**: Payment Service updates payment status
5. **Order Confirmation**: Order status updated to confirmed on successful payment

## Testing with Stripe

### Test Cards

- **Success**: `4242 4242 4242 4242`
- **Decline**: `4000 0000 0000 0002`
- **3D Secure**: `4000 0025 0000 3155`

Use any future expiry date, any 3-digit CVC, and any ZIP code.

### Local Webhook Testing

Use Stripe CLI to forward webhooks to your local server:

```bash
# Install Stripe CLI
brew install stripe/stripe-cli/stripe

# Login to Stripe
stripe login

# Forward webhooks
stripe listen --forward-to localhost:8086/api/v1/webhooks/stripe
```

## Architecture

### Service Communication

```
Order Service --> Cart Service (Get cart, Clear cart)
Payment Service --> Stripe API (Payment intents, Refunds)
Stripe --> Payment Service (Webhooks for payment events)
```

### Data Flow

1. **Order Creation Flow**:
   - User sends order request to Order Service
   - Order Service fetches cart from Cart Service
   - Order Service creates order with items snapshot
   - Order Service clears cart
   - Order created with pending status

2. **Payment Flow**:
   - User creates payment intent via Payment Service
   - Payment Service creates Stripe payment intent
   - Frontend processes payment with Stripe Elements
   - Stripe sends webhook to Payment Service
   - Payment Service updates payment status

## Database Schema

### Order Service (order_db)

**orders table:**
- id, user_id, order_number (unique), status
- subtotal, tax, shipping_cost, total
- shipping_address (JSONB), billing_address (JSONB)
- created_at, updated_at

**order_items table:**
- id, order_id (FK), product_id
- product_name, product_price, quantity, subtotal

**order_status_history table:**
- id, order_id (FK), status, comment, created_at

### Payment Service (payment_db)

**payments table:**
- id, order_id (unique), user_id, stripe_payment_intent_id (unique)
- amount, currency, status, payment_method
- created_at, updated_at

**refunds table:**
- id, payment_id (FK), stripe_refund_id (unique)
- amount, reason, status, created_at

## Health Checks

```bash
# Order Service
curl http://localhost:8085/health

# Payment Service
curl http://localhost:8086/health
```

## Troubleshooting

### Service Communication Issues

1. Ensure all services are running:
```bash
docker-compose ps
```

2. Check service logs:
```bash
docker-compose logs order-service
docker-compose logs payment-service
```

3. Verify network connectivity:
```bash
docker-compose exec order-service ping cart-service
```

### Stripe Issues

1. Verify Stripe keys are set:
```bash
docker-compose exec payment-service printenv | grep STRIPE
```

2. Check webhook signature:
```bash
# Make sure STRIPE_WEBHOOK_SECRET is set correctly
```

3. Test webhook locally:
```bash
stripe listen --forward-to localhost:8086/api/v1/webhooks/stripe
```

### Database Issues

1. Check database connectivity:
```bash
psql -U postgres -h localhost -p 5437 -d order_db
psql -U postgres -h localhost -p 5438 -d payment_db
```

2. View tables:
```sql
\dt
```

3. Check order status history:
```sql
SELECT * FROM order_status_history ORDER BY created_at DESC;
```

## Business Logic

### Order Creation
- Validates shipping and billing addresses
- Fetches cart from Cart Service
- Generates unique order number (ORD-YYYY-NNNNNN)
- Calculates tax (8%) and adds shipping ($10)
- Creates order with items snapshot
- Records initial status in history
- Clears cart after successful creation

### Payment Processing
- Creates Stripe payment intent with order metadata
- Returns client secret for frontend processing
- Handles webhook events for status updates
- Supports refunds with amount validation

### State Management
- Enforces strict state transitions
- Records all status changes in history
- Prevents invalid transitions

## Next Steps

Phase 5 will add:
- Inventory Service for stock management
- Notification Service for email/SMS alerts
- Stock reservation and low stock alerts
- Order confirmation emails

## Resources

- [Order Service Documentation](order-service/README.md)
- [Payment Service Documentation](payment-service/README.md)
- [Stripe Documentation](https://stripe.com/docs)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
