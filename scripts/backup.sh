#!/bin/bash

# Backup Script for E-Commerce Platform
# Creates backups of all databases and volumes

set -e

BLUE='\033[36m'
GREEN='\033[32m'
RESET='\033[0m'

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="ecommerce_backup_$TIMESTAMP"
BACKUP_PATH="$BACKUP_DIR/$BACKUP_NAME"

echo -e "${BLUE}E-Commerce Platform Backup${RESET}"
echo "=================================="
echo "Backup directory: $BACKUP_PATH"
echo ""

# Create backup directory
mkdir -p "$BACKUP_PATH"

# Backup databases
echo -e "${BLUE}Backing up databases...${RESET}"

echo "- Auth database"
docker-compose exec -T auth-db pg_dump -U postgres auth_db > "$BACKUP_PATH/auth_db.sql"

echo "- Product database"
docker-compose exec -T product-db pg_dump -U postgres product_db > "$BACKUP_PATH/product_db.sql"

echo "- Cart database"
docker-compose exec -T cart-db pg_dump -U postgres cart_db > "$BACKUP_PATH/cart_db.sql"

echo "- Wishlist database"
docker-compose exec -T wishlist-db pg_dump -U postgres wishlist_db > "$BACKUP_PATH/wishlist_db.sql"

echo "- Order database"
docker-compose exec -T order-db pg_dump -U postgres order_db > "$BACKUP_PATH/order_db.sql"

echo "- Payment database"
docker-compose exec -T payment-db pg_dump -U postgres payment_db > "$BACKUP_PATH/payment_db.sql"

echo "- Inventory database"
docker-compose exec -T inventory-db pg_dump -U postgres inventory_db > "$BACKUP_PATH/inventory_db.sql"

echo "- Notification database"
docker-compose exec -T notification-db pg_dump -U postgres notification_db > "$BACKUP_PATH/notification_db.sql"

echo "- Reporting database"
docker-compose exec -T reporting-db pg_dump -U postgres reporting_db > "$BACKUP_PATH/reporting_db.sql"

# Backup Redis data
echo ""
echo -e "${BLUE}Backing up Redis...${RESET}"
docker-compose exec -T redis redis-cli --rdb /data/dump.rdb SAVE
docker cp $(docker-compose ps -q redis):/data/dump.rdb "$BACKUP_PATH/redis_dump.rdb"

# Create metadata file
echo ""
echo -e "${BLUE}Creating backup metadata...${RESET}"
cat > "$BACKUP_PATH/backup_info.txt" <<EOF
E-Commerce Platform Backup
Date: $(date)
Timestamp: $TIMESTAMP
Databases: 9 (auth, product, cart, wishlist, order, payment, inventory, notification, reporting)
Redis: Yes
EOF

# Compress backup
echo ""
echo -e "${BLUE}Compressing backup...${RESET}"
cd "$BACKUP_DIR"
tar -czf "$BACKUP_NAME.tar.gz" "$BACKUP_NAME"
rm -rf "$BACKUP_NAME"
cd - > /dev/null

echo ""
echo -e "${GREEN}Backup completed successfully!${RESET}"
echo "Backup file: $BACKUP_DIR/$BACKUP_NAME.tar.gz"
echo "Size: $(du -h "$BACKUP_DIR/$BACKUP_NAME.tar.gz" | cut -f1)"
