# Payment Service

**Port:** 8086  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**Payment Gateway:** Stripe  
**ORM:** GORM

## Responsibilities

- **Payment Processing** - Process payments via Stripe
- **Payment Intent Creation** - Create Stripe payment intents
- **Refund Management** - Process refunds
- **Payment Status Tracking** - Track payment status
- **Webhook Handling** - Handle Stripe webhooks

## Database Schema

### payments
- id (UUID, PK)
- order_id (unique)
- user_id
- stripe_payment_intent_id
- amount
- currency
- status (pending, succeeded, failed, refunded)
- payment_method
- created_at
- updated_at

### refunds
- id (UUID, PK)
- payment_id (FK)
- stripe_refund_id
- amount
- reason
- status
- created_at

## API Endpoints

- `POST /api/v1/payments/intent` - Create payment intent
- `POST /api/v1/payments/confirm` - Confirm payment
- `GET /api/v1/payments/:id` - Get payment details
- `POST /api/v1/payments/:id/refund` - Process refund (admin)
- `POST /api/v1/payments/webhook` - Stripe webhook

## Payment Flow

1. Client creates order
2. Payment service creates Stripe payment intent
3. Client confirms payment on frontend
4. Stripe webhook notifies payment service
5. Payment service publishes `payment.completed` event
6. Order service confirms order

## Events Published

- `payment.created`
- `payment.completed`
- `payment.failed`
- `payment.refunded`

## Stripe Integration

```go
stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

// Create payment intent
params := &stripe.PaymentIntentParams{
    Amount:   stripe.Int64(amount),
    Currency: stripe.String("usd"),
}
pi, _ := paymentintent.New(params)
```

## Webhook Verification

```go
event, err := webhook.ConstructEvent(
    payload, 
    signature, 
    webhookSecret,
)
```

## Directory Structure

```
payment-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── events/
│   └── stripe/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```
