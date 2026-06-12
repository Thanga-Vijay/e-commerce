#!/bin/bash

# Undeploy E-Commerce Platform from Kubernetes

set -e

NAMESPACE="ecommerce"
RED='\033[31m'
YELLOW='\033[33m'
RESET='\033[0m'

echo -e "${RED}E-Commerce Platform - Kubernetes Teardown${RESET}"
echo "=========================================="
echo ""

echo -e "${YELLOW}⚠ This will delete all resources in the $NAMESPACE namespace!${RESET}"
echo -e "${YELLOW}Press Ctrl+C to cancel, or Enter to continue...${RESET}"
read

echo "Deleting HPA..."
kubectl delete -f k8s/hpa/ --ignore-not-found=true

echo "Deleting Ingress..."
kubectl delete -f k8s/ingress/ --ignore-not-found=true

echo "Deleting frontend..."
kubectl delete -f k8s/frontend/ --ignore-not-found=true

echo "Deleting microservices..."
kubectl delete -f k8s/services/ --ignore-not-found=true

echo "Deleting Kafka..."
kubectl delete -f k8s/kafka/ --ignore-not-found=true

echo "Deleting Redis..."
kubectl delete -f k8s/redis/ --ignore-not-found=true

echo "Deleting databases..."
kubectl delete -f k8s/databases/ --ignore-not-found=true

echo "Deleting ConfigMaps and Secrets..."
kubectl delete -f k8s/configmap.yaml --ignore-not-found=true
kubectl delete -f k8s/secrets.yaml --ignore-not-found=true

echo "Deleting namespace..."
kubectl delete -f k8s/namespace.yaml --ignore-not-found=true

echo ""
echo "✓ Teardown complete!"
