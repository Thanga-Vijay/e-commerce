# 🍎 MacBook Local Development Setup Guide

Complete guide to running the E-Commerce Platform on your MacBook using **k3d** and **Floci**.

---

## 📋 Prerequisites

### Required Software

1. **Docker Desktop for Mac**
   ```bash
   # Install via Homebrew
   brew install --cask docker
   
   # Or download from: https://www.docker.com/products/docker-desktop
   ```

2. **kubectl**
   ```bash
   brew install kubectl
   ```

3. **k3d**
   ```bash
   brew install k3d
   ```

4. **Helm** (optional but recommended)
   ```bash
   brew install helm
   ```

5. **AWS CLI** (for Floci interaction)
   ```bash
   brew install awscli
   ```

### System Requirements

- **RAM:** Minimum 16GB (32GB recommended)
- **CPU:** 4+ cores recommended
- **Disk:** 50GB free space minimum
- **macOS:** 12.0 (Monterey) or later

### Docker Desktop Configuration

Configure Docker Desktop for optimal performance:

1. Open Docker Desktop → Settings → Resources
2. Set the following:
   - **CPUs:** 6-8 cores
   - **Memory:** 12-16 GB
   - **Swap:** 2 GB
   - **Disk image size:** 100 GB

---

## 🚀 Quick Start (5 Steps)

### Step 1: Clone and Navigate

```bash
cd ~/Projects  # or your preferred directory
git clone <your-repo-url>
cd e-commerce
```

### Step 2: Make Scripts Executable

```bash
chmod +x k3d-setup.sh floci-setup.sh k8s/deploy.sh
```

### Step 3: Create k3d Cluster

```bash
./k3d-setup.sh
```

**What this does:**
- ✅ Creates k3d cluster with 1 server + 2 agent nodes
- ✅ Installs NGINX Ingress Controller
- ✅ Creates namespaces (ecommerce, monitoring, floci)
- ✅ Sets up local registry at registry.localhost:5001
- ✅ Configures /etc/hosts entries

**Time:** ~3-5 minutes

### Step 4: Deploy Floci (Cloud Services Emulator)

```bash
./floci-setup.sh
```

**What this does:**
- ✅ Deploys Floci to Kubernetes
- ✅ Provides AWS-compatible services (S3, SQS, SNS, etc.)
- ✅ Exposes on port 30456

**Time:** ~1-2 minutes

### Step 5: Deploy the E-Commerce Platform

```bash
# Deploy databases first
kubectl apply -f k8s/databases/

# Wait for databases to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n ecommerce --timeout=300s

# Deploy Redis
kubectl apply -f k8s/redis/

# Deploy Kafka
kubectl apply -f k8s/kafka/

# Wait for Kafka to be ready
kubectl wait --for=condition=ready pod -l app=kafka -n ecommerce --timeout=300s

# Deploy microservices
kubectl apply -f k8s/services/

# Deploy frontend
kubectl apply -f k8s/frontend/

# Deploy ingress
kubectl apply -f k8s/ingress/

# Deploy monitoring (optional)
kubectl apply -f k8s/monitoring/
```

**Time:** ~5-10 minutes

---

## 🏗️ K3d Cluster Architecture

### Node Configuration

```
k3d Cluster: ecommerce-cluster
├── Server Node (1)
│   └── Control Plane + Workloads
└── Agent Nodes (2)
    ├── Agent 1: Microservices (auth, product, cart, wishlist)
    └── Agent 2: Microservices (order, payment, inventory, notification, reporting)
```

### Why This Configuration?

**1 Server + 2 Agents** is optimal for MacBook because:

- ✅ **Resource Efficient:** Uses ~8-10 GB RAM total
- ✅ **High Availability:** 3 nodes provide redundancy
- ✅ **Service Distribution:** Balances 9 microservices across agents
- ✅ **Realistic:** Mimics production multi-node setup
- ✅ **Not Overkill:** More nodes would consume too much RAM

### Port Mappings

| Service | Port | Access URL |
|---------|------|------------|
| HTTP | 80 | http://ecommerce.local |
| HTTPS | 443 | https://ecommerce.local |
| PostgreSQL | 5432 | localhost:5432 |
| Kafka | 9092 | localhost:9092 |
| Redis | 6379 | localhost:6379 |
| Prometheus | 9090 | http://localhost:9090 |
| Grafana | 3001 | http://localhost:3001 |
| Floci | 30456 | http://localhost:30456 |

---

## 🎯 Resource Allocation

### Expected Resource Usage

```
Component               Pods    CPU (avg)   Memory (avg)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Microservices (9)       9       1.5 cores   3 GB
Databases (9)           9       1.0 cores   2 GB
Redis                   1       0.1 cores   256 MB
Kafka + Zookeeper       2       0.5 cores   1.5 GB
Frontend                1       0.1 cores   128 MB
NGINX Ingress           1       0.2 cores   256 MB
Monitoring Stack        3       0.5 cores   1.5 GB
Floci                   1       0.3 cores   512 MB
k3d Overhead            -       0.5 cores   1 GB
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
TOTAL                   27      4.7 cores   10.2 GB
```

### Recommended Docker Settings

For best performance on MacBook:

```
CPUs: 6-8 (leaves 2-4 for macOS)
Memory: 12-16 GB
Swap: 2 GB
```

---

## 🔧 Daily Development Workflow

### Start Everything

```bash
# 1. Ensure Docker Desktop is running

# 2. Start k3d cluster (if stopped)
k3d cluster start ecommerce-cluster

# 3. Check cluster health
kubectl get nodes
kubectl get pods -A
```

### Check Service Status

```bash
# All services
kubectl get pods -n ecommerce

# Specific service
kubectl get pods -n ecommerce -l app=auth-service

# Logs
kubectl logs -n ecommerce -l app=auth-service -f

# Describe pod (for troubleshooting)
kubectl describe pod -n ecommerce <pod-name>
```

### Access Services

```bash
# Frontend
open http://ecommerce.local

# Prometheus
open http://localhost:9090

# Grafana
open http://localhost:3001

# Floci Health Check
curl http://localhost:30456/health
```

### Stop Everything

```bash
# Stop k3d cluster (keeps data)
k3d cluster stop ecommerce-cluster

# Delete cluster (removes everything)
k3d cluster delete ecommerce-cluster
```

---

## 🛠️ Useful Commands

### Cluster Management

```bash
# List clusters
k3d cluster list

# Get cluster info
k3d cluster list ecommerce-cluster

# Import local Docker images
k3d image import <image-name> -c ecommerce-cluster

# Access cluster registry
docker tag myimage:latest registry.localhost:5001/myimage:latest
docker push registry.localhost:5001/myimage:latest
```

### Kubectl Shortcuts

```bash
# Create aliases (add to ~/.zshrc)
alias k='kubectl'
alias kge='kubectl get pods -n ecommerce'
alias kgm='kubectl get pods -n monitoring'
alias kl='kubectl logs -n ecommerce'
alias kdesc='kubectl describe -n ecommerce'

# Watch pods
watch kubectl get pods -n ecommerce

# Port forward to specific pod
kubectl port-forward -n ecommerce svc/auth-service 8081:8081
```

### Database Access

```bash
# Connect to PostgreSQL
kubectl exec -it -n ecommerce <postgres-pod> -- psql -U postgres -d auth_db

# Connect to Redis
kubectl exec -it -n ecommerce <redis-pod> -- redis-cli

# Backup database
kubectl exec -n ecommerce <postgres-pod> -- pg_dump -U postgres auth_db > backup.sql
```

### Debugging

```bash
# Get pod events
kubectl get events -n ecommerce --sort-by='.lastTimestamp'

# Check resource usage
kubectl top nodes
kubectl top pods -n ecommerce

# Restart deployment
kubectl rollout restart deployment -n ecommerce auth-service

# Scale deployment
kubectl scale deployment -n ecommerce auth-service --replicas=2

# Execute commands in pod
kubectl exec -it -n ecommerce <pod-name> -- /bin/sh
```

---

## 🧪 Testing the Setup

### Health Checks

```bash
# Check all services are running
kubectl get pods -n ecommerce -o wide

# Test auth service
curl http://localhost:8081/health

# Test frontend
curl http://ecommerce.local

# Test Floci
aws --endpoint-url=http://localhost:30456 s3 ls
```

### Create Test Data

```bash
# Register a user
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!"
  }'
```

---

## 🔍 Monitoring

### Access Monitoring Stack

1. **Prometheus:** http://localhost:9090
   - View metrics
   - Query PromQL
   - Check targets

2. **Grafana:** http://localhost:3001
   - Default credentials: admin / admin
   - Pre-configured dashboards
   - Service metrics

### Key Metrics to Watch

- **Pod CPU/Memory Usage**
- **Request Rate per Service**
- **Database Connection Pools**
- **Kafka Message Throughput**
- **Error Rates**

---

## 🐛 Troubleshooting

### Cluster Won't Start

```bash
# Check Docker is running
docker ps

# Check available resources
docker system df

# Clean up Docker
docker system prune -a

# Delete and recreate cluster
k3d cluster delete ecommerce-cluster
./k3d-setup.sh
```

### Pods in CrashLoopBackOff

```bash
# Check logs
kubectl logs -n ecommerce <pod-name> --previous

# Check events
kubectl describe pod -n ecommerce <pod-name>

# Common fixes:
# 1. Check database is ready
# 2. Check environment variables
# 3. Check resource limits
# 4. Check image pull issues
```

### Out of Memory

```bash
# Check resource usage
kubectl top nodes
kubectl top pods -A

# Increase Docker Desktop memory
# Or reduce number of services running
```

### Services Not Accessible

```bash
# Check ingress
kubectl get ingress -n ecommerce

# Check /etc/hosts
cat /etc/hosts | grep ecommerce

# Add if missing
echo "127.0.0.1 ecommerce.local api.ecommerce.local" | sudo tee -a /etc/hosts

# Restart ingress controller
kubectl rollout restart deployment -n ingress-nginx ingress-nginx-controller
```

---

## 📊 Performance Optimization

### For Better MacBook Performance

1. **Reduce Replicas** (development only):
   ```bash
   # Edit deployments to use 1 replica instead of 3
   kubectl scale deployment -n ecommerce --replicas=1 --all
   ```

2. **Disable Monitoring** (if not needed):
   ```bash
   kubectl delete namespace monitoring
   ```

3. **Use Resource Limits**:
   ```yaml
   resources:
     requests:
       cpu: 100m
       memory: 256Mi
     limits:
       cpu: 500m
       memory: 512Mi
   ```

4. **Enable BuildKit** for faster Docker builds:
   ```bash
   export DOCKER_BUILDKIT=1
   ```

---

## 🎓 Next Steps

1. **Explore the APIs:**
   - Import Postman collection from `/docs/api/`
   - Test all endpoints

2. **Develop a Service:**
   - Make changes to a microservice
   - Build Docker image
   - Deploy to k3d

3. **Add New Features:**
   - Follow the microservices pattern
   - Add Kafka events
   - Update frontend

4. **Learn Kubernetes:**
   - Experiment with scaling
   - Try rolling updates
   - Configure HPA (Horizontal Pod Autoscaler)

---

## 📚 Additional Resources

- [k3d Documentation](https://k3d.io)
- [Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/)
- [Docker Desktop for Mac](https://docs.docker.com/desktop/mac/)
- [Project Documentation](/docs/README.md)

---

## 🆘 Getting Help

- Check logs: `kubectl logs -n ecommerce -l app=<service-name>`
- Check events: `kubectl get events -n ecommerce`
- Describe resources: `kubectl describe pod/deployment/service -n ecommerce <name>`
- Project docs: `/docs/`

---

**Happy Coding! 🚀**
