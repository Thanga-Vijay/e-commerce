# E-Commerce Platform - Phase 5 Implementation

This document provides instructions for running the Inventory Service and Notification Service (Phase 5).

## Services Overview

**Phase 5** adds two new microservices:
- **Inventory Service** (Port 8087): Stock management with warehouses and transactions
- **Notification Service** (Port 8088): Email notifications with templates and retry logic

## Quick Start with Docker Compose

Start all services including Phases 2, 3, 4, and 5:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs for specific service
docker-compose logs -f inventory-service
docker-compose logs -f notification-service

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
- **PostgreSQL Inventory DB**: localhost:5439
- **PostgreSQL Notification DB**: localhost:5440
- **Redis**: localhost:6379

## Manual Setup

### 1. Set Up Databases

```bash
# Create databases
createdb inventory_db
createdb notification_db

# Or using psql
psql -U postgres -c "CREATE DATABASE inventory_db;"
psql -U postgres -c "CREATE DATABASE notification_db;"
```

### 2. Run Migrations

**Inventory Service:**
```bash
cd inventory-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/inventory_db?sslmode=disable" up
```

**Notification Service:**
```bash
cd notification-service
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/notification_db?sslmode=disable" up
```

### 3. Configure SMTP for Notifications

For Gmail:
1. Enable 2-factor authentication on your Google account
2. Generate an App Password: https://myaccount.google.com/apppasswords
3. Use the app password in your `.env` file

For other providers (SendGrid, AWS SES, Mailgun):
- Update SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD accordingly

### 4. Configure Environment

**Inventory Service:**
```bash
cd inventory-service
cp .env.example .env
# Edit .env with your configuration
```

**Notification Service:**
```bash
cd notification-service
cp .env.example .env
# Edit .env with your SMTP configuration
```

### 5. Install Dependencies

```bash
# Inventory Service
cd inventory-service
go mod download

# Notification Service
cd notification-service
go mod download
```

### 6. Run Services

**Terminal 7 - Inventory Service:**
```bash
cd inventory-service
go run cmd/main.go
```

**Terminal 8 - Notification Service:**
```bash
cd notification-service
go run cmd/main.go
```

## Inventory Service API

### Admin Endpoints

**Create Inventory:**
```bash
POST /api/v1/inventory
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "productId": "PRODUCT_UUID",
  "sku": "SKU-12345",
  "quantityAvailable": 100,
  "warehouseId": "WAREHOUSE_UUID",
  "lowStockThreshold": 10
}
```

**Get All Inventory (Paginated):**
```bash
GET /api/v1/inventory?page=1&limit=10
Authorization: Bearer ADMIN_TOKEN
```

**Get Low Stock Items:**
```bash
GET /api/v1/inventory/low-stock
Authorization: Bearer ADMIN_TOKEN
```

**Update Inventory:**
```bash
PUT /api/v1/inventory/:id
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "quantityAvailable": 150,
  "lowStockThreshold": 15
}
```

**Reserve Stock:**
```bash
POST /api/v1/inventory/product/:productId/reserve
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "quantity": 5,
  "referenceId": "ORDER_UUID",
  "notes": "Reserved for order ORD-2026-001234"
}
```

**Release Stock:**
```bash
POST /api/v1/inventory/product/:productId/release
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "quantity": 2,
  "referenceId": "ORDER_UUID",
  "notes": "Order cancelled"
}
```

**Confirm Sale:**
```bash
POST /api/v1/inventory/product/:productId/confirm-sale
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "quantity": 3,
  "referenceId": "ORDER_UUID",
  "notes": "Order completed"
}
```

**Adjust Stock:**
```bash
POST /api/v1/inventory/product/:productId/adjust
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "quantity": -5,
  "notes": "Damaged goods removed"
}
```

### Public Endpoints

**Get Inventory by Product ID:**
```bash
GET /api/v1/inventory/product/:productId
Authorization: Bearer YOUR_ACCESS_TOKEN
```

**Get Inventory by SKU:**
```bash
GET /api/v1/inventory/sku/:sku
Authorization: Bearer YOUR_ACCESS_TOKEN
```

**Get Transaction History:**
```bash
GET /api/v1/inventory/product/:productId/transactions?limit=50
Authorization: Bearer YOUR_ACCESS_TOKEN
```

## Notification Service API

### Send Notification

**Welcome Email:**
```bash
POST /api/v1/notifications
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "userId": "USER_UUID",
  "type": "email",
  "template": "welcome",
  "recipient": "user@example.com",
  "data": {
    "Name": "John Doe"
  }
}
```

**Email Verification:**
```bash
POST /api/v1/notifications
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "userId": "USER_UUID",
  "type": "email",
  "template": "email_verification",
  "recipient": "user@example.com",
  "data": {
    "Name": "John Doe",
    "VerificationLink": "https://example.com/verify?token=abc123"
  }
}
```

**Password Reset:**
```bash
POST /api/v1/notifications
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "userId": "USER_UUID",
  "type": "email",
  "template": "password_reset",
  "recipient": "user@example.com",
  "data": {
    "Name": "John Doe",
    "ResetLink": "https://example.com/reset?token=xyz789"
  }
}
```

**Order Confirmation:**
```bash
POST /api/v1/notifications
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "userId": "USER_UUID",
  "type": "email",
  "template": "order_confirmation",
  "recipient": "user@example.com",
  "data": {
    "Name": "John Doe",
    "OrderNumber": "ORD-2026-001234",
    "Total": 2709.98,
    "Items": [
      {
        "Name": "MacBook Pro 16",
        "Quantity": 1,
        "Price": 2499.99
      }
    ]
  }
}
```

**Order Shipped:**
```bash
POST /api/v1/notifications
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json

{
  "userId": "USER_UUID",
  "type": "email",
  "template": "order_shipped",
  "recipient": "user@example.com",
  "data": {
    "Name": "John Doe",
    "OrderNumber": "ORD-2026-001234",
    "TrackingNumber": "1Z999AA1012345678",
    "Carrier": "UPS"
  }
}
```

**Low Stock Alert (Admin):**
```bash
POST /api/v1/notifications
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "type": "email",
  "template": "low_stock_alert",
  "recipient": "admin@example.com",
  "data": {
    "ProductName": "MacBook Pro 16",
    "SKU": "SKU-12345",
    "Quantity": 5,
    "Threshold": 10
  }
}
```

### Get Notifications

**Get My Notifications:**
```bash
GET /api/v1/notifications/my?page=1&limit=10
Authorization: Bearer YOUR_ACCESS_TOKEN
```

**Get Notification by ID:**
```bash
GET /api/v1/notifications/:id
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Admin Endpoints

**Process Pending Notifications:**
```bash
POST /api/v1/notifications/process-pending?limit=10
Authorization: Bearer ADMIN_TOKEN
```

**Retry Failed Notifications:**
```bash
POST /api/v1/notifications/retry-failed?limit=10
Authorization: Bearer ADMIN_TOKEN
```

## Available Email Templates

1. **welcome** - Welcome new users
2. **email_verification** - Verify email addresses
3. **password_reset** - Password reset instructions
4. **order_confirmation** - Order confirmation with items
5. **payment_receipt** - Payment received confirmation
6. **order_shipped** - Shipping notification with tracking
7. **order_delivered** - Delivery confirmation
8. **low_stock_alert** - Low stock alert for admins

## Inventory Features

### Stock Management
- Create inventory records for products
- Track quantity available, reserved, and sold
- Set low stock thresholds for alerts
- Associate inventory with warehouses

### Transaction Types
- **purchase** - Initial inventory or restocking
- **sale** - Confirmed sale
- **adjustment** - Manual adjustments (damage, returns, etc.)
- **return** - Customer returns
- **reserve** - Reserve stock for orders
- **release** - Release reserved stock (cancelled orders)

### Workflow
1. **Reserve Stock**: When order created, reserve stock
2. **Confirm Sale**: When payment succeeds, confirm sale
3. **Release Stock**: If order cancelled, release reserved stock

### Transaction History
- Complete audit trail of all stock movements
- Track previous and new quantities
- Reference IDs link to orders
- Notes provide context for adjustments

## Notification Features

### Email System
- HTML email templates with responsive design
- SMTP integration (Gmail, SendGrid, AWS SES, etc.)
- Template rendering with dynamic data
- Automatic retry with exponential backoff

### Retry Mechanism
- Failed notifications automatically retry
- Maximum 3 retry attempts
- Exponential backoff: 1min, 2min, 4min
- Status tracking: pending, sent, failed

### Status Tracking
- Track notification delivery status
- Error messages for failed sends
- Sent timestamps
- User notification history

## Architecture

### Service Communication

```
Order Service --> Inventory Service (Reserve/Release/Confirm stock)
All Services --> Notification Service (Send emails)
```

### Data Flow

**Order Creation with Inventory:**
1. User creates order
2. Order Service reserves stock via Inventory Service
3. If stock insufficient, order fails
4. Payment processed
5. On success: Order Service confirms sale via Inventory Service
6. On failure: Order Service releases stock via Inventory Service

**Notification Flow:**
1. Service creates notification request
2. Notification Service renders template
3. Notification Service sends via SMTP
4. On success: Mark as sent
5. On failure: Mark as failed, schedule retry
6. Background job retries failed notifications

## Database Schemas

### Inventory Service (inventory_db)

**warehouses table:**
- id, name, code (unique), address (JSONB), is_active
- Default warehouse created on migration

**inventory table:**
- id, product_id (unique), sku
- quantity_available, quantity_reserved, quantity_sold
- warehouse_id (FK), low_stock_threshold
- created_at, updated_at

**inventory_transactions table:**
- id, inventory_id (FK), transaction_type
- quantity, previous_qty, new_qty
- reference_id, notes, created_at

### Notification Service (notification_db)

**notifications table:**
- id, user_id, type, template, recipient
- subject, body (HTML), status
- retry_count, error_message, sent_at
- created_at

## Health Checks

```bash
# Inventory Service
curl http://localhost:8087/health

# Notification Service
curl http://localhost:8088/health
```

## Troubleshooting

### Inventory Service Issues

1. Check database connectivity:
```bash
psql -U postgres -h localhost -p 5439 -d inventory_db
```

2. View inventory records:
```sql
SELECT * FROM inventory;
```

3. Check transaction history:
```sql
SELECT * FROM inventory_transactions ORDER BY created_at DESC LIMIT 10;
```

### Notification Service Issues

1. **SMTP Authentication Failed:**
   - Verify SMTP credentials are correct
   - For Gmail, use App Password, not regular password
   - Check 2FA is enabled if using Gmail

2. **Emails Not Sending:**
   - Check SMTP_HOST and SMTP_PORT are correct
   - Test SMTP connection manually
   - Check firewall rules allow outbound SMTP

3. **Check failed notifications:**
```sql
SELECT * FROM notifications WHERE status = 'failed' ORDER BY created_at DESC;
```

4. **View notification history:**
```sql
SELECT id, recipient, template, status, retry_count, created_at 
FROM notifications 
ORDER BY created_at DESC 
LIMIT 20;
```

### SMTP Configuration Examples

**Gmail:**
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**SendGrid:**
```
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

**AWS SES:**
```
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your-ses-smtp-username
SMTP_PASSWORD=your-ses-smtp-password
```

## Use Cases

### Inventory Management
1. Admin creates inventory for new products
2. System reserves stock when orders placed
3. System confirms sale when payment succeeds
4. System releases stock if order cancelled
5. Admin adjusts stock for damages/returns
6. Low stock alerts sent to admins

### Notification System
1. User registers → Send welcome email
2. User forgets password → Send reset link
3. Order created → Send confirmation email
4. Payment received → Send receipt
5. Order shipped → Send tracking info
6. Order delivered → Send delivery confirmation
7. Low stock detected → Alert admins

## Next Steps

Phase 6 will add:
- Reporting Service for analytics and metrics
- Dashboard with revenue, orders, customer stats
- Data aggregation jobs
- Export functionality (CSV, PDF)

## Resources

- [Inventory Service Documentation](inventory-service/README.md)
- [Notification Service Documentation](notification-service/README.md)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [SMTP Configuration Guide](https://support.google.com/mail/answer/7126229)
