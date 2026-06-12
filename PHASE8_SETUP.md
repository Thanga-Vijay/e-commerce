# E-Commerce Platform - Phase 8 Implementation

This document provides comprehensive instructions for Docker containerization and orchestration (Phase 8).

## Overview

Phase 8 implements production-ready Docker containerization with:
- Multi-environment support (dev, staging, production)
- Advanced networking and security
- Resource limits and health checks
- Automated backup and restore
- Monitoring and logging
- Nginx reverse proxy with rate limiting

## Architecture

### Network Topology

```
┌─────────────────────────────────────┐
│         Nginx Reverse Proxy         │
│    (Port 80, Rate Limiting)         │
└──────────────┬──────────────────────┘
               │
        ┌──────┴──────┐
        │             │
   ┌────▼─────┐  ┌───▼──────┐
   │ Frontend │  │ Backend  │
   │ Network  │  │ Network  │
   └────┬─────┘  └────┬─────┘
        │             │
   ┌────▼─────┐  ┌───▼──────────────┐
   │ Frontend │  │ Microservices    │
   │  (3000)  │  │ (8081-8089)      │
   └──────────┘  └────┬─────────────┘
                      │
                 ┌────▼──────┐
                 │ Databases │
                 │ & Redis   │
                 └───────────┘
```

### Container Resources

**Databases:**
- CPU Limit: 0.5
- Memory Limit: 512MB
- Memory Reservation: 256MB

**Services:**
- CPU Limit: 0.5
- Memory Limit: 256MB
- Memory Reservation: 128MB

**Frontend:**
- CPU Limit: 0.25
- Memory Limit: 128MB

## Quick Start

### Development Environment

```bash
# Using Make (recommended)
make dev-up-build

# Or using docker-compose
docker-compose up -d --build

# View logs
make dev-logs

# Check health
./scripts/health-check.sh
```

### Production Environment

```bash
# Copy environment file
cp .env.prod.example .env.prod

# Edit with your production values
nano .env.prod

# Start production environment
make prod-up-build

# Or using docker-compose
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

### Staging Environment

```bash
# Copy environment file
cp .env.staging.example .env.staging

# Edit with your staging values
nano .env.staging

# Start staging environment
make staging-up
```

## Makefile Commands

### Development Commands
```bash
make dev-build          # Build all services
make dev-up             # Start development environment
make dev-up-build       # Build and start
make dev-down           # Stop development environment
make dev-logs           # View all logs
make dev-logs-service SERVICE=auth-service  # View specific service logs
make dev-restart        # Restart all services
make dev-ps             # List running containers
```

### Production Commands
```bash
make prod-build         # Build production images
make prod-up            # Start production environment
make prod-up-build      # Build and start production
make prod-down          # Stop production environment
make prod-logs          # View production logs
make prod-ps            # List production containers
```

### Database Commands
```bash
make db-migrate         # Run database migrations
make db-seed            # Seed databases
make db-backup          # Backup all databases
make db-restore FILE=backup.sql DB=auth_db  # Restore database
```

### Maintenance Commands
```bash
make clean              # Remove all containers and volumes
make clean-volumes      # Remove all volumes (data loss!)
make rebuild            # Rebuild everything from scratch
make restart-service SERVICE=auth-service  # Restart specific service
```

### Monitoring Commands
```bash
make stats              # Show resource usage
make health             # Check service health
make inspect-service SERVICE=auth-service  # Inspect service
```

### Image Commands
```bash
make pull               # Pull latest images
make push               # Push to registry
make tag VERSION=1.0.0  # Tag images with version
```

## Environment Configuration

### Development (.env)
- Uses default passwords
- Debug logging enabled
- Hot reload enabled
- Development ports exposed
- Includes dev tools (pgAdmin, Redis Commander)

### Staging (.env.staging)
- Staging passwords
- Test Stripe keys
- Staging SMTP (Mailtrap)
- Standard logging

### Production (.env.prod)
- Strong passwords (required)
- Live Stripe keys
- Production SMTP
- Error-only logging
- No dev tools

## Security Features

### Network Isolation
- **Backend Network**: Internal services only
- **Frontend Network**: Public-facing services
- Services communicate through defined networks
- Databases not exposed to frontend network

### Rate Limiting
- API endpoints: 10 requests/second
- Auth endpoints: 5 requests/second
- Burst allowance for traffic spikes
- Per-IP tracking

### Resource Limits
- CPU limits prevent resource exhaustion
- Memory limits prevent OOM issues
- Reservations ensure minimum resources
- Automatic restart on failure

### Health Checks
- Database: pg_isready checks
- Redis: PING command
- Services: HTTP health endpoints
- Start period allows service initialization

### Logging
- JSON format for parsing
- Maximum 10MB per file
- 3 file rotation
- Centralized log collection ready

## Backup and Restore

### Automated Backup

```bash
# Manual backup
./scripts/backup.sh

# Scheduled backup (add to cron)
0 2 * * * /path/to/scripts/backup.sh

# Backup location
./backups/ecommerce_backup_<timestamp>.tar.gz
```

### Restore from Backup

```bash
# List available backups
ls -lh ./backups/

# Restore specific backup
./scripts/restore.sh ./backups/ecommerce_backup_20260612_020000.tar.gz
```

### Backup Contents
- All 9 PostgreSQL databases (SQL dumps)
- Redis data (RDB file)
- Metadata file with timestamp
- Compressed tar.gz archive

## Health Monitoring

### Health Check Script

```bash
# Run health check
./scripts/health-check.sh

# Returns exit code 0 if all healthy
# Returns exit code 1 if any failures
```

### Health Check Endpoints

All services expose `/health` endpoints:
- **Auth**: http://localhost:8081/health
- **Product**: http://localhost:8082/health
- **Cart**: http://localhost:8083/health
- **Wishlist**: http://localhost:8084/health
- **Order**: http://localhost:8085/health
- **Payment**: http://localhost:8086/health
- **Inventory**: http://localhost:8087/health
- **Notification**: http://localhost:8088/health
- **Reporting**: http://localhost:8089/health

### Automated Monitoring

Integrate with monitoring systems:
```bash
# Prometheus healthcheck exporter
# Grafana dashboard
# Uptime monitoring (UptimeRobot, Pingdom)
```

## Nginx Reverse Proxy

### Features
- Single entry point for all services
- Path-based routing
- Rate limiting per endpoint
- Security headers
- Gzip compression
- Health check endpoint

### Configuration

```nginx
# Rate limiting zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=5r/s;

# Security headers
add_header X-Frame-Options "SAMEORIGIN";
add_header X-Content-Type-Options "nosniff";
add_header X-XSS-Protection "1; mode=block";
```

### Usage

```bash
# Build nginx
docker build -t ecommerce-nginx ./nginx/

# Run nginx
docker run -p 80:80 --network ecommerce-network ecommerce-nginx
```

## Production Deployment

### Pre-deployment Checklist

- [ ] Update `.env.prod` with strong passwords
- [ ] Configure production Stripe keys
- [ ] Set up production SMTP
- [ ] Generate secure JWT secret (32+ bytes)
- [ ] Configure domain in nginx
- [ ] Set up SSL certificates (Let's Encrypt)
- [ ] Configure firewall rules
- [ ] Set up backup schedule
- [ ] Configure monitoring
- [ ] Test health checks

### Deployment Steps

```bash
# 1. Prepare environment
cp .env.prod.example .env.prod
# Edit .env.prod with production values

# 2. Build images
make prod-build

# 3. Tag images with version
make tag VERSION=1.0.0

# 4. Push to registry (optional)
make push

# 5. Deploy
make prod-up

# 6. Verify health
./scripts/health-check.sh

# 7. Run migrations
make db-migrate

# 8. Set up backup cron job
crontab -e
# Add: 0 2 * * * /path/to/scripts/backup.sh
```

### Rolling Updates

```bash
# Update specific service
docker-compose -f docker-compose.prod.yml up -d --no-deps --build auth-service

# Scale service (requires orchestration)
docker-compose -f docker-compose.prod.yml up -d --scale product-service=3
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
make dev-logs-service SERVICE=auth-service

# Inspect container
docker inspect auth-service

# Check health
docker ps --filter "health=unhealthy"
```

### Database Connection Issues

```bash
# Test database connectivity
docker-compose exec auth-db pg_isready -U postgres

# Check database logs
docker-compose logs auth-db

# Connect to database
docker-compose exec auth-db psql -U postgres auth_db
```

### Memory Issues

```bash
# Check resource usage
make stats

# Adjust memory limits in docker-compose.prod.yml
deploy:
  resources:
    limits:
      memory: 512M  # Increase if needed
```

### Network Issues

```bash
# Inspect networks
make network-inspect

# Check network connectivity
docker-compose exec auth-service ping product-service

# Restart networking
docker-compose down
docker-compose up -d
```

### Health Check Failing

```bash
# Check service health endpoint
curl http://localhost:8081/health

# Increase health check timeout
healthcheck:
  interval: 30s
  timeout: 10s  # Increase if needed
  retries: 5
  start_period: 60s  # Increase for slow starts
```

## Best Practices

### Security
1. **Never commit secrets** - Use .env files (git-ignored)
2. **Use strong passwords** - Generate with `openssl rand -hex 32`
3. **Limit exposed ports** - Only expose what's needed
4. **Regular updates** - Keep base images updated
5. **Scan images** - Use `docker scan` or Trivy
6. **Use secrets management** - Docker secrets or Vault

### Performance
1. **Resource limits** - Prevent resource exhaustion
2. **Health checks** - Enable automatic recovery
3. **Logging rotation** - Prevent disk fill
4. **Image optimization** - Multi-stage builds
5. **Caching** - Use build cache effectively
6. **Volume optimization** - Use delegated mounts for dev

### Maintenance
1. **Regular backups** - Automated daily backups
2. **Monitor logs** - Centralized logging
3. **Health monitoring** - Automated health checks
4. **Update strategy** - Rolling updates
5. **Disaster recovery** - Test restore procedures
6. **Documentation** - Keep runbooks updated

## Monitoring Integration

### Prometheus

```yaml
# Add to docker-compose
prometheus:
  image: prom/prometheus
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
  ports:
    - "9090:9090"
```

### Grafana

```yaml
grafana:
  image: grafana/grafana
  ports:
    - "3001:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
```

### ELK Stack

```yaml
elasticsearch:
  image: elasticsearch:8.5.0
  environment:
    - discovery.type=single-node

logstash:
  image: logstash:8.5.0
  volumes:
    - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf

kibana:
  image: kibana:8.5.0
  ports:
    - "5601:5601"
```

## Next Steps

### Phase 9: Kafka Event Streaming
- Event-driven architecture
- Kafka clusters
- Event producers and consumers
- Event sourcing patterns

### Phase 10: Kubernetes Deployment
- Kubernetes manifests
- Helm charts
- Auto-scaling
- Load balancing

### Phase 11: Monitoring & Logging
- Prometheus + Grafana
- ELK Stack
- Distributed tracing
- APM integration

### Phase 12: CI/CD Pipeline
- GitHub Actions
- Automated testing
- Automated deployment
- GitOps workflow

### Phase 13: API Gateway & Load Balancing
- Kong or Traefik
- Load balancing
- API versioning
- Circuit breakers

## Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Docker Security](https://docs.docker.com/engine/security/)
- [Container Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Health Checks](https://docs.docker.com/engine/reference/builder/#healthcheck)
