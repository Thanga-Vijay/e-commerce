#!/bin/bash

# Build and Load Docker Images into KIND Cluster
# This script builds all microservice images and loads them into KIND

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;36m'
NC='\033[0m'

CLUSTER_NAME="ecommerce-local"
REGISTRY="localhost:5000"
VERSION="${VERSION:-latest}"

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Building and Loading Images to KIND                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if KIND cluster exists
if ! kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo -e "${RED}✗ KIND cluster '$CLUSTER_NAME' not found${NC}"
    echo "  Run ./kind-setup.sh first"
    exit 1
fi

# Services to build
SERVICES=(
    "auth-service"
    "product-service"
    "cart-service"
    "wishlist-service"
    "order-service"
    "payment-service"
    "inventory-service"
    "notification-service"
    "reporting-service"
)

# Build and load each service
for service in "${SERVICES[@]}"; do
    echo ""
    echo -e "${BLUE}Building ${service}...${NC}"
    
    # Build the image
    docker build -t ecommerce/${service}:${VERSION} \
        -f ${service}/Dockerfile \
        ${service}
    
    # Tag for local registry
    docker tag ecommerce/${service}:${VERSION} ${REGISTRY}/ecommerce/${service}:${VERSION}
    
    # Load into KIND cluster
    echo -e "${YELLOW}Loading ${service} into KIND...${NC}"
    kind load docker-image ecommerce/${service}:${VERSION} --name ${CLUSTER_NAME}
    
    # Push to local registry (optional)
    # docker push ${REGISTRY}/ecommerce/${service}:${VERSION}
    
    echo -e "${GREEN}✓ ${service} ready${NC}"
done

# Build frontend
echo ""
echo -e "${BLUE}Building frontend...${NC}"
docker build -t ecommerce/frontend:${VERSION} \
    -f frontend/Dockerfile \
    frontend

docker tag ecommerce/frontend:${VERSION} ${REGISTRY}/ecommerce/frontend:${VERSION}

echo -e "${YELLOW}Loading frontend into KIND...${NC}"
kind load docker-image ecommerce/frontend:${VERSION} --name ${CLUSTER_NAME}

echo -e "${GREEN}✓ frontend ready${NC}"

# Verify images in cluster
echo ""
echo -e "${BLUE}Verifying images in KIND cluster...${NC}"
docker exec -it ${CLUSTER_NAME}-control-plane crictl images | grep ecommerce || echo "No images found"

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  All Images Built and Loaded Successfully!            ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Deploy to cluster: cd k8s && ./deploy.sh"
echo "  2. Check pods: kubectl get pods -n e-commerce"
echo ""
echo -e "${BLUE}Available images:${NC}"
for service in "${SERVICES[@]}"; do
    echo "  - ecommerce/${service}:${VERSION}"
done
echo "  - ecommerce/frontend:${VERSION}"
echo ""
