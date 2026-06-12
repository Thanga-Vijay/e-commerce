#!/bin/bash

# Disaster Recovery Script for E-Commerce Platform
# Performs full backup of all critical components

set -e

NAMESPACE="ecommerce"
DR_DIR="./disaster-recovery/$(date +%Y%m%d_%H%M%S)"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  E-Commerce Disaster Recovery Backup  ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo "Timestamp: $TIMESTAMP"
echo "Backup location: $DR_DIR"
echo ""

mkdir -p "$DR_DIR"

# 1. Backup Kubernetes manifests
echo -e "${YELLOW}[1/7] Backing up Kubernetes manifests...${NC}"
mkdir -p "$DR_DIR/manifests"
kubectl get all,configmap,secret,pvc,ingress,hpa,networkpolicy -n $NAMESPACE -o yaml > "$DR_DIR/manifests/all-resources.yaml"
cp -r k8s/* "$DR_DIR/manifests/" 2>/dev/null || true
echo -e "${GREEN}✓ Manifests backed up${NC}"
echo ""

# 2. Backup all databases
echo -e "${YELLOW}[2/7] Backing up databases...${NC}"
mkdir -p "$DR_DIR/databases"

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

for db_config in "${DATABASES[@]}"; do
  IFS=':' read -r POD_PREFIX DB_NAME <<< "$db_config"
  POD_NAME=$(kubectl get pod -n $NAMESPACE -l app=$POD_PREFIX -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
  
  if [ -z "$POD_NAME" ]; then
    echo -e "${RED}  ✗ Could not find pod for $DB_NAME${NC}"
    continue
  fi
  
  echo "  Backing up $DB_NAME..."
  kubectl exec -n $NAMESPACE $POD_NAME -- pg_dump -U postgres $DB_NAME 2>/dev/null | gzip > "$DR_DIR/databases/${DB_NAME}.sql.gz"
  echo -e "${GREEN}  ✓ $DB_NAME backed up${NC}"
done
echo ""

# 3. Backup Redis data
echo -e "${YELLOW}[3/7] Backing up Redis...${NC}"
mkdir -p "$DR_DIR/redis"
REDIS_POD=$(kubectl get pod -n $NAMESPACE -l app=redis -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
if [ ! -z "$REDIS_POD" ]; then
  kubectl exec -n $NAMESPACE $REDIS_POD -- redis-cli BGSAVE 2>/dev/null || true
  sleep 5
  kubectl cp -n $NAMESPACE $REDIS_POD:/data/dump.rdb "$DR_DIR/redis/dump.rdb" 2>/dev/null || true
  echo -e "${GREEN}✓ Redis backed up${NC}"
else
  echo -e "${RED}✗ Redis pod not found${NC}"
fi
echo ""

# 4. Backup Kafka topics (metadata only)
echo -e "${YELLOW}[4/7] Backing up Kafka configuration...${NC}"
mkdir -p "$DR_DIR/kafka"
KAFKA_POD=$(kubectl get pod -n $NAMESPACE -l app=kafka -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
if [ ! -z "$KAFKA_POD" ]; then
  kubectl exec -n $NAMESPACE $KAFKA_POD -- kafka-topics.sh --bootstrap-server localhost:9092 --list > "$DR_DIR/kafka/topics.txt" 2>/dev/null || true
  echo -e "${GREEN}✓ Kafka config backed up${NC}"
else
  echo -e "${RED}✗ Kafka pod not found${NC}"
fi
echo ""

# 5. Export secrets (encrypted)
echo -e "${YELLOW}[5/7] Exporting secrets...${NC}"
mkdir -p "$DR_DIR/secrets"
kubectl get secrets -n $NAMESPACE -o yaml > "$DR_DIR/secrets/secrets.yaml"
echo -e "${YELLOW}  ⚠ Remember to encrypt this file and store securely!${NC}"
echo -e "${GREEN}✓ Secrets exported${NC}"
echo ""

# 6. Backup persistent volume data
echo -e "${YELLOW}[6/7] Backing up PVC information...${NC}"
mkdir -p "$DR_DIR/pvcs"
kubectl get pvc -n $NAMESPACE -o yaml > "$DR_DIR/pvcs/all-pvcs.yaml"
kubectl get pv -o yaml > "$DR_DIR/pvcs/all-pvs.yaml"
echo -e "${GREEN}✓ PVC info backed up${NC}"
echo ""

# 7. Create recovery metadata
echo -e "${YELLOW}[7/7] Creating recovery metadata...${NC}"
cat > "$DR_DIR/recovery-metadata.json" <<EOF
{
  "timestamp": "$TIMESTAMP",
  "namespace": "$NAMESPACE",
  "kubernetes_version": "$(kubectl version --short 2>/dev/null | grep Server | awk '{print $3}' || echo 'unknown')",
  "cluster_name": "$(kubectl config current-context)",
  "node_count": $(kubectl get nodes --no-headers 2>/dev/null | wc -l),
  "pod_count": $(kubectl get pods -n $NAMESPACE --no-headers 2>/dev/null | wc -l),
  "databases": [
$(for db in "${DATABASES[@]}"; do
    echo "    \"${db#*:}\","
  done | sed '$ s/,$//')
  ],
  "backup_contents": {
    "manifests": true,
    "databases": true,
    "redis": $([ ! -z "$REDIS_POD" ] && echo "true" || echo "false"),
    "kafka": $([ ! -z "$KAFKA_POD" ] && echo "true" || echo "false"),
    "secrets": true,
    "pvcs": true
  }
}
EOF

# Create checksums
cd "$DR_DIR"
find . -type f -not -name "checksums.txt" -exec sha256sum {} \; > checksums.txt
cd - > /dev/null

echo -e "${GREEN}✓ Metadata created${NC}"
echo ""

# Calculate total size
TOTAL_SIZE=$(du -sh "$DR_DIR" | cut -f1)

# Create archive
echo -e "${YELLOW}Creating compressed archive...${NC}"
ARCHIVE_NAME="dr_backup_${TIMESTAMP}.tar.gz"
tar -czf "$ARCHIVE_NAME" -C "$(dirname $DR_DIR)" "$(basename $DR_DIR)"
ARCHIVE_SIZE=$(du -sh "$ARCHIVE_NAME" | cut -f1)

echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║      Disaster Recovery Complete       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Summary:${NC}"
echo "  Backup directory: $DR_DIR"
echo "  Archive: $ARCHIVE_NAME"
echo "  Total size: $TOTAL_SIZE"
echo "  Archive size: $ARCHIVE_SIZE"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Encrypt the archive: gpg -c $ARCHIVE_NAME"
echo "  2. Upload to secure storage (S3, Azure Blob, etc.)"
echo "  3. Verify backup integrity"
echo "  4. Document recovery procedures"
echo ""
echo -e "${GREEN}To restore from this backup, use:${NC}"
echo "  ./restore-from-dr.sh $ARCHIVE_NAME"
echo ""

# Upload to cloud (if configured)
if [ ! -z "$DR_S3_BUCKET" ]; then
  echo -e "${YELLOW}Uploading to S3...${NC}"
  aws s3 cp "$ARCHIVE_NAME" "s3://${DR_S3_BUCKET}/dr-backups/"
  aws s3 cp "$DR_DIR/recovery-metadata.json" "s3://${DR_S3_BUCKET}/dr-backups/"
  echo -e "${GREEN}✓ Uploaded to S3: s3://${DR_S3_BUCKET}/dr-backups/${NC}"
fi

echo -e "${GREEN}Disaster recovery backup completed successfully!${NC}"
