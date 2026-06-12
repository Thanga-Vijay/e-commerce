# E-Commerce Platform - Phase 10: Kubernetes Deployment

This document provides comprehensive instructions for deploying the e-commerce platform to Kubernetes (Phase 10).

## Overview

Phase 10 implements production-ready Kubernetes deployment with:
- **Namespace isolation** for the entire platform
- **StatefulSets** for databases, Redis, and Kafka (data persistence)
- **Deployments** for stateless microservices
- **Services** for internal and external communication
- **Ingress** with TLS termination and path-based routing
- **HorizontalPodAutoscaler** for automatic scaling
- **ConfigMaps and Secrets** for configuration management
- **Resource limits** and health checks

## Architecture

### Kubernetes Components

```
┌─────────────────────────────────────────────┐
│           Ingress Controller                │
│  (Nginx + TLS + Rate Limiting)             │
└──────────────┬──────────────────────────────┘
               │
        ┌──────┴──────┐
        │             │
   ┌────▼─────┐  ┌───▼──────────┐
   │ Frontend │  │ API Services │
   │ (3-20)   │  │ (2-15 each)  │
   └──────────┘  └───┬──────────┘
                     │
              ┌──────┴──────┐
              │             │
         ┌────▼─────┐  ┌───▼─────┐
         │ Kafka    │  │ Databases│
         │ (3 pods) │  │ (9 pods) │
         └──────────┘  └──────────┘
```

## Prerequisites

### Required Tools
- **kubectl** 1.28+ installed and configured
- **Kubernetes cluster** (local or cloud):
  - Minikube (local development)
  - KIND (Kubernetes in Docker)
  - GKE, EKS, or AKS (production)
- **Metrics Server** (for HPA)
- **Ingress Controller** (Nginx)
- **cert-manager** (optional, for TLS)

### Cluster Requirements
- Minimum 8GB RAM
- 4 CPU cores
- 50GB storage
- LoadBalancer support (cloud) or MetalLB (on-prem)

## Quick Start

### 1. Setup Kubernetes Cluster

#### Using Minikube
```bash
# Start Minikube with sufficient resources
minikube start --cpus=4 --memory=8192 --disk-size=50g

# Enable addons
minikube addons enable ingress
minikube addons enable metrics-server
```

#### Using KIND
```bash
# Create cluster
kind create cluster --name ecommerce

# Install Nginx Ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

# Install Metrics Server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### 2. Configure Secrets

**Important:** Update secrets with your actual credentials before deploying!

```bash
# Edit secrets file
cd k8s
nano secrets.yaml

# Update these values:
# - POSTGRES_PASSWORD
# - REDIS_PASSWORD
# - JWT_SECRET
# - STRIPE_SECRET_KEY
# - STRIPE_WEBHOOK_SECRET
# - SMTP credentials
```

### 3. Deploy Platform

```bash
# Make deployment script executable
chmod +x k8s/deploy.sh

# Run deployment
./k8s/deploy.sh
```

Or deploy manually:

```bash
# 1. Create namespace
kubectl apply -f k8s/namespace.yaml

# 2. Create ConfigMaps and Secrets
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# 3. Deploy databases
kubectl apply -f k8s/databases/

# 4. Deploy Redis
kubectl apply -f k8s/redis/

# 5. Deploy Kafka
kubectl apply -f k8s/kafka/

# 6. Deploy microservices
kubectl apply -f k8s/services/

# 7. Deploy frontend
kubectl apply -f k8s/frontend/

# 8. Deploy Ingress
kubectl apply -f k8s/ingress/

# 9. Deploy HPA
kubectl apply -f k8s/hpa/
```

### 4. Verify Deployment

```bash
# Check all resources
kubectl get all -n ecommerce

# Check pods status
kubectl get pods -n ecommerce

# Check services
kubectl get svc -n ecommerce

# Check ingress
kubectl get ingress -n ecommerce
```

## Kubernetes Resources

### Namespace

**File:** `k8s/namespace.yaml`

Single namespace `ecommerce` for all resources, providing isolation and easier management.

### ConfigMaps

**File:** `k8s/configmap.yaml`

Contains non-sensitive configuration:
- Database hostnames and ports
- Redis configuration
- Kafka brokers
- Service URLs
- Application settings (log level, environment)

### Secrets

**File:** `k8s/secrets.yaml`

Contains sensitive data (base64 encoded):
- Database passwords
- Redis password
- JWT secret
- Stripe API keys
- SMTP credentials

**Security Best Practice:** Use external secret management (Sealed Secrets, External Secrets Operator, or Vault) in production.

### Databases (StatefulSets)

**Files:** `k8s/databases/*.yaml`

9 PostgreSQL databases deployed as StatefulSets:
- auth-db
- product-db
- cart-db
- wishlist-db
- order-db
- payment-db
- inventory-db
- notification-db
- reporting-db

Each database:
- Has persistent storage (5Gi per database)
- Uses headless service for stable network identity
- Has liveness and readiness probes
- Resource limits: 256Mi-512Mi RAM, 250m-500m CPU

### Redis (StatefulSet)

**File:** `k8s/redis/redis.yaml`

- Single replica with persistent storage (1Gi)
- Password-protected
- Headless service for stable identity
- Resource limits: 128Mi-256Mi RAM, 100m-250m CPU

### Kafka (StatefulSet)

**File:** `k8s/kafka/kafka.yaml`

- Zookeeper (1 replica, 2Gi + 1Gi storage)
- Kafka (3 replicas, 10Gi storage each)
- Headless services for stable network identities
- Replication factor: 3
- 3 partitions per topic by default
- Resource limits: 512Mi-1Gi RAM, 500m-1000m CPU per Kafka pod

### Microservices (Deployments)

**Files:** `k8s/services/*.yaml`

9 microservices deployed as Deployments:
- auth-service (2-10 replicas)
- product-service (3-15 replicas)
- cart-service (2-10 replicas)
- wishlist-service (2-8 replicas)
- order-service (2-12 replicas)
- payment-service (2-10 replicas)
- inventory-service (2-8 replicas)
- notification-service (2-8 replicas)
- reporting-service (2-8 replicas)

Each service:
- Has ClusterIP service for internal communication
- Uses ConfigMaps and Secrets for configuration
- Has liveness and readiness probes (`/health` endpoint)
- Resource limits: 128Mi-256Mi RAM, 100m-500m CPU
- Pulls latest image (change to specific version for production)

### Frontend (Deployment)

**File:** `k8s/frontend/frontend.yaml`

- 3-20 replicas (scales based on traffic)
- LoadBalancer service for external access
- Resource limits: 64Mi-128Mi RAM, 50m-250m CPU
- Serves static React application

### Ingress

**File:** `k8s/ingress/ingress.yaml`

Path-based routing with TLS:
- **ecommerce.yourdomain.com** → Frontend
- **api.ecommerce.yourdomain.com** → Backend services
  - `/api/v1/auth` → Auth Service
  - `/api/v1/products` → Product Service
  - `/api/v1/cart` → Cart Service
  - `/api/v1/orders` → Order Service
  - `/api/v1/payments` → Payment Service
  - And more...

Features:
- TLS termination (cert-manager integration)
- Rate limiting (100 req/s)
- Proxy timeouts configured
- Force HTTPS redirect

### Horizontal Pod Autoscaler (HPA)

**File:** `k8s/hpa/hpa.yaml`

Auto-scaling for all services based on:
- CPU utilization (70%)
- Memory utilization (80%)

Scale policies:
- Scale up aggressively (100% increase or +2 pods every 30s)
- Scale down conservatively (50% decrease every 60s, 5min stabilization)

Ranges:
- Auth: 2-10 pods
- Product: 3-15 pods
- Frontend: 3-20 pods
- Others: 2-8 pods

## Configuration

### Update Domain Names

Edit `k8s/ingress/ingress.yaml`:

```yaml
spec:
  tls:
  - hosts:
    - ecommerce.yourdomain.com    # Change this
    - api.ecommerce.yourdomain.com # Change this
```

### Update Image References

For production, use specific image versions instead of `latest`:

```yaml
# Instead of:
image: ecommerce/auth-service:latest

# Use:
image: your-registry.com/ecommerce/auth-service:v1.0.0
```

### Adjust Resource Limits

Based on your workload, adjust resource requests/limits in deployment files:

```yaml
resources:
  requests:
    memory: "256Mi"  # Minimum guaranteed
    cpu: "250m"
  limits:
    memory: "512Mi"  # Maximum allowed
    cpu: "500m"
```

## TLS/SSL Configuration

### Using cert-manager

1. **Install cert-manager:**

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

2. **Create ClusterIssuer:**

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

3. **Apply ClusterIssuer:**

```bash
kubectl apply -f cluster-issuer.yaml
```

Certificates will be automatically created when Ingress is deployed.

## Monitoring

### View Pod Status

```bash
# Watch pod status
kubectl get pods -n ecommerce -w

# Describe pod
kubectl describe pod <pod-name> -n ecommerce

# View pod logs
kubectl logs -f <pod-name> -n ecommerce

# View logs from all replicas
kubectl logs -f -l app=auth-service -n ecommerce
```

### Check HPA Status

```bash
# View HPA status
kubectl get hpa -n ecommerce

# Describe HPA
kubectl describe hpa auth-service-hpa -n ecommerce

# Watch HPA autoscaling
kubectl get hpa -n ecommerce -w
```

### Check Services

```bash
# Get services
kubectl get svc -n ecommerce

# Get ingress
kubectl get ingress -n ecommerce

# Get ingress IP/hostname
kubectl get ingress ecommerce-ingress -n ecommerce -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

### Resource Usage

```bash
# View node resources
kubectl top nodes

# View pod resources
kubectl top pods -n ecommerce

# View pod resources sorted by CPU
kubectl top pods -n ecommerce --sort-by=cpu

# View pod resources sorted by memory
kubectl top pods -n ecommerce --sort-by=memory
```

## Scaling

### Manual Scaling

```bash
# Scale deployment manually
kubectl scale deployment auth-service --replicas=5 -n ecommerce

# Scale StatefulSet
kubectl scale statefulset kafka --replicas=5 -n ecommerce
```

### Adjust HPA

```bash
# Edit HPA
kubectl edit hpa auth-service-hpa -n ecommerce

# Or update YAML and apply
kubectl apply -f k8s/hpa/hpa.yaml
```

## Database Management

### Run Migrations

```bash
# Connect to auth service pod
kubectl exec -it <auth-service-pod> -n ecommerce -- sh

# Inside pod, run migrations
migrate -path migrations -database "postgresql://user:password@auth-db-service:5432/auth_db?sslmode=disable" up
```

### Backup Databases

```bash
# Backup single database
kubectl exec auth-db-0 -n ecommerce -- pg_dump -U postgres auth_db > auth_db_backup.sql

# Backup all databases script
for db in auth product cart wishlist order payment inventory notification reporting; do
  kubectl exec ${db}-db-0 -n ecommerce -- pg_dump -U postgres ${db}_db > ${db}_db_backup.sql
done
```

### Restore Database

```bash
# Restore database
kubectl exec -i auth-db-0 -n ecommerce -- psql -U postgres auth_db < auth_db_backup.sql
```

## Troubleshooting

### Pods Not Starting

```bash
# Check pod events
kubectl describe pod <pod-name> -n ecommerce

# Common issues:
# - ImagePullBackOff: Check image name and availability
# - CrashLoopBackOff: Check logs for errors
# - Pending: Check resource availability
```

### Service Unreachable

```bash
# Check service endpoints
kubectl get endpoints <service-name> -n ecommerce

# Test service from another pod
kubectl run -it --rm debug --image=busybox --restart=Never -n ecommerce -- wget -O- http://auth-service:8081/health
```

### Database Connection Issues

```bash
# Check database pod
kubectl get pod -l app=auth-db -n ecommerce

# Check database service
kubectl get svc auth-db-service -n ecommerce

# Test connection from service pod
kubectl exec -it <auth-service-pod> -n ecommerce -- sh
# Inside pod:
nc -zv auth-db-service 5432
```

### Ingress Not Working

```bash
# Check ingress controller
kubectl get pods -n ingress-nginx

# Check ingress resource
kubectl describe ingress ecommerce-ingress -n ecommerce

# Check ingress logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
```

### HPA Not Scaling

```bash
# Check metrics server
kubectl get apiservice v1beta1.metrics.k8s.io -o yaml

# Check HPA status
kubectl describe hpa auth-service-hpa -n ecommerce

# View metrics
kubectl top pods -n ecommerce
```

## Production Best Practices

### 1. Use Specific Image Versions

```yaml
# Bad
image: ecommerce/auth-service:latest

# Good
image: ecommerce/auth-service:v1.2.3
```

### 2. Set Resource Requests and Limits

Always define both requests and limits to prevent resource starvation and ensure quality of service.

### 3. Use External Secret Management

Replace Kubernetes Secrets with:
- **Sealed Secrets**: Encrypt secrets in Git
- **External Secrets Operator**: Sync from external stores
- **HashiCorp Vault**: Enterprise secret management

### 4. Implement Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: auth-service-netpol
  namespace: ecommerce
spec:
  podSelector:
    matchLabels:
      app: auth-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8081
```

### 5. Use PodDisruptionBudgets

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: auth-service-pdb
  namespace: ecommerce
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: auth-service
```

### 6. Implement Readiness Gates

Ensure pods are truly ready before receiving traffic.

### 7. Use StatefulSets for Kafka

Kafka brokers need stable network identities and persistent storage.

### 8. Regular Backups

Automate database backups with CronJobs:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: db-backup
  namespace: ecommerce
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:15-alpine
            command: ["/bin/sh", "-c"]
            args:
            - pg_dump -h auth-db-service -U postgres auth_db > /backup/auth_db_$(date +\%Y\%m\%d).sql
```

## Cleanup

### Undeploy Everything

```bash
# Using script
chmod +x k8s/undeploy.sh
./k8s/undeploy.sh

# Or manually
kubectl delete namespace ecommerce
```

### Delete Cluster (if using KIND/Minikube)

```bash
# KIND
kind delete cluster --name ecommerce

# Minikube
minikube delete
```

## Summary

Phase 10 delivers production-ready Kubernetes deployment with:
- ✅ **25+ Kubernetes manifests** for complete platform
- ✅ **Namespace isolation** for resource organization
- ✅ **9 StatefulSet databases** with persistent storage
- ✅ **3-replica Kafka cluster** for event streaming
- ✅ **9 microservice Deployments** with auto-scaling
- ✅ **Ingress with TLS** for secure external access
- ✅ **HPA** for 10 services (automatic scaling)
- ✅ **ConfigMaps and Secrets** for configuration
- ✅ **Health checks** for all services
- ✅ **Resource limits** to prevent resource exhaustion
- ✅ **Deployment scripts** for easy deployment

Your e-commerce platform is now cloud-native and production-ready! 🚀
