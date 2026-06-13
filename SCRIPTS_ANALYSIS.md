# E-Commerce Platform Scripts Analysis

Comprehensive analysis of all scripts for KIND local testing and AWS EKS production deployment.

## Overview

The platform includes 10 scripts across two environments:
- **Docker Compose** (Local Development): 4 scripts
- **Kubernetes** (KIND/EKS): 6 scripts
- **New**: KIND-specific scripts for local K8s testing

---

## Script Categories

### Category 1: Docker Compose Scripts (Development)

#### 1. `scripts/backup.sh`
**Purpose:** Backup all databases and Redis for Docker Compose environment

**Key Features:**
- Backs up 9 PostgreSQL databases
- Backs up Redis RDB file
- Creates compressed tar.gz archive
- Generates metadata file

**Usage:**
```bash
./scripts/backup.sh
# Creates: ./backups/ecommerce_backup_YYYYMMDD_HHMMSS.tar.gz
```

**Requirements:**
- Docker Compose running
- All database containers operational
- Sufficient disk space

**Compatibility:**
- ✅ Local Docker Compose
- ❌ Kubernetes/KIND (use backup-databases.sh instead)
- ❌ AWS EKS (use Velero or AWS Backup)

---

#### 2. `scripts/restore.sh`
**Purpose:** Restore databases from Docker Compose backup

**Key Features:**
- Extracts backup archive
- Stops services during restore
- Drops and recreates databases
- Restores all 9 databases
- Restarts services

**Usage:**
```bash
./scripts/restore.sh backups/ecommerce_backup_20240612_120000.tar.gz
```

**Warnings:**
- ⚠️ Destructive operation - drops existing databases
- ⚠️ Requires confirmation before proceeding
- ⚠️ Services are stopped during restore

**Compatibility:**
- ✅ Local Docker Compose
- ❌ Kubernetes/KIND
- ❌ AWS EKS

---

#### 3. `scripts/setup-secrets.sh`
**Purpose:** Generate secure secrets for Docker Compose production deployment

**Key Features:**
- Generates random PostgreSQL password (32 bytes)
- Generates random Redis password (32 bytes)
- Generates JWT secret (64 bytes)
- Prompts for Stripe API keys
- Prompts for SMTP password
- Creates secret files with 600 permissions

**Usage:**
```bash
./scripts/setup-secrets.sh
```

**Outputs:** (in `./secrets/` directory)
- `postgres_password.txt`
- `redis_password.txt`
- `jwt_secret.txt`
- `stripe_secret_key.txt`
- `stripe_webhook_secret.txt`
- `smtp_password.txt`

**Security:**
- All files created with 600 permissions (owner read/write only)
- Directory created with 700 permissions
- Never commit these files to git

**Compatibility:**
- ✅ Local Docker Compose
- ⚠️ Manual adaptation needed for K8s Secrets (base64 encode values)

---

#### 4. `scripts/health-check.sh`
**Purpose:** Check health status of all services

**Key Features:**
- Checks HTTP endpoints with curl
- Verifies database connectivity
- Color-coded output (✓ ✗ ⚠)
- Reports HTTP status codes

**Usage:**
```bash
./scripts/health-check.sh
```

**Checks:**
- Auth Service: http://localhost:8081/health
- Product Service: http://localhost:8082/health
- Cart Service: http://localhost:8083/health
- Wishlist Service: http://localhost:8084/health
- Order Service: http://localhost:8085/health
- Payment Service: http://localhost:8086/health
- Inventory Service: http://localhost:8087/health
- Notification Service: http://localhost:8088/health
- Reporting Service: http://localhost:8089/health
- Frontend: http://localhost:3000
- All 9 databases

**Compatibility:**
- ✅ Local Docker Compose
- ⚠️ Kubernetes/KIND (needs port-forward or ingress URLs)
- ⚠️ AWS EKS (needs Load Balancer URLs)

**Adaptation for K8s:**
```bash
# Option 1: Port forward
kubectl port-forward -n ecommerce svc/auth-service 8081:8081

# Option 2: Update script to use Ingress URLs
AUTH_URL="http://api.yourdomain.com/api/v1/auth/health"
```

---

### Category 2: Kubernetes Scripts (KIND/EKS)

#### 5. `k8s/deploy.sh`
**Purpose:** Deploy entire e-commerce platform to Kubernetes

**Deployment Order:**
1. Namespace (`ecommerce`)
2. ConfigMaps and Secrets
3. PostgreSQL databases (9 StatefulSets)
4. Redis (1 StatefulSet)
5. Kafka + Zookeeper (StatefulSets)
6. Microservices (9 Deployments)
7. Frontend (1 Deployment)
8. Ingress
9. HorizontalPodAutoscalers

**Usage:**
```bash
cd k8s
./deploy.sh
```

**Prerequisites:**
- kubectl installed
- Cluster accessible (KIND or EKS)
- Secrets updated with real credentials

**Wait Times:**
- Databases: 300s timeout
- Redis: 120s timeout
- Kafka: 30s sleep
- Services: 180s timeout

**Compatibility:**
- ✅ KIND (local Kubernetes)
- ✅ AWS EKS
- ✅ Any Kubernetes cluster
- ⚠️ Requires Ingress controller installed

**KIND-Specific Notes:**
- Ingress NGINX installed by kind-setup.sh
- Metrics Server patched for KIND by kind-setup.sh
- Uses local-path-provisioner for storage

**EKS-Specific Notes:**
- Requires AWS Load Balancer Controller
- Requires EBS CSI Driver
- Update images to use ECR registry
- Update Ingress to use ALB

---

#### 6. `k8s/undeploy.sh`
**Purpose:** Tear down entire e-commerce platform from Kubernetes

**Deletion Order:**
1. HPA (Horizontal Pod Autoscalers)
2. Ingress
3. Frontend
4. Microservices
5. Kafka
6. Redis
7. Databases
8. ConfigMaps and Secrets
9. Namespace

**Usage:**
```bash
cd k8s
./undeploy.sh
```

**Warnings:**
- ⚠️ Deletes all resources in `ecommerce` namespace
- ⚠️ PersistentVolumeClaims are deleted (data loss)
- ⚠️ Requires confirmation before proceeding

**Data Preservation:**
Before undeploying, backup:
```bash
./scripts/backup-databases.sh
./scripts/disaster-recovery.sh
```

**Compatibility:**
- ✅ KIND
- ✅ AWS EKS
- ✅ Any Kubernetes cluster

---

#### 7. `scripts/backup-databases.sh`
**Purpose:** Backup all 9 PostgreSQL databases in Kubernetes

**Key Features:**
- Backs up all 9 databases via kubectl exec
- Uses pg_dump for PostgreSQL backups
- Compresses backups with gzip
- Creates metadata JSON file
- 30-day retention policy

**Usage:**
```bash
./scripts/backup-databases.sh [backup-location]
# Default: ./backups/YYYYMMDD_HHMMSS/
```

**Output Structure:**
```
backups/20240612_120000/
├── auth_db_20240612_120000.sql.gz
├── product_db_20240612_120000.sql.gz
├── cart_db_20240612_120000.sql.gz
├── wishlist_db_20240612_120000.sql.gz
├── order_db_20240612_120000.sql.gz
├── payment_db_20240612_120000.sql.gz
├── inventory_db_20240612_120000.sql.gz
├── notification_db_20240612_120000.sql.gz
├── reporting_db_20240612_120000.sql.gz
└── metadata.json
```

**Database Mapping:**
- `auth-db` pod → `auth_db` database
- `product-db` pod → `product_db` database
- `cart-db` pod → `cart_db` database
- `wishlist-db` pod → `wishlist_db` database
- `order-db` pod → `order_db` database
- `payment-db` pod → `payment_db` database
- `inventory-db` pod → `inventory_db` database
- `notification-db` pod → `notification_db` database
- `reporting-db` pod → `reporting_db` database

**Compatibility:**
- ✅ KIND
- ✅ AWS EKS
- ✅ Any Kubernetes cluster
- ✅ Can be run manually or as CronJob

**EKS Enhancement:**
Add S3 upload:
```bash
# At end of script
aws s3 cp "$BACKUP_DIR" s3://my-backups/databases/ --recursive
```

---

#### 8. `scripts/restore-database.sh`
**Purpose:** Restore a single database from backup

**Key Features:**
- Restores one database at a time
- Supports gzip-compressed backups
- Validates backup file existence
- Requires confirmation before restore
- Scales down services during restore (recommended)

**Usage:**
```bash
./scripts/restore-database.sh backups/20240612_120000/auth_db_20240612_120000.sql.gz auth_db
```

**Parameters:**
1. `<backup-file>`: Path to SQL backup (plain or .gz)
2. `<database-name>`: Target database (auth_db, product_db, etc.)

**Supported Databases:**
- `auth_db`
- `product_db`
- `cart_db`
- `wishlist_db`
- `order_db`
- `payment_db`
- `inventory_db`
- `notification_db`
- `reporting_db`

**Process:**
1. Validates backup file
2. Determines pod name from database name
3. Gets pod from Kubernetes
4. Decompresses if needed
5. Copies SQL to pod
6. Drops and recreates database
7. Restores from SQL dump
8. Cleans up temp files

**Safety:**
- Requires explicit "yes" confirmation
- Shows warning about data overwrite
- Validates pod existence

**Compatibility:**
- ✅ KIND
- ✅ AWS EKS
- ✅ Any Kubernetes cluster

**Best Practice:**
Scale down services first:
```bash
kubectl scale deployment auth-service -n ecommerce --replicas=0
./scripts/restore-database.sh backup.sql.gz auth_db
kubectl scale deployment auth-service -n ecommerce --replicas=2
```

---

#### 9. `scripts/disaster-recovery.sh`
**Purpose:** Complete disaster recovery backup of entire platform

**Backup Components:**
1. Kubernetes manifests (all resources YAML)
2. All 9 PostgreSQL databases
3. Redis RDB file
4. Kafka topics and data
5. PersistentVolumeClaim snapshots
6. ConfigMaps and Secrets
7. Application logs

**Usage:**
```bash
./scripts/disaster-recovery.sh
# Creates: ./disaster-recovery/YYYYMMDD_HHMMSS/
```

**Output Structure:**
```
disaster-recovery/20240612_120000/
├── manifests/
│   ├── all-resources.yaml
│   └── k8s/ (copy of k8s directory)
├── databases/
│   ├── auth_db.sql.gz
│   ├── product_db.sql.gz
│   └── ...
├── redis/
│   └── dump.rdb
├── kafka/
│   └── topics.json
├── pvcs/
│   └── snapshot-list.json
├── logs/
│   ├── auth-service.log
│   └── ...
└── metadata.json
```

**Metadata File:**
```json
{
  "timestamp": "20240612_120000",
  "namespace": "ecommerce",
  "cluster_info": "...",
  "components": [
    "manifests",
    "databases",
    "redis",
    "kafka",
    "pvcs",
    "logs"
  ]
}
```

**Features:**
- [1/7] Kubernetes manifests (kubectl get all)
- [2/7] Database backups (pg_dump)
- [3/7] Redis backup (BGSAVE + copy RDB)
- [4/7] Kafka backup (topic list + configs)
- [5/7] PVC snapshots (metadata)
- [6/7] Application logs (last 1000 lines per pod)
- [7/7] Create compressed archive

**Compatibility:**
- ✅ KIND
- ✅ AWS EKS
- ✅ Any Kubernetes cluster

**EKS Enhancements:**
```bash
# Add AWS-specific backups
# 1. EBS snapshots
aws ec2 create-snapshots --instance-specification ...

# 2. Upload to S3
aws s3 cp "$DR_DIR" s3://my-dr-bucket/ --recursive

# 3. Store in Glacier for long-term
aws s3 cp ... --storage-class GLACIER
```

**Automation:**
Run as CronJob (already created in k8s/backup/backup-cronjobs.yaml):
```yaml
- Daily backup at 2 AM
- Weekly full DR backup (Sunday 3 AM)
- 30-day retention
```

---

### Category 3: KIND-Specific Scripts (NEW)

#### 10. `kind-setup.sh` ⭐ NEW
**Purpose:** Create and configure KIND cluster for local K8s testing

**Cluster Configuration:**
- Name: `ecommerce-local`
- Nodes: 4 (1 control-plane + 3 workers)
- Kubernetes Version: 1.28
- Container Runtime: containerd

**Node Labels:**
- `ecommerce-local-worker`: `tier=application`
- `ecommerce-local-worker2`: `tier=application`
- `ecommerce-local-worker3`: `tier=database`

**Port Mappings:**
| Service | Container Port | Host Port |
|---------|---------------|-----------|
| HTTP Ingress | 80 | 80 |
| HTTPS Ingress | 443 | 443 |
| Frontend | 30000 | 3000 |
| Grafana | 30001 | 3001 |
| Prometheus | 30002 | 9090 |
| Jaeger UI | 30003 | 16686 |
| Kafka UI | 30004 | 8090 |

**Installed Components:**
1. **Ingress NGINX Controller**
   - For HTTP/HTTPS routing
   - KIND-specific deployment

2. **Metrics Server**
   - For HorizontalPodAutoscaler
   - Patched with `--kubelet-insecure-tls` for KIND

3. **Local Container Registry**
   - Registry: `localhost:5000`
   - Connected to KIND network
   - For faster image loading

4. **StorageClass**
   - Provider: `rancher.io/local-path`
   - Default storage class
   - Volume binding: `WaitForFirstConsumer`

**Usage:**
```bash
./kind-setup.sh
```

**Prerequisites:**
- KIND installed (v0.20.0+)
- kubectl installed
- Docker Desktop running

**Duration:** 2-3 minutes

**Idempotency:**
- Checks if cluster exists
- Prompts to delete and recreate
- Or uses existing cluster

**Next Steps:**
```bash
# 1. Build and load images
./scripts/build-and-load-images.sh

# 2. Deploy platform
cd k8s && ./deploy.sh

# 3. Deploy monitoring
cd k8s/monitoring && ./deploy-monitoring.sh
```

---

#### 11. `scripts/build-and-load-images.sh` ⭐ NEW
**Purpose:** Build Docker images and load them into KIND cluster

**Images Built:**
- `ecommerce/auth-service:latest`
- `ecommerce/product-service:latest`
- `ecommerce/cart-service:latest`
- `ecommerce/wishlist-service:latest`
- `ecommerce/order-service:latest`
- `ecommerce/payment-service:latest`
- `ecommerce/inventory-service:latest`
- `ecommerce/notification-service:latest`
- `ecommerce/reporting-service:latest`
- `ecommerce/frontend:latest`

**Process for Each Service:**
1. Build Docker image
2. Tag for local registry (`localhost:5000/...`)
3. Load image into KIND cluster
4. (Optional) Push to local registry

**Usage:**
```bash
# Build with default version (latest)
./scripts/build-and-load-images.sh

# Build with specific version
VERSION=v1.0.0 ./scripts/build-and-load-images.sh
```

**Build Locations:**
- Services: `services/<service-name>/Dockerfile`
- Frontend: `frontend/Dockerfile`

**Verification:**
```bash
# Check images in KIND
docker exec -it ecommerce-local-control-plane crictl images | grep ecommerce
```

**Performance:**
- Parallel builds not implemented (sequential)
- Total time: ~10-15 minutes for all services
- Uses Docker BuildKit for layer caching

**Optimization Tips:**
```bash
# Build only changed services
SERVICES=("auth-service" "product-service")
for service in "${SERVICES[@]}"; do
  docker build ...
  kind load docker-image ...
done
```

**Compatibility:**
- ✅ KIND
- ❌ AWS EKS (use ECR instead)

**For EKS:**
Use ECR and push:
```bash
# See eks-migration-guide.md for full instructions
aws ecr get-login-password | docker login ...
docker build -t $ECR_URL/ecommerce/auth-service:v1.0.0 ...
docker push $ECR_URL/ecommerce/auth-service:v1.0.0
```

---

## Script Compatibility Matrix

| Script | Docker Compose | KIND | EKS | Purpose |
|--------|---------------|------|-----|---------|
| `backup.sh` | ✅ | ❌ | ❌ | Docker backup |
| `restore.sh` | ✅ | ❌ | ❌ | Docker restore |
| `setup-secrets.sh` | ✅ | ⚠️ | ⚠️ | Generate secrets |
| `health-check.sh` | ✅ | ⚠️ | ⚠️ | Health checks |
| `deploy.sh` | ❌ | ✅ | ✅ | K8s deployment |
| `undeploy.sh` | ❌ | ✅ | ✅ | K8s teardown |
| `backup-databases.sh` | ❌ | ✅ | ✅ | K8s DB backup |
| `restore-database.sh` | ❌ | ✅ | ✅ | K8s DB restore |
| `disaster-recovery.sh` | ❌ | ✅ | ✅ | Full K8s backup |
| `kind-setup.sh` | ❌ | ✅ | ❌ | KIND cluster |
| `build-and-load-images.sh` | ❌ | ✅ | ❌ | KIND images |

Legend:
- ✅ Fully compatible
- ⚠️ Needs modification
- ❌ Not compatible

---

## Workflow: KIND Local Testing

### Step 1: Create KIND Cluster
```bash
./kind-setup.sh
```

**What it does:**
- Creates 4-node cluster
- Installs Ingress NGINX
- Installs Metrics Server
- Creates local registry
- Configures storage

**Duration:** 2-3 minutes

---

### Step 2: Build and Load Images
```bash
./scripts/build-and-load-images.sh
```

**What it does:**
- Builds 10 Docker images
- Loads them into KIND
- Tags for local registry

**Duration:** 10-15 minutes

---

### Step 3: Deploy Platform
```bash
cd k8s
./deploy.sh
```

**What it does:**
- Creates namespace
- Deploys databases (9)
- Deploys Redis
- Deploys Kafka
- Deploys services (9)
- Deploys frontend
- Deploys ingress
- Deploys HPA

**Duration:** 5-8 minutes

---

### Step 4: Deploy Monitoring (Optional)
```bash
cd k8s/monitoring
kubectl apply -f prometheus-config.yaml
kubectl apply -f prometheus.yaml
kubectl apply -f grafana.yaml
kubectl apply -f loki.yaml
kubectl apply -f jaeger.yaml
```

**What it does:**
- Installs Prometheus
- Installs Grafana
- Installs Loki (logs)
- Installs Jaeger (tracing)

**Duration:** 3-5 minutes

---

### Step 5: Verify Deployment
```bash
# Check pods
kubectl get pods -n ecommerce

# Check services
kubectl get svc -n ecommerce

# Check ingress
kubectl get ingress -n ecommerce
```

**Access Services:**
- Frontend: http://localhost:3000
- Grafana: http://localhost:3001
- Prometheus: http://localhost:9090
- Jaeger: http://localhost:16686

---

### Step 6: Testing
```bash
# Run health checks
./scripts/health-check.sh  # Needs modification for K8s

# Test API endpoints
curl http://localhost/api/v1/auth/health
curl http://localhost/api/v1/products/health

# View logs
kubectl logs -f deployment/auth-service -n ecommerce
```

---

### Step 7: Backup (Important!)
```bash
# Backup databases
./scripts/backup-databases.sh

# Full disaster recovery backup
./scripts/disaster-recovery.sh
```

---

### Step 8: Teardown (When Done)
```bash
# Delete application
cd k8s
./undeploy.sh

# Delete KIND cluster
kind delete cluster --name ecommerce-local
```

---

## Workflow: EKS Production Deployment

### Prerequisites
1. AWS CLI configured
2. eksctl installed
3. kubectl installed
4. helm installed

**See:** [eks-migration-guide.md](eks-migration-guide.md) for complete instructions

---

### Step 1: Create EKS Cluster
```bash
# Using eksctl
eksctl create cluster -f eks-cluster.yaml
```

**Duration:** 15-20 minutes

---

### Step 2: Install Add-ons
```bash
# AWS Load Balancer Controller
helm install aws-load-balancer-controller ...

# EBS CSI Driver
eksctl create addon --name aws-ebs-csi-driver ...

# Cert-Manager (for TLS)
kubectl apply -f https://github.com/cert-manager/cert-manager/...
```

**Duration:** 5-10 minutes

---

### Step 3: Create ECR Repositories
```bash
# Create repos for all services
for service in "${SERVICES[@]}"; do
  aws ecr create-repository --repository-name ecommerce/${service}
done
```

---

### Step 4: Build and Push Images
```bash
# Login to ECR
aws ecr get-login-password | docker login ...

# Build and push
for service in "${SERVICES[@]}"; do
  docker build -t $ECR_URL/ecommerce/${service}:v1.0.0 ...
  docker push $ECR_URL/ecommerce/${service}:v1.0.0
done
```

**Duration:** 15-20 minutes

---

### Step 5: Update Manifests for EKS
```bash
# Update image references to ECR
sed -i "s|image: ecommerce/|image: $ECR_URL/ecommerce/|g" k8s/services/*.yaml

# Update Ingress to use ALB
# Use k8s/ingress/ingress-alb.yaml instead of ingress.yaml
```

---

### Step 6: Deploy to EKS
```bash
cd k8s
./deploy.sh
```

**Duration:** 8-12 minutes

---

### Step 7: Configure DNS
```bash
# Get ALB DNS
ALB_DNS=$(kubectl get ingress -n ecommerce -o jsonpath='...')

# Update Route53 A record
# Point ecommerce.yourdomain.com to $ALB_DNS
```

---

### Step 8: Set Up Monitoring
```bash
# CloudWatch Container Insights
kubectl apply -f cloudwatch-agent.yaml

# Prometheus + Grafana
cd k8s/monitoring
kubectl apply -f *.yaml
```

---

### Step 9: Configure Backup
```bash
# Velero for K8s backup
velero install --provider aws ...

# Schedule daily backups
velero schedule create daily-backup --schedule="0 2 * * *"

# Database backups
# CronJob already in k8s/backup/backup-cronjobs.yaml
```

---

## Key Differences: KIND vs EKS

### Storage
- **KIND**: local-path-provisioner (hostPath)
- **EKS**: AWS EBS CSI Driver (gp3 volumes)

### Load Balancer
- **KIND**: Ingress NGINX with NodePort/HostPort
- **EKS**: AWS Application Load Balancer (ALB)

### Container Registry
- **KIND**: localhost:5000 (local registry)
- **EKS**: Amazon ECR

### Networking
- **KIND**: Localhost ports (3000, 9090, etc.)
- **EKS**: Route53 DNS + ALB + ACM TLS

### Monitoring
- **KIND**: Prometheus + Grafana in-cluster
- **EKS**: CloudWatch + Prometheus + Grafana

### Backup
- **KIND**: Local backups to disk
- **EKS**: S3 + Velero + EBS Snapshots

### Auto-Scaling
- **KIND**: Manual HPA (limited by node resources)
- **EKS**: HPA + Cluster Autoscaler + Spot Instances

### Cost
- **KIND**: $0 (runs on local machine)
- **EKS**: $200-500/month
  - EKS Control Plane: $73/month
  - EC2 Nodes: $100-300/month
  - EBS Volumes: $20-50/month
  - ALB: $20/month
  - Data Transfer: $10-50/month

---

## Troubleshooting

### KIND Cluster Creation Fails
```bash
# Check Docker
docker ps

# Delete existing cluster
kind delete cluster --name ecommerce-local

# Recreate
./kind-setup.sh
```

### Images Not Loading
```bash
# Verify images built
docker images | grep ecommerce

# Manually load image
kind load docker-image ecommerce/auth-service:latest --name ecommerce-local

# Check images in cluster
docker exec ecommerce-local-control-plane crictl images | grep ecommerce
```

### Pods Not Starting
```bash
# Check pod status
kubectl get pods -n ecommerce

# Describe pod
kubectl describe pod <pod-name> -n ecommerce

# Check logs
kubectl logs <pod-name> -n ecommerce

# Common issues:
# 1. ImagePullBackOff: Image not loaded into KIND
# 2. CrashLoopBackOff: Application error (check logs)
# 3. Pending: Insufficient resources or PVC not bound
```

### Database Connection Issues
```bash
# Check database pods
kubectl get pods -n ecommerce -l tier=database

# Test database connection
kubectl exec -it auth-db-0 -n ecommerce -- psql -U postgres -c "SELECT 1"

# Check secrets
kubectl get secret ecommerce-secrets -n ecommerce -o yaml

# Verify ConfigMap
kubectl get configmap ecommerce-config -n ecommerce -o yaml
```

### Ingress Not Working
```bash
# Check Ingress controller
kubectl get pods -n ingress-nginx

# Check Ingress resource
kubectl get ingress -n ecommerce
kubectl describe ingress ecommerce-ingress -n ecommerce

# Test service directly (bypass Ingress)
kubectl port-forward -n ecommerce svc/frontend-service 8080:80
# Access: http://localhost:8080
```

### Backup Fails
```bash
# Check pod access
kubectl get pods -n ecommerce

# Test manual backup
kubectl exec auth-db-0 -n ecommerce -- pg_dump -U postgres auth_db

# Check disk space
df -h
```

### Performance Issues
```bash
# Check resource usage
kubectl top nodes
kubectl top pods -n ecommerce

# Check HPA status
kubectl get hpa -n ecommerce

# Increase replicas manually
kubectl scale deployment auth-service -n ecommerce --replicas=5
```

---

## Best Practices

### Security
1. **Never commit secrets to git**
   ```bash
   echo "secrets/" >> .gitignore
   echo "k8s/secrets.yaml" >> .gitignore  # Update before commit
   ```

2. **Use strong passwords**
   ```bash
   # Generate secure password
   openssl rand -hex 32
   ```

3. **Rotate secrets regularly**
   ```bash
   # Update secrets
   kubectl delete secret ecommerce-secrets -n ecommerce
   kubectl create secret generic ecommerce-secrets --from-literal=...
   
   # Rolling restart
   kubectl rollout restart deployment -n ecommerce
   ```

4. **Use RBAC**
   - Already configured in `k8s/security/rbac.yaml`
   - ServiceAccounts for apps, CI/CD, backups

5. **Network Policies**
   - Already configured in `k8s/security/network-policies.yaml`
   - Zero-trust default deny + explicit allows

---

### Backup
1. **Regular backups**
   ```bash
   # Daily database backups
   crontab -e
   0 2 * * * /path/to/scripts/backup-databases.sh
   
   # Weekly full DR backup
   0 3 * * 0 /path/to/scripts/disaster-recovery.sh
   ```

2. **Test restores regularly**
   ```bash
   # Quarterly DR test
   # 1. Create test namespace
   # 2. Restore backup to test namespace
   # 3. Verify data integrity
   # 4. Delete test namespace
   ```

3. **Offsite backups**
   ```bash
   # Upload to S3 (EKS)
   aws s3 cp backups/ s3://my-backups/ --recursive
   
   # Or cloud storage (KIND)
   rclone copy backups/ remote:backups/
   ```

---

### Monitoring
1. **Set up alerts**
   - Already configured in `k8s/monitoring/prometheus-config.yaml`
   - Pod down, high CPU, high memory, disk space

2. **Dashboard access**
   ```bash
   # Grafana
   kubectl port-forward -n ecommerce svc/grafana 3001:3000
   
   # Prometheus
   kubectl port-forward -n ecommerce svc/prometheus 9090:9090
   
   # Jaeger
   kubectl port-forward -n ecommerce svc/jaeger 16686:16686
   ```

3. **Log aggregation**
   - Loki + Promtail already configured
   - View logs in Grafana

---

### Performance
1. **Resource limits**
   - All deployments have resource requests/limits
   - Adjust based on actual usage

2. **Horizontal scaling**
   - HPA configured for all services
   - Scales based on CPU (70%) and memory (80%)

3. **Database optimization**
   ```bash
   # Connection pooling in app code
   # Indexes on frequently queried columns
   # Regular VACUUM and ANALYZE
   kubectl exec auth-db-0 -n ecommerce -- psql -U postgres -c "VACUUM ANALYZE"
   ```

---

## Summary

### Scripts by Environment

**Docker Compose (Local Dev):**
- `backup.sh` - Backup databases
- `restore.sh` - Restore databases
- `setup-secrets.sh` - Generate secrets
- `health-check.sh` - Health checks

**Kubernetes (KIND + EKS):**
- `deploy.sh` - Deploy platform
- `undeploy.sh` - Teardown platform
- `backup-databases.sh` - Backup databases
- `restore-database.sh` - Restore single DB
- `disaster-recovery.sh` - Full backup

**KIND-Specific:**
- `kind-setup.sh` - Create KIND cluster
- `build-and-load-images.sh` - Build and load images

### Migration Path

1. **Local Development** (Docker Compose)
   - Use `backup.sh`, `restore.sh`, `health-check.sh`
   - Fast iteration, no K8s complexity

2. **Local Kubernetes Testing** (KIND)
   - Run `kind-setup.sh`
   - Run `build-and-load-images.sh`
   - Run `deploy.sh`
   - Test K8s features (HPA, Ingress, etc.)

3. **Production Deployment** (AWS EKS)
   - Create EKS cluster
   - Push images to ECR
   - Update manifests for EKS
   - Deploy with `deploy.sh`
   - Configure AWS services (ALB, Route53, etc.)

---

## Additional Resources

- [KIND Documentation](https://kind.sigs.k8s.io/)
- [EKS Best Practices](https://aws.github.io/aws-eks-best-practices/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [PHASE10_KUBERNETES.md](PHASE10_KUBERNETES.md) - K8s deployment guide
- [PHASE14_DISASTER_RECOVERY.md](PHASE14_DISASTER_RECOVERY.md) - DR procedures
- [eks-migration-guide.md](eks-migration-guide.md) - EKS migration guide

---

**Document Version:** 1.0  
**Last Updated:** 2024-06-12  
**Author:** GitHub Copilot  
**Status:** Production Ready ✅
