# 🚀 Quick Start Guide - k3d + LocalStack Setup

Get the E-Commerce Platform running in **5 minutes** with k3d and LocalStack!

---

## 📋 Prerequisites Check

Run these commands to verify prerequisites:

```bash
# Check Docker
docker --version
# Should show: Docker version 24.0+ or higher

# Check kubectl
kubectl version --client
# Should show: v1.28+ or higher

# Check k3d
k3d version
# Should show: v5.0+ or higher

# Check if Docker is running
docker ps
# Should show list of containers (or empty list if none running)
```

**Missing something?**
- Docker: https://docs.docker.com/get-docker/
- kubectl: https://kubernetes.io/docs/tasks/tools/
- k3d: https://k3d.io/v5.6.0/#installation

---

## 🎯 Step-by-Step Setup

### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/e-commerce.git
cd e-commerce
```

### Step 2: Make Scripts Executable

**Linux/macOS:**
```bash
chmod +x k3d-setup.sh localstack-setup.sh scripts/*.sh
```

**Windows (Git Bash or WSL):**
```bash
chmod +x k3d-setup.sh localstack-setup.sh scripts/*.sh
```

### Step 3: Create k3d Cluster

```bash
./k3d-setup.sh
```

**What happens:**
- ✅ Creates k3d cluster with 1 server + 2 agents
- ✅ Installs NGINX Ingress Controller
- ✅ Creates namespaces (ecommerce, monitoring, localstack)
- ✅ Sets up local registry at registry.localhost:5001
- ✅ Configures /etc/hosts entries

**Expected output:**
```
🚀 Setting up k3d cluster for E-Commerce Platform
📦 Creating k3d cluster...
⏳ Waiting for cluster to be ready...
✅ k3d cluster setup complete!
```

**Time:** ~2-3 minutes

### Step 4: Deploy LocalStack

```bash
./localstack-setup.sh
```

**What happens:**
- ✅ Deploys LocalStack to Kubernetes
- ✅ Waits for LocalStack to be ready
- ✅ Exposes LocalStack on port 30456

**Expected output:**
```
🚀 Setting up LocalStack for AWS Services Emulation
📦 Deploying LocalStack...
⏳ Waiting for LocalStack to be ready...
✅ LocalStack deployed successfully!
```

**Time:** ~1 minute

### Step 5: Initialize AWS Resources

```bash
./scripts/init-localstack.sh
```

**What happens:**
- ✅ Creates S3 buckets (images, backups, logs)
- ✅ Creates ECR repositories for all services
- ✅ Creates Secrets Manager secrets
- ✅ Creates SSM parameters
- ✅ Creates CloudWatch log groups

**Expected output:**
```
🔧 Initializing LocalStack AWS Resources
✅ LocalStack is ready!
📦 Creating S3 buckets...
🐳 Creating ECR repositories...
✅ LocalStack initialization complete!
```

**Time:** ~30 seconds

### Step 6: Build Docker Images

```bash
make build-all
```

**What happens:**
- ✅ Builds Docker images for 9 Go services
- ✅ Builds React frontend image

**Time:** ~5-10 minutes (first time, then cached)

### Step 7: Push Images to Local Registry

```bash
make push-all
```

**What happens:**
- ✅ Tags images with registry.localhost:5001
- ✅ Pushes all images to local registry

**Time:** ~1-2 minutes

### Step 8: Deploy Infrastructure

```bash
make deploy-infra
```

**What happens:**
- ✅ Deploys PostgreSQL (9 databases)
- ✅ Deploys Redis cache
- ✅ Deploys Kafka message broker
- ✅ Applies ConfigMap and Secrets

**Expected output:**
```
🏗️  Deploying infrastructure...
⏳ Waiting for infrastructure to be ready...
pod/postgres-0 condition met
pod/kafka-0 condition met
```

**Time:** ~2-3 minutes

### Step 9: Deploy Microservices

```bash
make deploy-services
```

**What happens:**
- ✅ Deploys 9 Go microservices
- ✅ Deploys React frontend
- ✅ Configures Ingress routes

**Time:** ~2 minutes

### Step 10: Verify Deployment

```bash
make status
```

**Expected output:**
```
📊 Deployment Status:

🟢 Pods:
NAME                              READY   STATUS    RESTARTS   AGE
auth-service-xxx                  1/1     Running   0          2m
product-service-xxx               1/1     Running   0          2m
cart-service-xxx                  1/1     Running   0          2m
...

🌐 Services:
NAME               TYPE        CLUSTER-IP      PORT(S)    AGE
auth-service       ClusterIP   10.43.x.x       8081/TCP   2m
product-service    ClusterIP   10.43.x.x       8082/TCP   2m
...

🚪 Ingress:
NAME              HOSTS               ADDRESS          PORTS     AGE
ecommerce         ecommerce.local     192.168.x.x      80        2m
```

### Step 11: Access the Application

Open your browser:

- **Frontend:** http://ecommerce.local
- **API Health:** http://api.ecommerce.local/auth/health
- **LocalStack:** http://localhost:30456/_localstack/health

---

## 🎉 Success!

You now have:
- ✅ 9 microservices running
- ✅ React frontend running
- ✅ PostgreSQL with 9 databases
- ✅ Redis cache
- ✅ Kafka event streaming
- ✅ LocalStack (AWS emulation)
- ✅ NGINX Ingress

---

## 🔍 Verification Commands

```bash
# Check all pods are running
kubectl get pods -n ecommerce

# Check services
kubectl get svc -n ecommerce

# View logs of a service
kubectl logs -f deployment/auth-service -n ecommerce

# Test API endpoints
curl http://api.ecommerce.local/auth/health
curl http://api.ecommerce.local/product/health

# Check LocalStack
curl http://localhost:30456/_localstack/health

# List S3 buckets
aws --endpoint-url=http://localhost:30456 s3 ls
```

---

## 🛠️ One-Command Setup

If you want to do everything in one command:

```bash
make setup
```

This runs all steps automatically:
1. Creates k3d cluster
2. Deploys LocalStack
3. Initializes AWS resources
4. Builds all images
5. Pushes to registry
6. Deploys infrastructure
7. Deploys services

**Total time:** ~15-20 minutes

---

## 🧹 Cleanup

When you're done:

```bash
# Delete all Kubernetes resources
make clean

# Delete k3d cluster
make k3d-down

# Or delete everything
k3d cluster delete ecommerce-cluster
```

---

## 🔧 Troubleshooting

### Pods not starting?

```bash
# Check pod status
kubectl get pods -n ecommerce

# Describe problematic pod
kubectl describe pod <pod-name> -n ecommerce

# Check logs
kubectl logs <pod-name> -n ecommerce
```

### Can't access application?

```bash
# Check /etc/hosts has entries
cat /etc/hosts | grep ecommerce

# Should show:
# 127.0.0.1 ecommerce.local api.ecommerce.local

# Check ingress
kubectl get ingress -n ecommerce
```

### LocalStack not working?

```bash
# Check LocalStack pod
kubectl get pods -n localstack

# Port-forward if needed
kubectl port-forward -n localstack svc/localstack 4566:4566

# Test health
curl http://localhost:4566/_localstack/health
```

### Images not pulling?

```bash
# Check registry
docker ps | grep registry

# Should show registry container

# Re-tag and push
docker tag ecommerce/auth-service:latest registry.localhost:5001/ecommerce/auth-service:latest
docker push registry.localhost:5001/ecommerce/auth-service:latest
```

---

## 📚 Next Steps

- **Explore API:** See API documentation in README.md
- **Run Tests:** `make test`
- **View Logs:** `make logs`
- **Monitor:** Deploy Prometheus/Grafana
- **Develop:** Build custom features

---

## 💡 Useful Make Commands

```bash
make help              # Show all available commands
make status            # Check deployment status
make logs              # Stream logs from all pods
make health            # Run health checks
make test              # Run all tests
make clean             # Clean up everything
make k3d-restart       # Restart k3d cluster
make forward-postgres  # Port-forward PostgreSQL
make forward-redis     # Port-forward Redis
```

---

## 🎓 Learning Resources

- **k3d Docs:** https://k3d.io/
- **LocalStack Docs:** https://docs.localstack.cloud/
- **Kubernetes Basics:** https://kubernetes.io/docs/tutorials/
- **Docker Basics:** https://docs.docker.com/get-started/

---

**Need help?** Open an issue on GitHub or check the main README.md!
