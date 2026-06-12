#!/bin/bash

# Restore a PostgreSQL database from backup
# Usage: ./restore-database.sh <backup-file> <database-name>

set -e

if [ $# -lt 2 ]; then
  echo "Usage: $0 <backup-file> <database-name>"
  echo "Example: $0 backups/20240612/auth_db_20240612_120000.sql.gz auth_db"
  exit 1
fi

BACKUP_FILE=$1
DB_NAME=$2
NAMESPACE="ecommerce"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}=== Database Restore ===${NC}"
echo "Backup file: $BACKUP_FILE"
echo "Database: $DB_NAME"
echo ""

# Validate backup file
if [ ! -f "$BACKUP_FILE" ]; then
  echo -e "${RED}Error: Backup file not found: $BACKUP_FILE${NC}"
  exit 1
fi

# Determine pod name based on database
case $DB_NAME in
  auth_db) POD_PREFIX="auth-db" ;;
  product_db) POD_PREFIX="product-db" ;;
  cart_db) POD_PREFIX="cart-db" ;;
  wishlist_db) POD_PREFIX="wishlist-db" ;;
  order_db) POD_PREFIX="order-db" ;;
  payment_db) POD_PREFIX="payment-db" ;;
  inventory_db) POD_PREFIX="inventory-db" ;;
  notification_db) POD_PREFIX="notification-db" ;;
  reporting_db) POD_PREFIX="reporting-db" ;;
  *)
    echo -e "${RED}Error: Unknown database: $DB_NAME${NC}"
    exit 1
    ;;
esac

# Get pod name
POD_NAME=$(kubectl get pod -n $NAMESPACE -l app=$POD_PREFIX -o jsonpath='{.items[0].metadata.name}')

if [ -z "$POD_NAME" ]; then
  echo -e "${RED}Error: Could not find pod for $POD_PREFIX${NC}"
  exit 1
fi

echo "Pod: $POD_NAME"
echo ""

# Confirm restore
echo -e "${RED}WARNING: This will overwrite the existing database!${NC}"
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
  echo "Restore cancelled"
  exit 0
fi

echo ""
echo "Starting restore..."

# Decompress if needed
TEMP_FILE="/tmp/restore_${DB_NAME}_$$.sql"
if [[ $BACKUP_FILE == *.gz ]]; then
  echo "Decompressing backup..."
  gunzip -c "$BACKUP_FILE" > "$TEMP_FILE"
else
  cp "$BACKUP_FILE" "$TEMP_FILE"
fi

# Scale down dependent services
echo "Scaling down services using this database..."
SERVICE_NAME="${POD_PREFIX%-db}-service"
ORIGINAL_REPLICAS=$(kubectl get deployment $SERVICE_NAME -n $NAMESPACE -o jsonpath='{.spec.replicas}' 2>/dev/null || echo "0")

if [ "$ORIGINAL_REPLICAS" != "0" ]; then
  kubectl scale deployment $SERVICE_NAME --replicas=0 -n $NAMESPACE 2>/dev/null || true
  echo "Waiting for pods to terminate..."
  sleep 10
fi

# Drop existing database (optional - comment out if you want to merge)
echo "Dropping existing database..."
kubectl exec -n $NAMESPACE $POD_NAME -- psql -U postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"

# Create fresh database
echo "Creating fresh database..."
kubectl exec -n $NAMESPACE $POD_NAME -- psql -U postgres -c "CREATE DATABASE $DB_NAME;"

# Restore from backup
echo "Restoring data..."
cat "$TEMP_FILE" | kubectl exec -i -n $NAMESPACE $POD_NAME -- psql -U postgres $DB_NAME

# Clean up temp file
rm -f "$TEMP_FILE"

# Scale services back up
if [ "$ORIGINAL_REPLICAS" != "0" ]; then
  echo "Scaling services back to $ORIGINAL_REPLICAS replicas..."
  kubectl scale deployment $SERVICE_NAME --replicas=$ORIGINAL_REPLICAS -n $NAMESPACE 2>/dev/null || true
fi

echo ""
echo -e "${GREEN}=== Restore Complete ===${NC}"
echo "Database $DB_NAME has been restored successfully"

# Verify restore
echo ""
echo "Verifying restore..."
TABLE_COUNT=$(kubectl exec -n $NAMESPACE $POD_NAME -- psql -U postgres $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
echo "Tables restored: $(echo $TABLE_COUNT | tr -d ' ')"

echo ""
echo -e "${GREEN}Restore completed successfully!${NC}"
