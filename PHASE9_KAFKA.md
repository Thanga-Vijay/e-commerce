# E-Commerce Platform - Phase 9: Event-Driven Architecture with Kafka

This document provides comprehensive instructions for implementing event-driven architecture with Apache Kafka (Phase 9).

## Overview

Phase 9 implements a complete event-driven architecture using Apache Kafka for asynchronous communication between microservices. This enables:
- **Loose coupling** between services
- **Scalability** through asynchronous processing
- **Reliability** with event replay and dead letter queues
- **Real-time** notifications and updates
- **Event sourcing** capabilities

## Architecture

### Event Flow Diagram

```
┌──────────────┐
│ Auth Service │──────► user.registered
└──────────────┘        user.updated
                        user.deleted

┌──────────────┐
│Order Service │──────► order.created ──────► Notification Service
└──────────────┘        order.confirmed ────► Inventory Service
                        order.shipped ───────► Reporting Service
                        order.delivered
                        order.cancelled

┌──────────────────┐
│ Payment Service  │──► payment.completed ──► Order Service
└──────────────────┘    payment.failed
                        payment.refunded

┌────────────────────┐
│ Inventory Service  │► inventory.updated
└────────────────────┘  inventory.low-stock
                        inventory.out-of-stock

┌────────────────────────┐
│ Notification Service   │► notification.email.sent
└────────────────────────┘  notification.email.failed
```

### Kafka Infrastructure

```
┌─────────────────────────────────────┐
│          Zookeeper (2181)           │
│    Coordination & Configuration     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Kafka Broker (9092)         │
│  - 3 Partitions per topic           │
│  - 7 days retention                 │
│  - Snappy compression               │
└──────────────┬──────────────────────┘
               │
        ┌──────┴──────┐
   ┌────▼────┐   ┌────▼────┐
   │ Topics  │   │   DLQ   │
   │  (26)   │   │ (1)     │
   └─────────┘   └─────────┘
```

## Topics

### User Events (3 topics)
- `user.registered` - New user registration
- `user.updated` - User profile updates
- `user.deleted` - User account deletion

### Product Events (3 topics)
- `product.created` - New product created
- `product.updated` - Product information updated
- `product.deleted` - Product removed from catalog

### Cart Events (3 topics)
- `cart.item.added` - Item added to cart
- `cart.item.removed` - Item removed from cart
- `cart.cleared` - Cart emptied

### Order Events (5 topics)
- `order.created` - New order placed
- `order.confirmed` - Payment confirmed, order processing
- `order.shipped` - Order shipped with tracking
- `order.delivered` - Order successfully delivered
- `order.cancelled` - Order cancelled with refund

### Payment Events (4 topics)
- `payment.initiated` - Payment process started
- `payment.completed` - Payment successfully processed
- `payment.failed` - Payment processing failed
- `payment.refunded` - Payment refunded

### Inventory Events (3 topics)
- `inventory.updated` - Stock quantity changed
- `inventory.low-stock` - Stock below threshold
- `inventory.out-of-stock` - Product out of stock

### Notification Events (3 topics)
- `notification.email.send` - Email notification request
- `notification.email.sent` - Email successfully sent
- `notification.email.failed` - Email sending failed

### Dead Letter Queue (1 topic)
- `dlq.events` - Failed messages for manual review

**Total: 26 topics + 1 DLQ**

## Quick Start

### Start Kafka Infrastructure

```bash
# Start Kafka and all dependencies
docker-compose -f docker-compose.yml -f docker-compose.kafka.yml up -d

# Check Kafka health
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# View topics in UI
open http://localhost:8090
```

### Verify Topics Created

```bash
# List all topics
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Describe specific topic
docker-compose exec kafka kafka-topics --describe --topic order.created --bootstrap-server localhost:9092
```

### Test Event Publishing

```bash
# Produce test event
docker-compose exec kafka kafka-console-producer --topic order.created --bootstrap-server localhost:9092
# Type message and press Ctrl+D

# Consume events
docker-compose exec kafka kafka-console-consumer --topic order.created --from-beginning --bootstrap-server localhost:9092
```

## Event Contracts

All events follow a standard envelope format:

```json
{
  "id": "uuid-v4",
  "type": "event.type",
  "source": "service-name",
  "timestamp": "2026-06-12T10:30:00Z",
  "version": "1.0",
  "data": {
    // Event-specific payload
  }
}
```

### Example: Order Created Event

```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "type": "order.created",
  "source": "order-service",
  "timestamp": "2026-06-12T10:30:00Z",
  "version": "1.0",
  "data": {
    "order_id": 123,
    "user_id": 456,
    "total_amount": 149.99,
    "status": "pending",
    "item_count": 3,
    "shipping_address": "123 Main St, City, State 12345"
  }
}
```

## Integration Guide

### 1. Add Kafka Dependency

Add to your service's `go.mod`:

```go
require (
    github.com/segmentio/kafka-go v0.4.47
    github.com/google/uuid v1.5.0
)
```

Run `go mod download` to install dependencies.

### 2. Initialize Kafka Producer

```go
import "github.com/ecommerce/events/kafka"

// In main.go or initialization
kafkaConfig := kafka.LoadConfig()
producer := kafka.NewProducer(kafkaConfig.Brokers, "your-service-name")
defer producer.Close()
```

### 3. Publish Events

```go
import (
    "context"
    "time"
    "github.com/ecommerce/events"
)

// Create event
event := events.OrderCreatedEvent{
    BaseEvent: events.BaseEvent{
        ID:        uuid.New().String(),
        Type:      events.EventOrderCreated,
        Source:    "order-service",
        Timestamp: time.Now().UTC(),
        Version:   "1.0",
    },
    OrderID:     order.ID,
    UserID:      userID,
    TotalAmount: totalAmount,
}

// Publish to Kafka
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := producer.PublishEvent(ctx, "order.created", events.EventOrderCreated, event)
```

### 4. Consume Events

```go
// Create consumer configuration
config := kafka.ConsumerConfig{
    Brokers: []string{"localhost:9093"},
    Topic:   "order.created",
    GroupID: "your-service-orders",
    Service: "your-service",
    Handler: handleOrderEvent,
}

// Create consumer
consumer := kafka.NewConsumer(config)

// Start consuming
ctx := context.Background()
go consumer.Start(ctx)

// Handler function
func handleOrderEvent(ctx context.Context, event map[string]interface{}) error {
    log.Printf("Processing order event: %v", event)
    // Your business logic here
    return nil
}
```

## Environment Variables

Add to your `.env` files:

```bash
# Kafka Configuration
KAFKA_BROKERS=kafka:9092
KAFKA_ENABLED=true

# For local development (outside Docker)
KAFKA_BROKERS=localhost:9093
```

## Dead Letter Queue (DLQ)

### Purpose
The DLQ captures messages that fail processing after multiple retry attempts, enabling:
- Manual inspection of failed messages
- Automatic retry with exponential backoff
- Audit trail of processing failures
- Prevention of message loss

### DLQ Message Format

```json
{
  "id": "uuid",
  "type": "message.failed",
  "source": "dlq-service",
  "timestamp": "2026-06-12T10:30:00Z",
  "version": "1.0",
  "data": {
    "original_topic": "order.created",
    "original_partition": 2,
    "original_offset": 12345,
    "original_message": "{ original event json }",
    "error_message": "Database connection timeout",
    "service": "inventory-service",
    "retry_count": 3,
    "failed_at": "2026-06-12T10:30:00Z"
  }
}
```

### DLQ Operations

```bash
# View DLQ messages
docker-compose exec kafka kafka-console-consumer --topic dlq.events --from-beginning --bootstrap-server localhost:9092

# Count DLQ messages
docker-compose exec kafka kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic dlq.events
```

### Automatic Retry Logic

- **Retry attempts**: 3 times
- **Retry interval**: 10 minutes
- **Final status**: `failed` after 3 attempts
- **Manual retry**: Available via API

## Monitoring Kafka

### Kafka UI

Access Kafka UI at: http://localhost:8090

Features:
- View all topics and partitions
- Monitor consumer lag
- Browse messages
- Create/delete topics
- View broker metrics

### Consumer Lag Monitoring

```bash
# Check consumer group lag
docker-compose exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group your-consumer-group

# List all consumer groups
docker-compose exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list
```

### Topic Metrics

```bash
# Topic size and message count
docker-compose exec kafka kafka-log-dirs --bootstrap-server localhost:9092 --describe

# Partition details
docker-compose exec kafka kafka-topics --describe --topic order.created --bootstrap-server localhost:9092
```

## Event Patterns

### 1. Event Notification Pattern

**Use Case**: Notify other services of state changes

```
Order Service ──► order.created ──► Notification Service
                                  ► Inventory Service
                                  ► Reporting Service
```

### 2. Event-Carried State Transfer

**Use Case**: Include full entity state in event

```json
{
  "type": "product.updated",
  "data": {
    "product_id": 123,
    "name": "Product Name",
    "price": 49.99,
    "stock": 100,
    "category": "Electronics",
    // Full product state
  }
}
```

### 3. Event Sourcing Pattern

**Use Case**: Store all state changes as events

```
order.created → order.confirmed → order.shipped → order.delivered
```

Rebuild state by replaying events.

### 4. CQRS Pattern

**Use Case**: Separate read and write models

```
Write: Order Service ──► order.created
                    ↓
Read:  Reporting Service ──► Materialized View
```

## Best Practices

### 1. Event Schema Versioning

```go
type BaseEvent struct {
    Version string `json:"version"` // Always include version
}

// Version 1.0
type OrderCreatedV1 struct {
    BaseEvent
    OrderID uint
}

// Version 2.0 (backward compatible)
type OrderCreatedV2 struct {
    BaseEvent
    OrderID     uint
    CustomerName string // New field
}
```

### 2. Idempotency

Ensure consumers can safely process the same message multiple times:

```go
func handleOrderEvent(ctx context.Context, event map[string]interface{}) error {
    orderID := event["order_id"]
    
    // Check if already processed
    var existing Order
    if err := db.Where("id = ?", orderID).First(&existing).Error; err == nil {
        log.Printf("Order %d already processed, skipping", orderID)
        return nil // Idempotent - already handled
    }
    
    // Process order...
}
```

### 3. Error Handling

```go
func handleEvent(ctx context.Context, event map[string]interface{}) error {
    // Transient errors - will be retried automatically
    if err := callExternalService(); err != nil {
        return fmt.Errorf("transient error: %w", err)
    }
    
    // Permanent errors - log and skip
    if invalidData() {
        log.Printf("Invalid data in event, skipping: %v", event)
        return nil // Don't retry
    }
    
    return nil
}
```

### 4. Event Ordering

Events in the same partition are guaranteed to be ordered:

```go
// Use consistent key for related events
message := kafka.Message{
    Topic: "order.events",
    Key:   []byte(fmt.Sprintf("order-%d", orderID)), // Same key = same partition
    Value: eventJSON,
}
```

### 5. Event Retention

Configure retention based on use case:

```bash
# 7 days (default)
kafka-configs --alter --add-config retention.ms=604800000 --topic order.created

# Infinite retention (event sourcing)
kafka-configs --alter --add-config retention.ms=-1 --topic event.store

# Compact logs (keep latest value per key)
kafka-configs --alter --add-config cleanup.policy=compact --topic product.state
```

## Troubleshooting

### Consumer Not Receiving Messages

```bash
# Check consumer group status
docker-compose exec kafka kafka-consumer-groups --describe --group your-group --bootstrap-server localhost:9092

# Reset consumer offset
docker-compose exec kafka kafka-consumer-groups --reset-offsets --to-earliest --execute --group your-group --topic your-topic --bootstrap-server localhost:9092
```

### High Consumer Lag

```bash
# Increase consumer instances (scale horizontally)
docker-compose up -d --scale your-service=3

# Check partition count
docker-compose exec kafka kafka-topics --describe --topic your-topic --bootstrap-server localhost:9092
```

### Topic Not Found

```bash
# Recreate topic
docker-compose exec kafka kafka-topics --create --topic your-topic --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092
```

### Connection Refused

```bash
# Check Kafka is running
docker-compose ps kafka

# Check Kafka logs
docker-compose logs kafka

# Verify network connectivity
docker-compose exec your-service ping kafka
```

## Performance Tuning

### Producer Settings

```go
writer := &kafka.Writer{
    Addr:         kafka.TCP("kafka:9092"),
    Balancer:     &kafka.LeastBytes{},    // Load balancing
    BatchSize:    100,                     // Batch messages
    BatchTimeout: 10 * time.Millisecond,  // Max wait time
    Compression:  kafka.Snappy,            // Compression
    RequiredAcks: kafka.RequireOne,        // Acknowledgment level
    Async:        true,                    // Non-blocking writes
}
```

### Consumer Settings

```go
reader := kafka.NewReader(kafka.ReaderConfig{
    Brokers:        []string{"kafka:9092"},
    Topic:          "your-topic",
    GroupID:        "your-group",
    MinBytes:       10e3,  // 10KB min fetch
    MaxBytes:       10e6,  // 10MB max fetch
    CommitInterval: time.Second,
    StartOffset:    kafka.LastOffset,
    MaxAttempts:    3,
})
```

## Testing

### Unit Testing with Mock Kafka

```go
// Mock producer for testing
type MockProducer struct{}

func (m *MockProducer) PublishEvent(ctx context.Context, topic string, eventType string, payload interface{}) error {
    // Record event for assertion
    return nil
}

// Use in tests
func TestOrderCreation(t *testing.T) {
    mockProducer := &MockProducer{}
    handler := &OrderHandler{
        KafkaProducer: mockProducer,
    }
    
    // Test order creation
    // Assert event was published
}
```

### Integration Testing

```bash
# Start test Kafka
docker-compose -f docker-compose.kafka.yml up -d

# Run integration tests
go test -tags=integration ./...

# Cleanup
docker-compose -f docker-compose.kafka.yml down -v
```

## Migration Guide

### Migrating Existing Services

1. **Add Kafka client library**
   ```bash
   go get github.com/segmentio/kafka-go
   ```

2. **Initialize producer in main.go**
   ```go
   producer := kafka.NewProducer(brokers, "service-name")
   defer producer.Close()
   ```

3. **Publish events after state changes**
   ```go
   // After creating order
   producer.PublishEvent(ctx, "order.created", eventType, orderData)
   ```

4. **Add consumers for dependent events**
   ```go
   consumer := kafka.NewConsumer(config)
   go consumer.Start(ctx)
   ```

5. **Test thoroughly**
   - Verify events are published
   - Verify consumers process events
   - Test failure scenarios

## Next Steps

### Phase 10: Kubernetes Deployment
- Kubernetes manifests for all services
- StatefulSets for Kafka
- ConfigMaps and Secrets
- Ingress configuration
- Auto-scaling policies

### Phase 11: Advanced Monitoring
- Distributed tracing with Jaeger
- Kafka metrics in Prometheus
- Custom Grafana dashboards
- Alert rules for Kafka lag

### Phase 12: CI/CD Pipeline
- Automated testing
- Container registry integration
- Blue-green deployments
- Canary releases

## Resources

- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Kafka-go Client](https://github.com/segmentio/kafka-go)
- [Event-Driven Architecture Patterns](https://martinfowler.com/articles/201701-event-driven.html)
- [Kafka Best Practices](https://docs.confluent.io/platform/current/kafka/deployment.html)

## Summary

Phase 9 delivers a production-ready event-driven architecture with:
- ✅ 26 event topics + 1 DLQ
- ✅ Kafka infrastructure with Zookeeper
- ✅ Producer and consumer utilities
- ✅ Event contracts and type definitions
- ✅ Dead letter queue with automatic retry
- ✅ Kafka UI for monitoring
- ✅ Integration examples for all services
- ✅ Comprehensive documentation

Your e-commerce platform now supports asynchronous, scalable, event-driven communication!
