#!/bin/bash

set -e

echo "🚀 Setting up k3d cluster for E-Commerce Platform"

# Check if k3d is installed
if ! command -v k3d &> /dev/null; then
    echo "❌ k3d is not installed. Installing..."
    curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Delete existing cluster if it exists
if k3d cluster list | grep -q ecommerce-cluster; then
    echo "🗑️  Deleting existing cluster..."
    k3d cluster delete ecommerce-cluster
fi

# Create data directory for persistence
mkdir -p ./data

# Create k3d cluster
echo "📦 Creating k3d cluster..."
k3d cluster create --config k3d-config.yaml

# Wait for cluster to be ready
echo "⏳ Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=300s

# Create namespaces
echo "📁 Creating namespaces..."
kubectl create namespace ecommerce --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace floci --dry-run=client -o yaml | kubectl apply -f -

# Label namespaces
kubectl label namespace ecommerce env=local --overwrite
kubectl label namespace monitoring env=local --overwrite
kubectl label namespace floci env=local --overwrite

# Install NGINX Ingress Controller
echo "🌐 Installing NGINX Ingress Controller..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.1/deploy/static/provider/cloud/deploy.yaml

# Wait for ingress controller
echo "⏳ Waiting for ingress controller..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=300s

# Configure hosts file
echo "📝 Configuring /etc/hosts (requires sudo)..."
if ! grep -q "ecommerce.local" /etc/hosts 2>/dev/null; then
    echo "127.0.0.1 ecommerce.local api.ecommerce.local" | sudo tee -a /etc/hosts
fi

echo "✅ k3d cluster setup complete!"
echo ""
echo "📌 Cluster Info:"
echo "   Name: ecommerce-cluster"
echo "   Nodes: 1 server + 2 agents"
echo "   Registry: registry.localhost:5001"
echo ""
echo "🌐 Access Points:"
echo "   Frontend: http://ecommerce.local"
echo "   API: http://api.ecommerce.local"
echo ""
echo "🔧 Next Steps:"
echo "   1. Deploy Floci: ./floci-setup.sh"
echo "   2. Deploy databases: kubectl apply -f k8s/databases/"
echo "   3. Deploy Kafka: kubectl apply -f k8s/kafka/"
echo "   4. Deploy services: kubectl apply -f k8s/services/"
echo "   5. Check status: kubectl get pods -n ecommerce"
