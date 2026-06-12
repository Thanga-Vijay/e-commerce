# Kafka Configuration

Apache Kafka setup for event-driven architecture.

## Components

- Kafka brokers (3 replicas)
- Zookeeper ensemble (3 replicas)
- Schema Registry
- Kafka UI (for management)

## Topics

All topics created with:
- Partitions: 3
- Replication Factor: 2
- Retention: 7 days

### User Events
- `user.created`
- `user.updated`
- `user.deleted`

### Product Events
- `product.created`
- `product.updated`
- `product.deleted`

### Inventory Events
- `inventory.updated`
- `inventory.lowstock`
- `inventory.outofstock`

### Order Events
- `order.created`
- `order.confirmed`
- `order.shipped`
- `order.delivered`
- `order.cancelled`

### Payment Events
- `payment.created`
- `payment.completed`
- `payment.failed`
- `payment.refunded`

### Notification Events
- `notification.send`
- `notification.sent`
- `notification.failed`

### Dead Letter Queue
- `dlq.notifications`
- `dlq.inventory`
- `dlq.orders`

## Setup

Create topics:
```bash
./create-topics.sh
```

## Access

Kafka UI:
```bash
kubectl port-forward svc/kafka-ui 8080:8080 -n kafka
```
URL: http://localhost:8080

## Event Schema Example

```json
{
  "eventId": "uuid",
  "eventType": "order.created",
  "timestamp": "2026-06-12T10:30:00Z",
  "version": "1.0",
  "payload": {
    "orderId": "uuid",
    "userId": "uuid",
    "total": 100.00,
    "items": []
  }
}
```

## Consumer Groups

- `notification-service-group`
- `inventory-service-group`
- `reporting-service-group`
- `order-service-group`

## Monitoring

Metrics exposed:
- Consumer lag
- Message throughput
- Broker health
- Partition distribution
