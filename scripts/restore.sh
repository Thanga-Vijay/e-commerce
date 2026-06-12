#!/bin/bash

# Restore Script for E-Commerce Platform
# Restores databases from backup

set -e

BLUE='\033[36m'
GREEN='\033[32m'
RED='\033[31m'
YELLOW='\033[33m'
RESET='\033[0m'

if [ -z "$1" ]; then
    echo -e "${RED}Error: Backup file not specified${RESET}"
    echo "Usage: $0 <backup_file.tar.gz>"
    echo ""
    echo "Available backups:"
    ls -lh ./backups/*.tar.gz 2>/dev/null || echo "  No backups found"
    exit 1
fi

BACKUP_FILE=$1
BACKUP_DIR="./backups"
TEMP_DIR="$BACKUP_DIR/temp_restore"

echo -e "${YELLOW}WARNING: This will overwrite existing data!${RESET}"
echo -e "${YELLOW}Press Ctrl+C to cancel, or Enter to continue...${RESET}"
read

echo ""
echo -e "${BLUE}E-Commerce Platform Restore${RESET}"
echo "=================================="
echo "Backup file: $BACKUP_FILE"
echo ""

# Extract backup
echo -e "${BLUE}Extracting backup...${RESET}"
mkdir -p "$TEMP_DIR"
tar -xzf "$BACKUP_FILE" -C "$TEMP_DIR"

# Find extracted directory
RESTORE_PATH=$(find "$TEMP_DIR" -mindepth 1 -maxdepth 1 -type d | head -1)

if [ -z "$RESTORE_PATH" ]; then
    echo -e "${RED}Error: Could not find extracted backup${RESET}"
    exit 1
fi

echo "Restore path: $RESTORE_PATH"
echo ""

# Stop services
echo -e "${BLUE}Stopping services...${RESET}"
docker-compose stop auth-service product-service cart-service wishlist-service order-service payment-service inventory-service notification-service reporting-service

# Restore databases
echo ""
echo -e "${BLUE}Restoring databases...${RESET}"

echo "- Auth database"
docker-compose exec -T auth-db psql -U postgres -c "DROP DATABASE IF EXISTS auth_db;"
docker-compose exec -T auth-db psql -U postgres -c "CREATE DATABASE auth_db;"
docker-compose exec -T auth-db psql -U postgres auth_db < "$RESTORE_PATH/auth_db.sql"

echo "- Product database"
docker-compose exec -T product-db psql -U postgres -c "DROP DATABASE IF EXISTS product_db;"
docker-compose exec -T product-db psql -U postgres -c "CREATE DATABASE product_db;"
docker-compose exec -T product-db psql -U postgres product_db < "$RESTORE_PATH/product_db.sql"

echo "- Cart database"
docker-compose exec -T cart-db psql -U postgres -c "DROP DATABASE IF EXISTS cart_db;"
docker-compose exec -T cart-db psql -U postgres -c "CREATE DATABASE cart_db;"
docker-compose exec -T cart-db psql -U postgres cart_db < "$RESTORE_PATH/cart_db.sql"

echo "- Wishlist database"
docker-compose exec -T wishlist-db psql -U postgres -c "DROP DATABASE IF EXISTS wishlist_db;"
docker-compose exec -T wishlist-db psql -U postgres -c "CREATE DATABASE wishlist_db;"
docker-compose exec -T wishlist-db psql -U postgres wishlist_db < "$RESTORE_PATH/wishlist_db.sql"

echo "- Order database"
docker-compose exec -T order-db psql -U postgres -c "DROP DATABASE IF EXISTS order_db;"
docker-compose exec -T order-db psql -U postgres -c "CREATE DATABASE order_db;"
docker-compose exec -T order-db psql -U postgres order_db < "$RESTORE_PATH/order_db.sql"

echo "- Payment database"
docker-compose exec -T payment-db psql -U postgres -c "DROP DATABASE IF EXISTS payment_db;"
docker-compose exec -T payment-db psql -U postgres -c "CREATE DATABASE payment_db;"
docker-compose exec -T payment-db psql -U postgres payment_db < "$RESTORE_PATH/payment_db.sql"

echo "- Inventory database"
docker-compose exec -T inventory-db psql -U postgres -c "DROP DATABASE IF EXISTS inventory_db;"
docker-compose exec -T inventory-db psql -U postgres -c "CREATE DATABASE inventory_db;"
docker-compose exec -T inventory-db psql -U postgres inventory_db < "$RESTORE_PATH/inventory_db.sql"

echo "- Notification database"
docker-compose exec -T notification-db psql -U postgres -c "DROP DATABASE IF EXISTS notification_db;"
docker-compose exec -T notification-db psql -U postgres -c "CREATE DATABASE notification_db;"
docker-compose exec -T notification-db psql -U postgres notification_db < "$RESTORE_PATH/notification_db.sql"

echo "- Reporting database"
docker-compose exec -T reporting-db psql -U postgres -c "DROP DATABASE IF EXISTS reporting_db;"
docker-compose exec -T reporting-db psql -U postgres -c "CREATE DATABASE reporting_db;"
docker-compose exec -T reporting-db psql -U postgres reporting_db < "$RESTORE_PATH/reporting_db.sql"

# Restore Redis
echo ""
echo -e "${BLUE}Restoring Redis...${RESET}"
docker-compose stop redis
docker cp "$RESTORE_PATH/redis_dump.rdb" $(docker-compose ps -q redis):/data/dump.rdb
docker-compose start redis

# Clean up
echo ""
echo -e "${BLUE}Cleaning up...${RESET}"
rm -rf "$TEMP_DIR"

# Start services
echo ""
echo -e "${BLUE}Starting services...${RESET}"
docker-compose up -d

echo ""
echo -e "${GREEN}Restore completed successfully!${RESET}"
echo "All services have been restarted."
