# E-Commerce Platform - Production-Grade Microservices Architecture

A complete enterprise-level e-commerce platform built with microservices architecture, featuring full observability, event-driven communication, and cloud-native deployment.

## Technology Stack

### Backend
- **Language**: Golang
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: Apache Kafka

### Frontend
- **Framework**: React.js with Vite
- **Language**: TypeScript
- **State Management**: Redux Toolkit
- **Routing**: React Router
- **HTTP Client**: Axios
- **UI Library**: Material UI

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes (KIND for local)
- **Monitoring**: Prometheus + Grafana
- **Logging**: Loki
- **Tracing**: OpenTelemetry
- **CI/CD**: GitHub Actions

## Architecture

This platform consists of 10 microservices:

1. **API Gateway** (Port 8080) - Request routing, JWT validation, rate limiting
2. **Auth Service** (Port 8081) - Authentication, authorization, RBAC
3. **Product Service** (Port 8082) - Product catalog, reviews, ratings
4. **Cart Service** (Port 8083) - Shopping cart management
5. **Wishlist Service** (Port 8084) - User wishlist management
6. **Order Service** (Port 8085) - Order creation and tracking
7. **Payment Service** (Port 8086) - Payment processing via Stripe
8. **Inventory Service** (Port 8087) - Stock management, low stock alerts
9. **Notification Service** (Port 8088) - Email and event notifications
10. **Reporting Service** (Port 8089) - Analytics and reporting

## Getting Started

See [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) for setup instructions.

## Documentation

- [Architecture Overview](docs/architecture/HIGH_LEVEL_ARCHITECTURE.md)
- [Service Documentation](docs/architecture/SERVICE_ARCHITECTURE.md)
- [Database Design](docs/database/ER_DIAGRAMS.md)
- [API Contracts](docs/api/API_CONTRACTS.md)
- [Event Contracts](docs/events/KAFKA_EVENTS.md)
- [Development Roadmap](docs/ROADMAP.md)

## Project Structure

```
e-commerce/
├── frontend/                 # React frontend application
├── api-gateway/             # API Gateway service
├── auth-service/            # Authentication service
├── product-service/         # Product catalog service
├── cart-service/            # Shopping cart service
├── wishlist-service/        # Wishlist service
├── order-service/           # Order management service
├── payment-service/         # Payment processing service
├── inventory-service/       # Inventory management service
├── notification-service/    # Notification service
├── reporting-service/       # Analytics and reporting service
├── infrastructure/          # Infrastructure as Code
│   ├── docker/             # Docker configurations
│   ├── kind/               # KIND cluster configurations
│   ├── kubernetes/         # Kubernetes manifests
│   ├── monitoring/         # Prometheus, Grafana, Loki
│   ├── kafka/              # Kafka configurations
│   ├── redis/              # Redis configurations
│   └── github-actions/     # CI/CD workflows
├── docs/                    # Documentation
└── scripts/                 # Utility scripts
```

## License

MIT
