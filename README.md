# 🛒 E-Commerce Platform - Cloud-Native Microservices

A production-ready, cloud-native e-commerce platform built with **Go microservices**, **React frontend**, deployed on **Kubernetes (k3d)** with **Floci** for cloud services emulation.

[![CI/CD](https://github.com/yourusername/e-commerce/workflows/CI/badge.svg)](https://github.com/yourusername/e-commerce/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [**🍎 MacBook Setup**](#-macbook-setup-new) ⭐
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Monitoring](#-monitoring)
- [API Documentation](#-api-documentation)
- [Troubleshooting](#-troubleshooting)
- [License](#-license)

---

## 🍎 MacBook Setup (NEW)

**Running locally on your MacBook?** We've optimized the setup for you!

### One-Line Setup
```bash
./setup-all.sh
```

### Manual Setup
```bash
# 1. Create k3d cluster (1 server + 2 agents)
./k3d-setup.sh

# 2. Deploy Floci (cloud services emulator)
./floci-setup.sh

# 3. Deploy everything
kubectl apply -f k8s/databases/
kubectl apply -f k8s/services/
kubectl apply -f k8s/frontend/
```

### 📚 Complete Guides
- **[MACBOOK_SETUP.md](MACBOOK_SETUP.md)** - Detailed setup guide with troubleshooting
- **[K3D_CLUSTER_GUIDE.md](K3D_CLUSTER_GUIDE.md)** - Why 1 server + 2 agents? Resource allocation explained
- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Essential commands cheat sheet
- **[DOCKER_COMPOSE_CLEANUP.md](DOCKER_COMPOSE_CLEANUP.md)** - Docker Compose vs k3d guide

### System Requirements
- **RAM:** 16 GB minimum (32 GB recommended)
- **CPU:** 4+ cores
- **Storage:** 50 GB free
- **macOS:** 12.0+ (Monterey or later)

### k3d Configuration
- **Cluster:** 1 server + 2 agents (optimal for MacBook)
- **Registry:** Local registry at registry.localhost:5001
- **Resources:** ~10 GB RAM, ~5-6 CPU cores
- **Services:** 9 microservices + databases + Kafka + monitoring

### Access After Setup
- Frontend: http://ecommerce.local
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001
- Floci: http://localhost:30456

---

## ✨ Features

### Core Functionality
- ✅ **User Authentication** - JWT-based auth with role-based access control (RBAC)
- ✅ **Product Catalog** - Browse, search, filter products by category
- ✅ **Shopping Cart** - Add, update, remove items with real-time sync
- ✅ **Wishlist** - Save products for later
- ✅ **Order Management** - Place orders, track status, view history
- ✅ **Payment Processing** - Secure payment gateway integration
- ✅ **Inventory Management** - Real-time stock tracking
- ✅ **Notifications** - Email/SMS notifications via Kafka events
- ✅ **Reporting & Analytics** - Sales reports, user analytics

### Platform Features
- 🚀 **Microservices Architecture** - 9 independent Go services
- 🔄 **Event-Driven** - Kafka for asynchronous communication
- 🐳 **Containerized** - Docker images for all services
- ☸️ **Kubernetes Native** - Deployed on k3d (lightweight K8s)
- ☁️ **AWS Services** - S3, ECR, Secrets Manager via LocalStack
- 📊 **Observability** - Prometheus metrics, Grafana dashboards
- 🔐 **Security** - mTLS, secrets management, security scanning
- 🔄 **CI/CD** - GitHub Actions with automated testing
- 🔧 **GitOps** - ArgoCD for declarative deployments

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                        │
│                    http://ecommerce.local                       │
└────────────────────────────┬────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  NGINX Ingress  │
                    └────────┬────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
┌───────▼───────┐   ┌───────▼───────┐   ┌───────▼───────┐
│ Auth Service  │   │Product Service│   │ Cart Service  │
│   :8081       │   │   :8082       │   │   :8083       │
└───────┬───────┘   └───────┬───────┘   └───────┬───────┘
        │                   │                    │
        └───────────────────┼────────────────────┘
                            │
                ┌───────────┴───────────┐
                │                       │
        ┌───────▼───────┐       ┌──────▼──────┐
        │  PostgreSQL   │       │    Kafka    │
        │  (9 databases)│       │  (Events)   │
        └───────────────┘       └─────────────┘
                │
        ┌───────▼───────┐
        │  LocalStack   │
        │  AWS Services │
        │  S3, ECR, SM  │
        └───────────────┘
```

### Microservices

| Service | Port | Description | Database |
|---------|------|-------------|----------|
| **Auth** | 8081 | User authentication, JWT tokens | auth_db |
| **Product** | 8082 | Product catalog, categories | product_db |
| **Cart** | 8083 | Shopping cart management | cart_db |
| **Wishlist** | 8084 | User wishlist | wishlist_db |
| **Order** | 8085 | Order processing, history | order_db |
| **Payment** | 8086 | Payment gateway integration | payment_db |
| **Inventory** | 8087 | Stock management | inventory_db |
| **Notification** | 8088 | Email/SMS notifications | notification_db |
| **Reporting** | 8089 | Analytics, reports | reporting_db |
| **Frontend** | 3000 | React SPA | - |

---

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.22
- **Framework:** Gin (HTTP routing)
- **Database:** PostgreSQL 15 (9 separate databases)
- **Cache:** Redis 7
- **Message Queue:** Apache Kafka 3.6
- **Authentication:** JWT (golang-jwt/jwt/v5)

### Frontend
- **Framework:** React 18
- **Build Tool:** Vite 5
- **Styling:** Tailwind CSS 3
- **HTTP Client:** Axios
- **Routing:** React Router 6
- **State Management:** Context API

### Infrastructure
- **Container Runtime:** Docker 24
- **Orchestration:** Kubernetes (k3d v1.30.8)
- **Ingress:** NGINX Ingress Controller
- **AWS Emulation:** LocalStack 3.x
- **CI/CD:** GitHub Actions
- **GitOps:** ArgoCD (optional)

### Monitoring & Observability
- **Metrics:** Prometheus
- **Visualization:** Grafana
- **Logging:** ELK Stack (optional)
- **Tracing:** Jaeger (optional)

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

### Required
- **Docker Desktop** 24.0+ ([Install](https://docs.docker.com/get-docker/))
- **kubectl** 1.28+ ([Install](https://kubernetes.io/docs/tasks/tools/))
- **k3d** 5.0+ ([Install](https://k3d.io/))
- **Git** 2.40+

### Optional (for development)
- **Go** 1.22+ ([Install](https://go.dev/dl/))
- **Node.js** 18+ ([Install](https://nodejs.org/))
- **AWS CLI** 2.0+ ([Install](https://aws.amazon.com/cli/))
- **Make** (for Makefile commands)

### System Requirements
- **OS:** Linux, macOS, or Windows with WSL2
- **RAM:** 8GB minimum, 16GB recommended
- **CPU:** 4 cores minimum
- **Disk:** 20GB free space

---

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/e-commerce.git
cd e-commerce
```

### 2. Set Up k3d Cluster

```bash
# Make scripts executable
chmod +x k3d-setup.sh localstack-setup.sh scripts/*.sh

# Create k3d cluster (takes ~3 minutes)
./k3d-setup.sh
```

**What this does:**
- Creates a k3d cluster with 1 server + 2 agents
- Installs NGINX Ingress Controller
- Creates namespaces: `ecommerce`, `monitoring`, `localstack`
- Sets up local registry at `registry.localhost:5001`
- Configures `/etc/hosts` entries

### 3. Deploy LocalStack (AWS Emulation)

```bash
# Deploy LocalStack to k3d
./localstack-setup.sh

# Initialize AWS resources (S3, ECR, Secrets)
./scripts/init-localstack.sh
```

**LocalStack provides:**
- S3 buckets for images, backups, logs
- ECR repositories for container images
- Secrets Manager for credentials
- CloudWatch for logging

### 4. Build and Push Images

```bash
# Build all service images
make build-all

# Tag and push to local registry
make push-all
```

### 5. Deploy to Kubernetes

```bash
# Deploy infrastructure (PostgreSQL, Redis, Kafka)
kubectl apply -f k8s/postgres/
kubectl apply -f k8s/redis/
kubectl apply -f k8s/kafka/

# Wait for infrastructure to be ready
kubectl wait --for=condition=Ready pod -l app=postgres -n ecommerce --timeout=300s
kubectl wait --for=condition=Ready pod -l app=kafka -n ecommerce --timeout=300s

# Deploy ConfigMap and Secrets
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# Deploy all microservices
kubectl apply -f k8s/deployments/

# Deploy Ingress
kubectl apply -f k8s/ingress/
```

### 6. Verify Deployment

```bash
# Check all pods are running
kubectl get pods -n ecommerce

# Check services
kubectl get svc -n ecommerce

# Check ingress
kubectl get ingress -n ecommerce
```

### 7. Access the Application

Open your browser and navigate to:

- **Frontend:** http://ecommerce.local
- **API:** http://api.ecommerce.local/health
- **LocalStack:** http://localhost:30456/_localstack/health

---

## 📁 Project Structure

```
e-commerce/
├── .github/workflows/       # CI/CD pipelines
├── k8s/                     # Kubernetes manifests
│   ├── localstack/          # LocalStack deployment
│   ├── deployments/         # Service deployments
│   ├── postgres/            # PostgreSQL StatefulSet
│   ├── redis/               # Redis deployment
│   ├── kafka/               # Kafka cluster
│   └── ingress/             # Ingress rules
├── auth-service/            # Authentication microservice
├── product-service/         # Product catalog microservice
├── cart-service/            # Shopping cart microservice
├── wishlist-service/        # Wishlist microservice
├── order-service/           # Order management microservice
├── payment-service/         # Payment processing microservice
├── inventory-service/       # Inventory management microservice
├── notification-service/    # Notification microservice
├── reporting-service/       # Analytics microservice
├── frontend/                # React frontend
├── scripts/                 # Utility scripts
├── k3d-config.yaml          # k3d cluster configuration
├── k3d-setup.sh             # k3d cluster setup script
├── localstack-setup.sh      # LocalStack setup script
├── Makefile                 # Make commands
└── README.md                # This file
```

---

## 💻 Development

### Local Development (without Kubernetes)

```bash
# Start infrastructure with Docker Compose
docker-compose up -d postgres redis kafka

# Run a service locally (example: auth-service)
cd auth-service
go mod download
go run main.go

# Run frontend locally
cd frontend
npm install
npm run dev
```

---

## 🧪 Testing

### Run Unit Tests

```bash
# Test all Go services
make test

# Test specific service
cd auth-service
go test -v ./...
go test -race ./...
```

---

## 🚢 Deployment

### Automated Deployment (GitHub Actions)

The CI/CD pipeline automatically:

1. **On PR:** Runs tests, linting, security scans
2. **On merge to `develop`:** Builds and tests all services
3. **On merge to `main`:** Builds, pushes to registry, deploys to cluster

---

## 📊 Monitoring

### Prometheus Metrics

Each service exposes metrics at `/metrics`:

```bash
curl http://api.ecommerce.local/auth/metrics
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. Pods not starting

```bash
# Check pod status
kubectl get pods -n ecommerce

# Describe pod to see events
kubectl describe pod <pod-name> -n ecommerce

# Check logs
kubectl logs <pod-name> -n ecommerce
```

#### 2. Cannot access application

```bash
# Check ingress
kubectl get ingress -n ecommerce

# Verify /etc/hosts
cat /etc/hosts | grep ecommerce
```

### Cleanup

```bash
# Delete all resources
kubectl delete namespace ecommerce monitoring localstack

# Delete k3d cluster
k3d cluster delete ecommerce-cluster

# Clean Docker resources
docker system prune -a --volumes
```

---

## 📄 License

This project is licensed under the MIT License.

---

**Made with ❤️ by the E-Commerce Team**
- ✅ **Health Monitoring**: Comprehensive health check system
- ✅ **Docker Secrets**: Secure secrets management
- ✅ **Makefile Commands**: 70+ commands for Docker & Kafka operations
- ✅ **Network Isolation**: Separate backend and frontend networks
- ✅ **Structured Logging**: JSON logs with rotation

## Next Phases

### Phase 10: Kubernetes Deployment
- Kubernetes manifests for all services
- StatefulSets for Kafka and databases
- Helm charts for easy deployment
- Ingress configuration with TLS
- Horizontal Pod Autoscaling (HPA)
- ConfigMaps and Secrets management

### Phase 11: Advanced Monitoring & Observability
- Distributed tracing with Jaeger
- Log aggregation with ELK Stack
- APM integration (Application Performance Monitoring)
- Custom Grafana dashboards
- Kafka metrics and lag monitoring
- SLI/SLO definition and tracking

### Phase 12: CI/CD Pipeline
- GitHub Actions workflows
- Automated testing (unit, integration, e2e)
- Container registry integration
- Blue-green deployments
- Canary releases
- Automated rollback mechanisms

### Phase 13: API Gateway & Service Mesh
- Kong or Traefik API Gateway
- Istio service mesh integration
- Circuit breakers and retry policies
- Rate limiting and throttling
- mTLS for service-to-service communication
- Advanced traffic management

## License

MIT
