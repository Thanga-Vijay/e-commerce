# E-Commerce Platform - Production-Grade Microservices Architecture

A complete enterprise-level e-commerce platform built with microservices architecture, featuring full observability, event-driven communication, and cloud-native deployment.

## 🎯 Project Status

**Current Phase**: Phase 8 Complete - Advanced Dockerization & Container Orchestration ✅

### Completed Phases
- ✅ Phase 1: Architecture Design & Database Schema
- ✅ Phase 2: Auth & Product Services
- ✅ Phase 3: Cart & Wishlist Services
- ✅ Phase 4: Order & Payment Services (Stripe Integration)
- ✅ Phase 5: Inventory & Notification Services
- ✅ Phase 6: Reporting Service with Analytics
- ✅ Phase 7: React Frontend (Complete UI)
- ✅ **Phase 8: Advanced Dockerization & Container Orchestration**

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: PostgreSQL 15 (9 separate databases)
- **Cache**: Redis 7 with password authentication
- **Message Queue**: Apache Kafka (Phase 9)

### Frontend
- **Framework**: React 18 with Vite
- **Language**: JavaScript/JSX
- **Styling**: Tailwind CSS 3.3
- **Routing**: React Router v6
- **HTTP Client**: Axios with interceptors
- **State Management**: Context API

### Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose (dev/staging/prod)
- **Reverse Proxy**: Nginx with rate limiting
- **Monitoring**: Prometheus + Grafana + Alertmanager
- **Logging**: JSON structured logging with rotation
- **Health Checks**: Comprehensive health monitoring
- **CI/CD**: GitHub Actions (Phase 12)

## Architecture

This platform consists of 9 microservices + frontend:

1. **Auth Service** (Port 8081) - JWT authentication, user management, RBAC
2. **Product Service** (Port 8082) - Product catalog, categories, Redis caching
3. **Cart Service** (Port 8083) - Shopping cart with service-to-service calls
4. **Wishlist Service** (Port 8084) - Wishlist management
5. **Order Service** (Port 8085) - Order processing, state machine workflow
6. **Payment Service** (Port 8086) - Stripe integration, webhooks, refunds
7. **Inventory Service** (Port 8087) - Stock management, warehouses
8. **Notification Service** (Port 8088) - Email notifications with SMTP
9. **Reporting Service** (Port 8089) - Dashboard analytics, CSV/PDF reports
10. **Frontend** (Port 3000) - React SPA with full e-commerce UI

## Quick Start

### Development Environment

```bash
# Using Make (recommended)
### Phase Documentation
- [Phase 1: Architecture](PHASE1_ARCHITECTURE.md) - System design, database schemas, API contracts
- [Phase 8: Docker & Orchestration](PHASE8_SETUP.md) - Container setup, deployment, monitoring

### Service Ports
- Frontend: http://localhost:3000
- Auth Service: http://localhost:8081
- Product Service: http://localhost:8082
- Cart Service: http://localhost:8083
- Wishlist Service: http://localhost:8084
- Order Service: http://localhost:8085
- Payment Service: http://localhost:8086
- Inventory Service: http://localhost:8087
- Notification Service: http://localhost:8088
- Reporting Service: http://localhost:8089

### Additional Tools
- Prometheus: http://localhost:9090 (when monitoring stack running)
- Grafana: http://localhost:3001 (when monitoring stack running)
- Alertmanager: http://localhost:9093 (when monitoring stack running
./scripts/health-check.sh

# View logs
make dev-logs
```

### Production Deployment

```bash
# Setup environment
cp .env.prod.example .env.prod
# Edit .env.prod with your production values

# Deploy production
make prod-up-build

# Run health check
./scripts/health-check.sh
```

See [PHASE8_SETUP.md](PHASE8_SETUP.md) for detailed Docker and deployment instructions.
     # React frontend application
├── auth-service/                  # Authentication service
├── product-service/               # Product catalog service
├── cart-service/                  # Shopping cart service
├── wishlist-service/              # Wishlist service
├── order-service/                 # Order management service
├── payment-service/               # Payment processing service
├── inventory-service/             # Inventory management service
├── notification-service/          # Notification service
├── reporting-service/             # Analytics and reporting service
├── nginx/                         # Nginx reverse proxy
├── monitoring/                    # Monitoring stack configs
│   ├── prometheus.yml            # Prometheus configuration
│   ├── alerts.yml                # Alert rules
│   ├── alertmanager.yml          # Alertmanager configuration
│   └── grafana/                  # Grafana datasources & dashboards
├── scripts/                       # Utility scripts
│   ├── health-check.sh           # Health check script
│   ├── backup.sh                 # Database backup script
│   ├── restore.sh                # Database restore script
│   └── setup-secrets.sh          # Docker secrets setup
├── docker-compose.yml             # Development orchestration
├── docker-compose.prod.yml        # Production orchestration
├── docker-compose.override.yml    # Development overrides
├── docker-compose.monitoring.yml  # Monitoring stack
├── docker-compose.secrets.yml     # Secrets configuration
├── Makefile                       # Docker management commands
├── .env.prod.example              # Production env template
├── .env.staging.example           # Staging env template
├── .dockerignore                  # Docker ignore rules
├── PHASE1_ARCHITECTURE.md         # Architecture documentation
├── PHASE8_SETUP.md                # Docker setup documentation
└── README.md                      # This file
```

## Features

### Phase 8 Highlights
- ✅ **Multi-Environment Support**: Dev, staging, and production configurations
- ✅ **Production Docker Compose**: Resource limits, health checks, networking
- ✅ **Nginx Reverse Proxy**: Rate limiting, security headers, gzip compression
- ✅ **Monitoring Stack**: Prometheus, Grafana, Alertmanager
- ✅ **Automated Backup/Restore**: Database backup and restore scripts
- ✅ **Health Monitoring**: Comprehensive health check system
- ✅ **Docker Secrets**: Secure secrets management
- ✅ **Makefile Commands**: 50+ commands for Docker operations
- ✅ **Network Isolation**: Separate backend and frontend networks
- ✅ **Structured Logging**: JSON logs with rotation

## Next Phases

### Phase 9: Event-Driven Architecture with Kafka
- Apache Kafka setup
- Event producers and consumers
- Event sourcing patterns
- Dead letter queues

### Phase 10: Kubernetes Deployment
- Kubernetes manifests
- Helm charts
- Ingress configuration
- Auto-scaling policies

### Phase 11: Advanced Monitoring & Observability
- Distributed tracing with Jaeger
- Log aggregation with ELK Stack
- APM integration
- Custom dashboards

### Phase 12: CI/CD Pipeline
- GitHub Actions workflows
- Automated testing
- Container registry integration
- Blue-green deployments

### Phase 13: API Gateway & Service Mesh
- Kong or Traefik API Gateway
- Istio service mesh
- Circuit breakers
- Retry policies payment-service/         # Payment processing service
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
