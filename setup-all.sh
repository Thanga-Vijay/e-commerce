#!/bin/bash

# 🚀 E-Commerce Platform - MacBook Quick Setup
# Run this script to set up everything in one go

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║   E-Commerce Platform - MacBook k3d Setup                      ║"
echo "║   1 Server + 2 Agents with Floci                               ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check prerequisites
echo "📋 Checking prerequisites..."
echo ""

check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}❌ $1 is not installed${NC}"
        echo "   Install: brew install $2"
        exit 1
    else
        echo -e "${GREEN}✅ $1 is installed${NC}"
    fi
}

check_command "docker" "docker"
check_command "kubectl" "kubectl"
check_command "k3d" "k3d"

# Check Docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}❌ Docker is not running${NC}"
    echo "   Please start Docker Desktop first"
    exit 1
else
    echo -e "${GREEN}✅ Docker is running${NC}"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Ask user confirmation
echo -e "${YELLOW}This will:${NC}"
echo "  1. Create k3d cluster (1 server + 2 agents)"
echo "  2. Install NGINX Ingress"
echo "  3. Deploy Floci (cloud services emulator)"
echo "  4. Deploy PostgreSQL databases (9)"
echo "  5. Deploy Redis"
echo "  6. Deploy Kafka + Zookeeper"
echo "  7. Deploy 9 microservices"
echo "  8. Deploy frontend"
echo "  9. Set up monitoring (Prometheus + Grafana)"
echo ""
echo -e "${YELLOW}Estimated time: 10-15 minutes${NC}"
echo -e "${YELLOW}Required resources: ~10 GB RAM, 6 CPU cores${NC}"
echo ""

read -p "Continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 0
fi

echo ""
echo "🚀 Starting setup..."
echo ""

# Step 1: Create k3d cluster
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Step 1/9: Creating k3d cluster..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./k3d-setup.sh || { echo -e "${RED}Failed to create k3d cluster${NC}"; exit 1; }
echo -e "${GREEN}✅ k3d cluster created${NC}"
echo ""

# Step 2: Deploy Floci
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "☁️  Step 2/9: Deploying Floci..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./floci-setup.sh || { echo -e "${RED}Failed to deploy Floci${NC}"; exit 1; }
echo -e "${GREEN}✅ Floci deployed${NC}"
echo ""

# Step 3: Deploy databases
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🗄️  Step 3/9: Deploying PostgreSQL databases..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/databases/ || { echo -e "${RED}Failed to deploy databases${NC}"; exit 1; }
echo "⏳ Waiting for databases to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n ecommerce --timeout=300s || true
echo -e "${GREEN}✅ Databases deployed${NC}"
echo ""

# Step 4: Deploy Redis
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💾 Step 4/9: Deploying Redis..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/redis/ || { echo -e "${RED}Failed to deploy Redis${NC}"; exit 1; }
kubectl wait --for=condition=ready pod -l app=redis -n ecommerce --timeout=300s || true
echo -e "${GREEN}✅ Redis deployed${NC}"
echo ""

# Step 5: Deploy Kafka
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📨 Step 5/9: Deploying Kafka..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/kafka/ || { echo -e "${RED}Failed to deploy Kafka${NC}"; exit 1; }
echo "⏳ Waiting for Kafka to be ready..."
kubectl wait --for=condition=ready pod -l app=kafka -n ecommerce --timeout=300s || true
echo -e "${GREEN}✅ Kafka deployed${NC}"
echo ""

# Step 6: Deploy ConfigMaps and Secrets
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔐 Step 6/9: Creating ConfigMaps and Secrets..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/configmap.yaml || true
kubectl apply -f k8s/secrets.yaml || true
echo -e "${GREEN}✅ ConfigMaps and Secrets created${NC}"
echo ""

# Step 7: Deploy microservices
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚙️  Step 7/9: Deploying microservices..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/services/ || { echo -e "${RED}Failed to deploy services${NC}"; exit 1; }
echo "⏳ Waiting for services to start..."
sleep 10
echo -e "${GREEN}✅ Microservices deployed${NC}"
echo ""

# Step 8: Deploy frontend
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎨 Step 8/9: Deploying frontend..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/frontend/ || { echo -e "${RED}Failed to deploy frontend${NC}"; exit 1; }
kubectl apply -f k8s/ingress/ || { echo -e "${RED}Failed to deploy ingress${NC}"; exit 1; }
echo -e "${GREEN}✅ Frontend deployed${NC}"
echo ""

# Step 9: Deploy monitoring
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Step 9/9: Deploying monitoring stack..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
kubectl apply -f k8s/monitoring/ || { echo -e "${YELLOW}⚠️  Monitoring deployment had issues (non-critical)${NC}"; }
echo -e "${GREEN}✅ Monitoring deployed${NC}"
echo ""

# Final status
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 Setup Complete!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 Deployment Status:"
kubectl get pods -n ecommerce
echo ""
echo "🌐 Access Points:"
echo "   Frontend:    http://ecommerce.local"
echo "   Prometheus:  http://localhost:9090"
echo "   Grafana:     http://localhost:3001 (admin/admin)"
echo "   Floci:       http://localhost:30456"
echo ""
echo "🔍 Useful Commands:"
echo "   Check status:  kubectl get pods -n ecommerce"
echo "   View logs:     kubectl logs -n ecommerce -l app=<service-name>"
echo "   Stop cluster:  k3d cluster stop ecommerce-cluster"
echo "   Delete cluster: k3d cluster delete ecommerce-cluster"
echo ""
echo "📚 Full documentation: MACBOOK_SETUP.md"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
