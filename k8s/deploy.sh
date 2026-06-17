#!/bin/bash

# Deploy E-Commerce Platform to Kubernetes
# This script deploys all resources in the correct order

set -e

NAMESPACE="ecommerce"
BLUE='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

echo -e "${BLUE}E-Commerce Platform - Kubernetes Deployment${RESET}"
echo "=============================================="
echo ""

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo -e "${YELLOW}kubectl not found. Please install kubectl first.${RESET}"
    exit 1
fi

# Check if cluster is accessible
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${YELLOW}Cannot connect to Kubernetes cluster. Please check your connection.${RESET}"
    exit 1
fi

echo -e "${GREEN}✓ kubectl installed and cluster accessible${RESET}"
echo ""

# 1. Create namespace
echo -e "${BLUE}1. Creating namespace...${RESET}"
kubectl apply -f namespace.yaml
echo ""

# 2. Create ConfigMaps and Secrets
echo -e "${BLUE}2. Creating ConfigMaps and Secrets...${RESET}"
echo -e "${YELLOW}⚠ Remember to update secrets.yaml with your actual credentials!${RESET}"
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml
echo ""

# 3. Deploy databases
echo -e "${BLUE}3. Deploying PostgreSQL databases...${RESET}"
kubectl apply -f databases/
echo "Waiting for databases to be ready..."
kubectl wait --for=condition=ready pod -l app=auth-db -n $NAMESPACE --timeout=300s || true
echo ""

# 4. Deploy Redis
echo -e "${BLUE}4. Deploying Redis...${RESET}"
kubectl apply -f redis/
echo "Waiting for Redis to be ready..."
kubectl wait --for=condition=ready pod -l app=redis -n $NAMESPACE --timeout=120s || true
echo ""

# 5. Deploy Kafka
echo -e "${BLUE}5. Deploying Kafka and Zookeeper...${RESET}"
kubectl apply -f kafka/
echo "Waiting for Kafka to be ready..."
sleep 30
echo ""

# 6. Deploy microservices
echo -e "${BLUE}6. Deploying microservices...${RESET}"
kubectl apply -f services/
echo "Waiting for services to be ready..."
kubectl wait --for=condition=ready pod -l app=auth-service -n $NAMESPACE --timeout=180s || true
echo ""

# 7. Deploy frontend
echo -e "${BLUE}7. Deploying frontend...${RESET}"
kubectl apply -f frontend/
echo ""

# 8. Deploy Ingress
echo -e "${BLUE}8. Deploying Ingress...${RESET}"
kubectl apply -f ingress/
echo ""

# 9. Deploy HPA
echo -e "${BLUE}9. Deploying Horizontal Pod Autoscalers...${RESET}"
kubectl apply -f hpa/
echo ""

# Display status
echo -e "${GREEN}✓ Deployment complete!${RESET}"
echo ""
echo -e "${BLUE}Checking deployment status...${RESET}"
kubectl get all -n $NAMESPACE
echo ""

echo -e "${GREEN}Next steps:${RESET}"
echo "1. Update DNS records to point to your Ingress IP"
echo "2. Configure TLS certificates (cert-manager)"
echo "3. Run database migrations"
echo "4. Monitor pod status: kubectl get pods -n $NAMESPACE -w"
echo "5. Check logs: kubectl logs -f <pod-name> -n $NAMESPACE"
