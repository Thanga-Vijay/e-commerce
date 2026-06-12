# Notification Service

**Port:** 8088  
**Language:** Golang  
**Framework:** Gin  
**Message Queue:** Kafka Consumer

## Responsibilities

- **Email Notifications** - Send transactional emails
- **Event Processing** - Consume Kafka events
- **Template Management** - Email templates
- **Notification History** - Track sent notifications
- **Retry Logic** - Retry failed notifications

## Email Templates

### User Events
- Welcome Email (on user registration)
- Email Verification
- Password Reset

### Order Events
- Order Confirmation
- Order Shipped
- Order Delivered

### Inventory Events
- Low Stock Alert (to admin)
- Out of Stock Alert (to admin)

### Payment Events
- Payment Confirmation
- Payment Failed
- Refund Processed

## Events Consumed

- `user.created` → Send welcome email
- `user.email_verification` → Send verification email
- `user.password_reset` → Send reset link
- `order.created` → Send order confirmation
- `order.shipped` → Send shipping notification
- `order.delivered` → Send delivery confirmation
- `payment.completed` → Send payment receipt
- `payment.failed` → Send payment failure notice
- `payment.refunded` → Send refund confirmation
- `inventory.lowstock` → Alert admin
- `inventory.outofstock` → Alert admin

## Database Schema

### notifications
- id (UUID, PK)
- user_id
- type (email, sms, push)
- template
- recipient
- subject
- status (pending, sent, failed)
- retry_count
- error_message
- sent_at
- created_at

## SMTP Configuration

Environment variables:
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USER` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `SMTP_FROM` - From email address

## Retry Mechanism

- Failed emails are retried 3 times
- Exponential backoff: 1m, 5m, 15m
- After 3 failures, mark as failed and alert admin

## API Endpoints

- `POST /api/v1/notifications/send` - Send notification (internal)
- `GET /api/v1/notifications/:userId` - Get user notifications
- `GET /api/v1/notifications/:id/status` - Get notification status

## Directory Structure

```
notification-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── consumers/
│   ├── templates/
│   └── mailer/
├── templates/
│   ├── welcome.html
│   ├── order-confirmation.html
│   └── ...
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```

## Sample Email Template

```html
<!-- templates/order-confirmation.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Order Confirmation</title>
</head>
<body>
    <h1>Thank you for your order!</h1>
    <p>Order Number: {{.OrderNumber}}</p>
    <p>Total: ${{.Total}}</p>
    <!-- More details -->
</body>
</html>
```
