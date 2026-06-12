# Docker Configurations

Docker and Docker Compose files for all services.

## Files

- `docker-compose.yml` - Main compose file for all services
- `docker-compose.dev.yml` - Development overrides
- `docker-compose.prod.yml` - Production overrides

## Services in Docker Compose

### Infrastructure Services
- PostgreSQL (multiple databases)
- Redis
- Kafka
- Zookeeper

### Application Services
- API Gateway
- Auth Service
- Product Service
- Cart Service
- Wishlist Service
- Order Service
- Payment Service
- Inventory Service
- Notification Service
- Reporting Service
- Frontend

### Monitoring Services
- Prometheus
- Grafana
- Loki

## Usage

Start all services:
```bash
docker-compose up -d
```

Start specific service:
```bash
docker-compose up -d auth-service
```

View logs:
```bash
docker-compose logs -f auth-service
```

Stop all services:
```bash
docker-compose down
```

Rebuild services:
```bash
docker-compose up -d --build
```
