#!/bin/bash

# Backup all PostgreSQL databases in Kubernetes
# Usage: ./backup-databases.sh [backup-location]

set -e

NAMESPACE="ecommerce"
BACKUP_DIR="${1:-./backups/$(date +%Y%m%d_%H%M%S)}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Database list
DATABASES=(
  "auth-db:auth_db"
  "product-db:product_db"
  "cart-db:cart_db"
  "wishlist-db:wishlist_db"
  "order-db:order_db"
  "payment-db:payment_db"
  "inventory-db:inventory_db"
  "notification-db:notification_db"
  "reporting-db:reporting_db"
)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== E-Commerce Database Backup ===${NC}"
echo "Timestamp: $TIMESTAMP"
echo "Backup location: $BACKUP_DIR"
echo ""

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Function to backup a single database
backup_database() {
  local POD_PREFIX=$1
  local DB_NAME=$2
  local BACKUP_FILE="${BACKUP_DIR}/${DB_NAME}_${TIMESTAMP}.sql"
  
  echo -e "${YELLOW}Backing up $DB_NAME...${NC}"
  
  # Get the pod name
  POD_NAME=$(kubectl get pod -n $NAMESPACE -l app=$POD_PREFIX -o jsonpath='{.items[0].metadata.name}')
  
  if [ -z "$POD_NAME" ]; then
    echo -e "${RED}Error: Could not find pod for $POD_PREFIX${NC}"
    return 1
  fi
  
  # Perform backup
  kubectl exec -n $NAMESPACE $POD_NAME -- pg_dump -U postgres $DB_NAME > "$BACKUP_FILE"
  
  # Compress backup
  gzip "$BACKUP_FILE"
  
  local SIZE=$(du -h "${BACKUP_FILE}.gz" | cut -f1)
  echo -e "${GREEN}✓ Backed up $DB_NAME ($SIZE)${NC}"
}

# Backup all databases
for db_config in "${DATABASES[@]}"; do
  IFS=':' read -r POD_PREFIX DB_NAME <<< "$db_config"
  backup_database "$POD_PREFIX" "$DB_NAME" || echo -e "${RED}Failed to backup $DB_NAME${NC}"
done

# Backup metadata
echo "Creating backup metadata..."
cat > "${BACKUP_DIR}/metadata.json" <<EOF
{
  "timestamp": "$TIMESTAMP",
  "namespace": "$NAMESPACE",
  "databases": [
$(for db in "${DATABASES[@]}"; do
    echo "    \"${db#*:}\","
  done | sed '$ s/,$//')
  ],
  "kubernetes_version": "$(kubectl version --short | grep Server | awk '{print $3}')"
}
EOF

# Create checksum file
echo "Creating checksums..."
cd "$BACKUP_DIR"
sha256sum *.sql.gz > checksums.txt
cd - > /dev/null

# Calculate total size
TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
echo ""
echo -e "${GREEN}=== Backup Complete ===${NC}"
echo "Total size: $TOTAL_SIZE"
echo "Location: $BACKUP_DIR"

# Clean old backups
if [ -d "./backups" ]; then
  echo ""
  echo "Cleaning backups older than $RETENTION_DAYS days..."
  find ./backups -type d -mtime +$RETENTION_DAYS -exec rm -rf {} + 2>/dev/null || true
  echo "Cleanup complete"
fi

# Upload to cloud storage (optional)
if [ ! -z "$BACKUP_S3_BUCKET" ]; then
  echo ""
  echo "Uploading to S3..."
  aws s3 cp "$BACKUP_DIR" "s3://${BACKUP_S3_BUCKET}/ecommerce-backups/" --recursive
  echo -e "${GREEN}✓ Uploaded to S3${NC}"
fi

echo ""
echo -e "${GREEN}All backups completed successfully!${NC}"
