#!/bin/bash

# Health Check Script for E-Commerce Platform
# Checks health of all services and reports status

set -e

BLUE='\033[36m'
GREEN='\033[32m'
RED='\033[31m'
YELLOW='\033[33m'
RESET='\033[0m'

echo -e "${BLUE}E-Commerce Platform Health Check${RESET}"
echo "=================================="
echo ""

FAILED=0

check_service() {
    local name=$1
    local url=$2
    local expected=$3

    printf "Checking %-20s ... " "$name"
    
    if response=$(curl -s -o /dev/null -w "%{http_code}" --max-time 5 "$url" 2>&1); then
        if [ "$response" = "$expected" ]; then
            echo -e "${GREEN}✓ OK${RESET} (HTTP $response)"
        else
            echo -e "${YELLOW}⚠ WARNING${RESET} (HTTP $response, expected $expected)"
            FAILED=$((FAILED + 1))
        fi
    else
        echo -e "${RED}✗ FAILED${RESET} (Connection failed)"
        FAILED=$((FAILED + 1))
    fi
}

check_database() {
    local name=$1
    local container=$2
    local db=$3

    printf "Checking %-20s ... " "$name"
    
    if docker-compose exec -T "$container" pg_isready -U postgres >/dev/null 2>&1; then
        echo -e "${GREEN}✓ OK${RESET}"
    else
        echo -e "${RED}✗ FAILED${RESET}"
        FAILED=$((FAILED + 1))
    fi
}

# Check Databases
echo -e "${BLUE}Databases:${RESET}"
check_database "Auth DB" "auth-db" "auth_db"
check_database "Product DB" "product-db" "product_db"
check_database "Cart DB" "cart-db" "cart_db"
check_database "Wishlist DB" "wishlist-db" "wishlist_db"
check_database "Order DB" "order-db" "order_db"
check_database "Payment DB" "payment-db" "payment_db"
check_database "Inventory DB" "inventory-db" "inventory_db"
check_database "Notification DB" "notification-db" "notification_db"
check_database "Reporting DB" "reporting-db" "reporting_db"
echo ""

# Check Redis
echo -e "${BLUE}Cache:${RESET}"
printf "Checking %-20s ... " "Redis"
if docker-compose exec -T redis redis-cli ping >/dev/null 2>&1; then
    echo -e "${GREEN}✓ OK${RESET}"
else
    echo -e "${RED}✗ FAILED${RESET}"
    FAILED=$((FAILED + 1))
fi
echo ""

# Check Services
echo -e "${BLUE}Microservices:${RESET}"
check_service "Auth Service" "http://localhost:8081/health" "200"
check_service "Product Service" "http://localhost:8082/health" "200"
check_service "Cart Service" "http://localhost:8083/health" "200"
check_service "Wishlist Service" "http://localhost:8084/health" "200"
check_service "Order Service" "http://localhost:8085/health" "200"
check_service "Payment Service" "http://localhost:8086/health" "200"
check_service "Inventory Service" "http://localhost:8087/health" "200"
check_service "Notification Service" "http://localhost:8088/health" "200"
check_service "Reporting Service" "http://localhost:8089/health" "200"
echo ""

# Check Frontend
echo -e "${BLUE}Frontend:${RESET}"
check_service "Frontend" "http://localhost:3000" "200"
echo ""

# Summary
echo "=================================="
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All services are healthy!${RESET}"
    exit 0
else
    echo -e "${RED}$FAILED service(s) failed health check${RESET}"
    exit 1
fi
