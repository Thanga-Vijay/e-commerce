# Docker Compose Cleanup Guide

This guide explains which Docker Compose files are needed for local MacBook development with k3d.

## ✅ Files to KEEP

### 1. `docker-compose.yml` (KEEP - Modified for local dev)
**Purpose:** Local development without Kubernetes  
**Use Case:** Quick testing of individual services  
**When to use:** 
- Developing a single microservice
- Testing service integration locally
- Running databases locally without k3d

**Note:** For MacBook k3d setup, you'll primarily use Kubernetes manifests in `/k8s` folder.

### 2. `docker-compose.override.yml` (KEEP)
**Purpose:** Local environment overrides  
**Use Case:** Developer-specific configurations  
**Automatically merged with docker-compose.yml**

## ❌ Files to REMOVE (For k3d-based development)

### 1. `docker-compose.prod.yml` ❌
**Reason:** Production configuration not needed for local MacBook development  
**Alternative:** Use Kubernetes manifests in `/k8s` for production-like environment

### 2. `docker-compose.kafka.yml` ❌
**Reason:** Kafka is deployed via Kubernetes (`/k8s/kafka/`)  
**Alternative:** `kubectl apply -f k8s/kafka/`

### 3. `docker-compose.monitoring.yml` ❌
**Reason:** Monitoring stack deployed via Kubernetes (`/k8s/monitoring/`)  
**Alternative:** `kubectl apply -f k8s/monitoring/`

### 4. `docker-compose.secrets.yml` ❌
**Reason:** Secrets managed by Kubernetes Secrets  
**Alternative:** `kubectl create secret` or use `/k8s/secrets.yaml`

## 📝 Recommendation

For **k3d-based MacBook development**, use this approach:

### Primary Method: Kubernetes (Recommended)
```bash
# Start k3d cluster
./k3d-setup.sh

# Deploy everything via Kubernetes
kubectl apply -f k8s/databases/
kubectl apply -f k8s/kafka/
kubectl apply -f k8s/services/
```

**Pros:**
- ✅ Production-like environment
- ✅ Service mesh capabilities
- ✅ Scaling and load balancing
- ✅ Better isolation
- ✅ Learn Kubernetes

### Fallback: Docker Compose (For quick tests)
```bash
# Only when you need quick local testing without k3d
docker-compose up -d auth-service auth-db redis
```

**Pros:**
- ✅ Faster startup for single services
- ✅ Less resource usage
- ✅ Simpler for isolated testing

## 🗂️ Recommended File Structure

```
e-commerce/
├── docker-compose.yml              # KEEP - Local dev (simplified)
├── docker-compose.override.yml     # KEEP - Local overrides
├── docker-compose.prod.yml         # REMOVE or move to /archive
├── docker-compose.kafka.yml        # REMOVE (use k8s/kafka/)
├── docker-compose.monitoring.yml   # REMOVE (use k8s/monitoring/)
├── docker-compose.secrets.yml      # REMOVE (use k8s/secrets.yaml)
└── k8s/                            # PRIMARY - Use these!
    ├── databases/
    ├── kafka/
    ├── monitoring/
    ├── services/
    └── ...
```

## 🚀 Migration Commands

If you want to clean up:

```bash
# Create archive directory
mkdir -p archive/docker-compose

# Move unused files
mv docker-compose.prod.yml archive/docker-compose/
mv docker-compose.kafka.yml archive/docker-compose/
mv docker-compose.monitoring.yml archive/docker-compose/
mv docker-compose.secrets.yml archive/docker-compose/

# Keep only essential files
# - docker-compose.yml (for fallback local testing)
# - docker-compose.override.yml
```

## 💡 When to Use What

| Scenario | Use This |
|----------|----------|
| Full e-commerce platform testing | **k3d + Kubernetes** |
| Develop single microservice | **Docker Compose** (single service) |
| Test service integration | **k3d + Kubernetes** |
| Quick database access | **Docker Compose** (db only) |
| Production-like environment | **k3d + Kubernetes** |
| CI/CD testing | **Kubernetes** |
| Resource-constrained laptop | **Docker Compose** (minimal services) |

## 🎯 Final Recommendation

For your MacBook k3d setup:

1. **DELETE or ARCHIVE:**
   - `docker-compose.prod.yml`
   - `docker-compose.kafka.yml`
   - `docker-compose.monitoring.yml`
   - `docker-compose.secrets.yml`

2. **KEEP (but use rarely):**
   - `docker-compose.yml` - Simplify to only core services for fallback testing
   - `docker-compose.override.yml` - For personal dev preferences

3. **USE PRIMARILY:**
   - `/k8s/*` manifests via k3d cluster
   - Follow the MacBook setup guide in `MACBOOK_SETUP.md`

This keeps your development environment clean and focused on the k3d + Kubernetes workflow!
