# KIND Quick Start Guide

Fast-track guide to test the e-commerce platform on local Kubernetes using KIND.

## Prerequisites

Install these tools first:

### 1. KIND (Kubernetes IN Docker)
```powershell
# Windows (using chocolatey)
choco install kind

# Or download binary
curl -Lo kind.exe https://kind.sigs.k8s.io/dl/v0.20.0/kind-windows-amd64
Move-Item .\kind.exe C:\Windows\System32\kind.exe

# Verify
kind version
```

### 2. kubectl
```powershell
# Windows (using chocolatey)
choco install kubernetes-cli

# Or download
curl -LO "https://dl.k8s.io/release/v1.28.0/bin/windows/amd64/kubectl.exe"
Move-Item .\kubectl.exe C:\Windows\System32\kubectl.exe

# Verify
kubectl version --client
```

### 3. Docker Desktop
- Download from: https://www.docker.com/products/docker-desktop/
- Install and start Docker Desktop
- Verify: `docker ps`

---

## Quick Start (30 minutes)

### Step 1: Create KIND Cluster (3 minutes)
```bash
cd C:\Users\ao45j\OneDrive - Cummins\Documents\K8S\e-commerce

# Make script executable (Git Bash)
chmod +x kind-setup.sh

# Run setup
./kind-setup.sh
```

**Expected Output:**
```
✓ KIND installed
✓ kubectl installed
✓ Docker installed
Creating KIND cluster (this may take 2-3 minutes)...
✓ Ingress NGINX installed
✓ Metrics Server installed
✓ Local registry created at localhost:5000
✓ Nodes labeled
✓ Storage configured

KIND Cluster Setup Complete!
```

**Verify:**
```bash
kubectl get nodes
# Should show 4 nodes: 1 control-plane + 3 workers

kubectl get pods -A
# Should show system pods running
```

---

### Step 2: Build and Load Images (15 minutes)
```bash
cd C:\Users\ao45j\OneDrive - Cummins\Documents\K8S\e-commerce

# Make script executable
chmod +x scripts/build-and-load-images.sh

# Build and load all images
./scripts/build-and-load-images.sh
```

**Expected Output:**
```
Building auth-service...
✓ auth-service ready

Building product-service...
✓ product-service ready

...

All Images Built and Loaded Successfully!
```

**Verify:**
```bash
# Check images in KIND
docker exec -it ecommerce-local-control-plane crictl images | grep ecommerce

# Should show:
# ecommerce/auth-service:latest
# ecommerce/product-service:latest
# ... (10 images total)
```

---

### Step 3: Deploy Platform (10 minutes)
```bash
cd C:\Users\ao45j\OneDrive - Cummins\Documents\K8S\e-commerce\k8s

# Make script executable
chmod +x deploy.sh

# Deploy
./deploy.sh
```

**Expected Output:**
```
1. Creating namespace...
namespace/ecommerce created

2. Creating ConfigMaps and Secrets...
configmap/ecommerce-config created
secret/ecommerce-secrets created

3. Deploying PostgreSQL databases...
statefulset.apps/auth-db created
...

✓ Deployment complete!
```

**Verify:**
```bash
# Check all pods
kubectl get pods -n ecommerce

# Should eventually show all pods Running:
# NAME                                 READY   STATUS    RESTARTS   AGE
# auth-db-0                            1/1     Running   0          2m
# product-db-0                         1/1     Running   0          2m
# cart-db-0                            1/1     Running   0          2m
# ... (9 databases)
# redis-0                              1/1     Running   0          2m
# kafka-0                              1/1     Running   0          2m
# zookeeper-0                          1/1     Running   0          2m
# auth-service-xxx                     1/1     Running   0          1m
# product-service-xxx                  1/1     Running   0          1m
# ... (9 services)
# frontend-xxx                         1/1     Running   0          1m

# Wait for all pods to be Running
kubectl wait --for=condition=ready pod --all -n ecommerce --timeout=10m
```

---

### Step 4: Access Services

#### Frontend
```bash
# Access in browser
http://localhost:3000
```

#### API Endpoints
```bash
# Health checks
curl http://localhost/api/v1/auth/health
curl http://localhost/api/v1/products/health
curl http://localhost/api/v1/cart/health
curl http://localhost/api/v1/orders/health

# Test auth registration
curl -X POST http://localhost/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "name": "Test User"
  }'
```

#### Monitoring (Optional)
```bash
# Deploy monitoring stack
cd k8s/monitoring
kubectl apply -f prometheus-config.yaml
kubectl apply -f prometheus.yaml
kubectl apply -f grafana.yaml

# Access Grafana
http://localhost:3001
# Default: admin / admin

# Access Prometheus
http://localhost:9090
```

---

## Troubleshooting

### Issue: Pods Stuck in Pending
```bash
# Check pod events
kubectl describe pod <pod-name> -n ecommerce

# Common cause: PVC not bound
kubectl get pvc -n ecommerce

# Solution: Wait for PVC to bind, or check StorageClass
kubectl get storageclass
```

### Issue: Pods CrashLoopBackOff
```bash
# Check logs
kubectl logs <pod-name> -n ecommerce

# Common causes:
# 1. Database not ready
# 2. Missing secrets
# 3. Application error

# Check secrets
kubectl get secret ecommerce-secrets -n ecommerce -o yaml

# Check configmap
kubectl get configmap ecommerce-config -n ecommerce -o yaml
```

### Issue: Cannot Access Frontend
```bash
# Check Ingress
kubectl get ingress -n ecommerce

# Check Ingress controller
kubectl get pods -n ingress-nginx

# Port forward as workaround
kubectl port-forward -n ecommerce svc/frontend-service 8080:80
# Access: http://localhost:8080
```

### Issue: Images Not Found
```bash
# Check if images loaded
docker exec -it ecommerce-local-control-plane crictl images | grep ecommerce

# If missing, rebuild and load
./scripts/build-and-load-images.sh
```

---

## Testing Checklist

- [ ] All pods running: `kubectl get pods -n ecommerce`
- [ ] Databases accessible: `kubectl exec -it auth-db-0 -n ecommerce -- psql -U postgres -c "SELECT 1"`
- [ ] Redis accessible: `kubectl exec -it redis-0 -n ecommerce -- redis-cli ping`
- [ ] Kafka running: `kubectl get pods -n ecommerce -l app=kafka`
- [ ] Frontend accessible: http://localhost:3000
- [ ] API health checks pass: `curl http://localhost/api/v1/auth/health`
- [ ] User registration works
- [ ] User login works
- [ ] Product listing works
- [ ] Cart operations work
- [ ] Order creation works

---

## Useful Commands

### Pod Management
```bash
# List all pods
kubectl get pods -n ecommerce

# Watch pod status
kubectl get pods -n ecommerce -w

# Describe pod
kubectl describe pod <pod-name> -n ecommerce

# View logs
kubectl logs <pod-name> -n ecommerce

# Follow logs
kubectl logs -f <pod-name> -n ecommerce

# Execute command in pod
kubectl exec -it <pod-name> -n ecommerce -- bash
```

### Service Management
```bash
# List services
kubectl get svc -n ecommerce

# Port forward to service
kubectl port-forward -n ecommerce svc/auth-service 8081:8081

# Test service
curl http://localhost:8081/health
```

### Database Management
```bash
# Connect to database
kubectl exec -it auth-db-0 -n ecommerce -- psql -U postgres auth_db

# Run SQL query
kubectl exec -it auth-db-0 -n ecommerce -- psql -U postgres -c "SELECT * FROM users;"

# Backup database
kubectl exec -it auth-db-0 -n ecommerce -- pg_dump -U postgres auth_db > backup.sql

# Restore database
cat backup.sql | kubectl exec -i auth-db-0 -n ecommerce -- psql -U postgres auth_db
```

### Scaling
```bash
# Scale deployment
kubectl scale deployment auth-service -n ecommerce --replicas=5

# Check HPA
kubectl get hpa -n ecommerce

# Describe HPA
kubectl describe hpa auth-service-hpa -n ecommerce
```

### Resource Usage
```bash
# Node resources
kubectl top nodes

# Pod resources
kubectl top pods -n ecommerce

# Resource requests/limits
kubectl describe deployment auth-service -n ecommerce | grep -A5 Limits
```

---

## Backup and Restore

### Backup All Databases
```bash
./scripts/backup-databases.sh

# Creates: backups/YYYYMMDD_HHMMSS/
```

### Restore Single Database
```bash
./scripts/restore-database.sh backups/20240612_120000/auth_db_20240612_120000.sql.gz auth_db
```

### Full Disaster Recovery Backup
```bash
./scripts/disaster-recovery.sh

# Creates: disaster-recovery/YYYYMMDD_HHMMSS/
```

---

## Cleanup

### Delete Application (Keep Cluster)
```bash
cd k8s
./undeploy.sh

# This deletes all resources but keeps the KIND cluster
```

### Delete KIND Cluster
```bash
kind delete cluster --name ecommerce-local

# This deletes the entire cluster
```

### Delete Local Images
```bash
# Delete ecommerce images
docker images | grep ecommerce | awk '{print $3}' | xargs docker rmi -f

# Prune unused images
docker image prune -a
```

---

## Next Steps

### 1. Deploy Monitoring
```bash
cd k8s/monitoring
kubectl apply -f prometheus-config.yaml
kubectl apply -f prometheus.yaml
kubectl apply -f grafana.yaml
kubectl apply -f loki.yaml
kubectl apply -f jaeger.yaml

# Access Grafana: http://localhost:3001
# Access Prometheus: http://localhost:9090
# Access Jaeger: http://localhost:16686
```

### 2. Set Up Continuous Integration
- See `.github/workflows/ci.yml`
- Set up GitHub Actions
- Add Docker Hub credentials
- Run tests on every push

### 3. Prepare for AWS EKS
- Review [eks-migration-guide.md](eks-migration-guide.md)
- Create AWS account
- Install AWS CLI
- Install eksctl
- Estimate costs

### 4. Security Hardening
```bash
# Apply security policies
kubectl apply -f k8s/security/network-policies.yaml
kubectl apply -f k8s/security/rbac.yaml
kubectl apply -f k8s/security/pod-security.yaml

# Verify
kubectl get networkpolicies -n ecommerce
kubectl get roles,rolebindings -n ecommerce
kubectl get pdb -n ecommerce
```

### 5. Performance Testing
```bash
# Install K6
choco install k6

# Run load tests
k6 run loadtests/auth-load-test.js
k6 run loadtests/product-load-test.js
k6 run loadtests/order-load-test.js

# Monitor HPA
kubectl get hpa -n ecommerce -w
```

---

## FAQ

### Q: Can I use Minikube instead of KIND?
**A:** Yes, but:
- KIND is faster and lighter
- KIND uses real Kubernetes (not VM-based)
- KIND supports multi-node clusters easily
- Scripts are optimized for KIND

### Q: How much disk space does this need?
**A:** Approximately:
- Docker images: 5-10 GB
- KIND cluster: 2-3 GB
- Database storage: 1-2 GB
- Total: ~10-15 GB

### Q: Can I run this on Windows without WSL?
**A:** Yes, using:
- Docker Desktop for Windows
- Git Bash (for shell scripts)
- Or PowerShell (convert scripts)

### Q: How do I update a single service?
```bash
# Rebuild service
docker build -t ecommerce/auth-service:latest services/auth-service/

# Load into KIND
kind load docker-image ecommerce/auth-service:latest --name ecommerce-local

# Rolling restart
kubectl rollout restart deployment auth-service -n ecommerce

# Watch rollout
kubectl rollout status deployment auth-service -n ecommerce
```

### Q: How do I add more worker nodes?
Edit `kind-config.yaml` and add:
```yaml
  - role: worker
    image: kindest/node:v1.28.0
```
Then recreate cluster: `./kind-setup.sh`

### Q: How do I check resource limits?
```bash
# Show resource requests/limits
kubectl describe deployment auth-service -n ecommerce | grep -A10 Limits

# Check actual usage
kubectl top pod auth-service-xxx -n ecommerce

# Check HPA status
kubectl describe hpa auth-service-hpa -n ecommerce
```

---

## Performance Tips

1. **Enable BuildKit**
   ```bash
   export DOCKER_BUILDKIT=1
   export COMPOSE_DOCKER_CLI_BUILD=1
   ```

2. **Use Docker Layer Caching**
   - Already configured in Dockerfiles
   - Rebuild only changed services

3. **Allocate More Resources to Docker Desktop**
   - Settings → Resources
   - Increase CPUs: 4-6
   - Increase Memory: 8-12 GB
   - Increase Disk: 50+ GB

4. **Reduce Replica Counts for Testing**
   ```bash
   # Edit k8s/services/*.yaml
   # Change replicas: 3 to replicas: 1
   ```

5. **Skip Building Unchanged Services**
   ```bash
   # Only build changed services
   docker build -t ecommerce/auth-service:latest services/auth-service/
   kind load docker-image ecommerce/auth-service:latest --name ecommerce-local
   kubectl rollout restart deployment auth-service -n ecommerce
   ```

---

## Support

- Issues: GitHub Issues
- Docs: `README.md`, `DEPLOYMENT_GUIDE.md`
- Scripts: `SCRIPTS_ANALYSIS.md`
- EKS Migration: `eks-migration-guide.md`

---

**Ready to start?**
```bash
./kind-setup.sh
```

🚀 Good luck!
