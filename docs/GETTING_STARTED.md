# Getting Started

Quick start guide for setting up the E-Commerce Platform development environment.

## Prerequisites

### Required Software

1. **Go 1.21+**
   ```bash
   # Windows (using Chocolatey)
   choco install golang
   
   # Verify installation
   go version
   ```

2. **Node.js 20+ and npm**
   ```bash
   # Windows (using Chocolatey)
   choco install nodejs
   
   # Verify installation
   node --version
   npm --version
   ```

3. **Docker Desktop**
   - Download from https://www.docker.com/products/docker-desktop
   - Enable Kubernetes in Docker Desktop settings

4. **KIND (Kubernetes IN Docker)**
   ```bash
   # Windows (PowerShell)
   choco install kind
   
   # Verify installation
   kind --version
   ```

5. **kubectl**
   ```bash
   # Windows
   choco install kubernetes-cli
   
   # Verify installation
   kubectl version --client
   ```

6. **Git**
   ```bash
   # Windows
   choco install git
   
   # Verify installation
   git --version
   ```

### Optional Tools

- **PostgreSQL Client:** For direct database access
- **Redis CLI:** For Redis debugging
- **Postman/Insomnia:** For API testing
- **VS Code:** Recommended IDE with Go and React extensions

---

## Repository Setup

### 1. Clone Repository

```bash
cd C:/Users/ao45j/OneDrive\ -\ Cummins/Documents/K8S
git clone <repository-url> e-commerce
cd e-commerce
```

### 2. Environment Variables

Create `.env` files for each service. Example for Auth Service:

```bash
# auth-service/.env
PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_db
JWT_SECRET=your-256-bit-secret-key-here
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**Copy `.env.example` to `.env` for each service and fill in values.**

---

## Local Development Setup

### Option 1: Docker Compose (Recommended)

**Start all infrastructure services:**

```bash
cd infrastructure/docker
docker-compose up -d postgres redis kafka zookeeper
```

**Start a specific microservice:**

```bash
cd auth-service
go mod download
go run cmd/main.go
```

**Start frontend:**

```bash
cd frontend
npm install
npm run dev
```

### Option 2: Individual Services

**1. PostgreSQL**

```bash
# Using Docker
docker run --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  -d postgres:15
```

**2. Redis**

```bash
# Using Docker
docker run --name redis \
  -p 6379:6379 \
  -d redis:7-alpine
```

**3. Kafka**

```bash
# Using Docker Compose
cd infrastructure/kafka
docker-compose up -d
```

---

## Running Services

### Backend Services

Each service follows the same pattern:

```bash
# Example: Auth Service
cd auth-service

# Install dependencies
go mod download

# Run database migrations
go run migrations/migrate.go up

# Run service
go run cmd/main.go

# Service will start on its designated port (8081 for auth)
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Open browser to http://localhost:5173
```

### API Gateway

```bash
cd api-gateway

go mod download
go run cmd/main.go

# Gateway will start on port 8080
```

---

## Database Migrations

### Create Migration

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration
migrate create -ext sql -dir migrations -seq create_users_table

# This creates:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql
```

### Run Migrations

```bash
# Up (apply migrations)
migrate -path ./migrations \
  -database "postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable" \
  up

# Down (rollback one migration)
migrate -path ./migrations \
  -database "postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable" \
  down 1

# Check version
migrate -path ./migrations \
  -database "postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable" \
  version
```

---

## Testing

### Unit Tests

```bash
# Run tests for a service
cd auth-service
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Ensure test database is running
cd auth-service
go test -v -tags=integration ./tests/integration/...
```

### API Tests

Import Postman collection from `docs/api/postman-collection.json`

---

## Docker Commands

### Build Service Image

```bash
# Example: Auth Service
cd auth-service
docker build -t auth-service:latest .
```

### Run with Docker

```bash
docker run -d \
  --name auth-service \
  -p 8081:8081 \
  --env-file .env \
  auth-service:latest
```

### Docker Compose

```bash
# Start all services
cd infrastructure/docker
docker-compose up -d

# View logs
docker-compose logs -f auth-service

# Stop all services
docker-compose down

# Rebuild and start
docker-compose up -d --build
```

---

## KIND Kubernetes Setup

### Create Cluster

```bash
cd infrastructure/kind

# Create cluster
kind create cluster --config kind-config.yaml --name ecommerce

# Verify cluster
kubectl cluster-info --context kind-ecommerce
kubectl get nodes
```

### Load Docker Images

```bash
# Build images
cd auth-service
docker build -t auth-service:latest .

# Load into KIND
kind load docker-image auth-service:latest --name ecommerce
```

### Deploy Services

```bash
cd infrastructure/kubernetes

# Create namespace
kubectl apply -f namespace.yaml

# Deploy infrastructure
kubectl apply -f deployments/postgres.yaml
kubectl apply -f deployments/redis.yaml
kubectl apply -f deployments/kafka.yaml

# Deploy services
kubectl apply -f deployments/auth-service.yaml
kubectl apply -f deployments/product-service.yaml
# ... deploy all services

# Deploy ingress
kubectl apply -f ingress/ingress.yaml
```

### Verify Deployment

```bash
# Check all pods
kubectl get pods -n ecommerce

# Check services
kubectl get svc -n ecommerce

# Check logs
kubectl logs -f deployment/auth-service -n ecommerce

# Port forward to service
kubectl port-forward svc/auth-service 8081:8081 -n ecommerce
```

### Delete Cluster

```bash
kind delete cluster --name ecommerce
```

---

## Monitoring Setup

### Prometheus & Grafana

```bash
cd infrastructure/monitoring

# Deploy monitoring stack
./setup-monitoring.sh

# Access Grafana
kubectl port-forward svc/grafana 3000:3000 -n monitoring

# Open http://localhost:3000
# Default credentials: admin / admin
```

### View Metrics

```bash
# Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n monitoring

# Open http://localhost:9090
```

---

## Common Issues & Solutions

### Port Already in Use

```bash
# Windows - Find process using port
netstat -ano | findstr :8081

# Kill process (replace PID)
taskkill /PID <PID> /F
```

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -h localhost -U postgres -d auth_db

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

### Kafka Not Starting

```bash
# Check Zookeeper is running first
docker-compose up -d zookeeper

# Wait 30 seconds, then start Kafka
docker-compose up -d kafka

# Check logs
docker-compose logs -f kafka
```

### Go Module Issues

```bash
# Clear module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy modules
go mod tidy
```

### Docker Build Fails

```bash
# Clear Docker cache
docker builder prune -a

# Rebuild without cache
docker build --no-cache -t service:latest .
```

---

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/user-authentication
```

### 2. Make Changes

Edit code, add tests, update documentation.

### 3. Run Tests

```bash
go test -v ./...
```

### 4. Commit Changes

```bash
git add .
git commit -m "feat: implement user authentication"
```

### 5. Push and Create PR

```bash
git push origin feature/user-authentication

# Create Pull Request on GitHub
```

### 6. Code Review

Wait for review, address comments, merge to main.

---

## Next Steps

1. **Read Documentation:**
   - [High-Level Architecture](docs/architecture/HIGH_LEVEL_ARCHITECTURE.md)
   - [Service Architecture](docs/architecture/SERVICE_ARCHITECTURE.md)
   - [API Contracts](docs/api/API_CONTRACTS.md)
   - [Development Roadmap](docs/ROADMAP.md)

2. **Follow the Roadmap:**
   - Start with Phase 1: Architecture & Setup
   - Move to Phase 2: Auth & Product Services
   - Continue through all phases

3. **Join the Team:**
   - Set up communication channels (Slack, Teams)
   - Attend daily standups
   - Participate in code reviews

4. **Learn the Stack:**
   - Go programming
   - React and TypeScript
   - Docker and Kubernetes
   - PostgreSQL, Redis, Kafka

---

## Useful Commands Reference

### Go

```bash
go mod init          # Initialize module
go mod download      # Download dependencies
go mod tidy          # Clean up dependencies
go run cmd/main.go   # Run application
go test ./...        # Run all tests
go build             # Build binary
go fmt ./...         # Format code
go vet ./...         # Run Go vet
```

### Docker

```bash
docker ps                    # List running containers
docker ps -a                 # List all containers
docker logs <container>      # View logs
docker exec -it <container> bash  # Enter container
docker-compose up -d         # Start services
docker-compose down          # Stop services
docker-compose logs -f       # Follow logs
docker system prune          # Clean up
```

### Kubernetes

```bash
kubectl get pods -n ecommerce        # List pods
kubectl get svc -n ecommerce         # List services
kubectl logs <pod> -n ecommerce      # View logs
kubectl describe pod <pod> -n ecommerce  # Pod details
kubectl exec -it <pod> -n ecommerce -- bash  # Enter pod
kubectl port-forward <pod> 8080:8080 -n ecommerce  # Port forward
kubectl delete pod <pod> -n ecommerce  # Delete pod
kubectl apply -f <file>              # Apply manifest
```

### Git

```bash
git status              # Check status
git add .               # Stage all changes
git commit -m "message" # Commit changes
git push                # Push to remote
git pull                # Pull from remote
git checkout -b <branch>  # Create branch
git merge <branch>      # Merge branch
git log --oneline       # View commit history
```

---

## Support

For questions or issues:
- Check documentation in `/docs`
- Review existing GitHub issues
- Create new issue with detailed description
- Contact team lead

---

**Happy coding! 🚀**
