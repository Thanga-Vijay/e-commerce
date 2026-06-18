#!/bin/bash

# KIND Cluster Setup Script for E-Commerce Platform
# This script sets up a complete KIND cluster for local testing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  E-Commerce Platform - KIND Cluster Setup             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if ! command -v kind &> /dev/null; then
    echo -e "${RED}✗ KIND not found. Please install KIND first:${NC}"
    echo "  https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
fi
echo -e "${GREEN}✓ KIND installed${NC}"

if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}✗ kubectl not found. Please install kubectl first${NC}"
    exit 1
fi
echo -e "${GREEN}✓ kubectl installed${NC}"

if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker not found. Please install Docker Desktop${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker installed${NC}"

# Check if cluster already exists
if kind get clusters | grep -q "ecommerce-local"; then
    echo -e "${YELLOW}Cluster 'ecommerce-local' already exists.${NC}"
    read -p "Do you want to delete and recreate it? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Deleting existing cluster...${NC}"
        kind delete cluster --name ecommerce-local
    else
        echo -e "${GREEN}Using existing cluster${NC}"
        kubectl cluster-info --context kind-ecommerce-local
        exit 0
    fi
fi

# Create data directory for persistent volumes
echo -e "${YELLOW}Creating data directory...${NC}"
mkdir -p ./data
echo -e "${GREEN}✓ Data directory created${NC}"

# Create KIND cluster
echo ""
echo -e "${BLUE}Creating KIND cluster (this may take 2-3 minutes)...${NC}"
kind create cluster --config kind-config.yaml --wait 5m

# Verify cluster
echo ""
echo -e "${YELLOW}Verifying cluster...${NC}"
kubectl cluster-info --context kind-ecommerce-local
kubectl get nodes

# Install Ingress NGINX Controller
echo ""
echo -e "${BLUE}Installing Ingress NGINX Controller...${NC}"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

# Wait for Ingress controller to be ready
echo -e "${YELLOW}Waiting for Ingress controller to be ready...${NC}"
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

echo -e "${GREEN}✓ Ingress NGINX installed${NC}"

# Install Metrics Server (for HPA)
echo ""
echo -e "${BLUE}Installing Metrics Server...${NC}"
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Patch metrics-server to work with KIND
kubectl patch -n kube-system deployment metrics-server --type=json \
  -p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'

echo -e "${GREEN}✓ Metrics Server installed${NC}"

# Create local container registry (optional but recommended)
echo ""
echo -e "${BLUE}Creating local container registry...${NC}"

# Check if registry already exists
if [ "$(docker ps -q -f name=kind-registry)" ]; then
    echo -e "${GREEN}✓ Registry already running${NC}"
else
    docker run -d --restart=always \
      -p 5001:5000 \
      --name kind-registry \
      registry:2

    # Connect registry to KIND network
    docker network connect kind kind-registry || true
    
    echo -e "${GREEN}✓ Local registry created at localhost:5001${NC}"
fi

# Label nodes for different workloads
echo ""
echo -e "${BLUE}Labeling nodes...${NC}"
kubectl label nodes ecommerce-local-worker tier=application --overwrite
kubectl label nodes ecommerce-local-worker2 tier=application --overwrite
kubectl label nodes ecommerce-local-worker3 tier=database --overwrite

echo -e "${GREEN}✓ Nodes labeled${NC}"

# Create StorageClass for local-path provisioner
echo ""
echo -e "${BLUE}Configuring storage...${NC}"
kubectl apply -f - <<EOF
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: rancher.io/local-path
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Delete
EOF

echo -e "${GREEN}✓ Storage configured${NC}"

# Display cluster information
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  KIND Cluster Setup Complete!                         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Cluster Information:${NC}"
echo "  Name: ecommerce-local"
echo "  Nodes: 4 (1 control-plane, 3 workers)"
echo "  Context: kind-ecommerce-local"
echo "  Registry: localhost:5001"
echo ""
echo -e "${GREEN}Access Points:${NC}"
echo "  Ingress HTTP: http://localhost:80"
echo "  Ingress HTTPS: https://localhost:443"
echo "  Frontend: http://localhost:3000"
echo "  Grafana: http://localhost:3001"
echo "  Prometheus: http://localhost:9090"
echo "  Jaeger: http://localhost:16686"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "  1. Build and load Docker images:"
echo "     ./scripts/build-and-load-images.sh"
echo ""
echo "  2. Deploy the e-commerce platform:"
echo "     cd k8s && ./deploy.sh"
echo ""
echo "  3. Deploy monitoring stack:"
echo "     cd k8s/monitoring && ./deploy-monitoring.sh"
echo ""
echo "  4. Check cluster status:"
echo "     kubectl get nodes"
echo "     kubectl get pods -A"
echo ""
echo -e "${GREEN}To delete the cluster:${NC}"
echo "  kind delete cluster --name ecommerce-local"
echo ""
