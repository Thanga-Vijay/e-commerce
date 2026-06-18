# KIND Testing - Quick Action Checklist

## ✅ Changes Analyzed and Fixed

All your changes for KIND testing have been analyzed. Here's what you did and what's been fixed:

---

## What You Changed (All Good!)

### 1. **ConfigMap** (`k8s/configmap.yaml`)
- ✅ Renamed to `app-config`
- ✅ Single Kafka broker for KIND
- ✅ Deployments already updated to use new name

### 2. **Secrets** (`k8s/secrets.yaml`)
- ✅ Renamed to `app-secrets`
- ✅ Using `stringData` (plain text, easier)
- ✅ Simple "changeme" passwords for testing
- ✅ Deployments already updated to use new name

### 3. **Kafka** (`k8s/kafka/kafka.yaml`)
- ✅ Reduced to 1 replica (from 3)
- ✅ Replication factor 1
- ✅ **Saves 2GB RAM**

### 4. **Ingress** (`k8s/ingress/ingress.yaml`)
- ✅ Using `.local` domains
- ✅ Disabled SSL redirects
- ✅ Ready for local testing

### 5. **KIND Config** (`kind-config.yaml`)
- ✅ Simplified configuration
- ✅ Updated to Kubernetes v1.30.8
- ✅ Kept essential port mappings

### 6. **Build Script** (`scripts/build-and-load-images.sh`)
- ✅ Correct paths (services in root directory)
- ✅ **Fixed: Registry port 5000 → 5001** ✨

### 7. **KIND Setup** (`kind-setup.sh`)
- ✅ Registry on port 5001
- ✅ **Fixed: Output messages now show correct port** ✨

---

## Required Actions Before Testing

### 1. Add Hosts File Entries (Required)

**Windows** (Run PowerShell as Administrator):
```powershell
Add-Content C:\Windows\System32\drivers\etc\hosts "127.0.0.1 ecommerce.local"
Add-Content C:\Windows\System32\drivers\etc\hosts "127.0.0.1 api.ecommerce.local"
```

**Or manually edit**: `C:\Windows\System32\drivers\etc\hosts`
```
127.0.0.1 ecommerce.local
127.0.0.1 api.ecommerce.local
```

### 2. Verify Docker Settings

Ensure Docker Desktop has enough resources:
- **CPUs**: 4-6 cores
- **Memory**: 8-12 GB
- **Disk**: 50+ GB

**To check/update**: Docker Desktop → Settings → Resources

---

## Ready to Test! 🚀

### Step 1: Create KIND Cluster (3 min)
```bash
cd "C:\Users\ao45j\OneDrive - Cummins\Documents\K8S\e-commerce"
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
✓ Local registry created at localhost:5001
✓ Nodes labeled
✓ Storage configured

KIND Cluster Setup Complete!
```

### Step 2: Build and Load Images (15 min)
```bash
./scripts/build-and-load-images.sh
```

**Expected Output:**
```
Building auth-service...
✓ auth-service ready

Building product-service...
✓ product-service ready

... (9 services + frontend)

All Images Built and Loaded Successfully!
```

### Step 3: Deploy Platform (10 min)
```bash
cd k8s
./deploy.sh
```

**Expected Output:**
```
1. Creating namespace...
2. Creating ConfigMaps and Secrets...
3. Deploying PostgreSQL databases...
4. Deploying Redis...
5. Deploying Kafka and Zookeeper...
6. Deploying microservices...
7. Deploying frontend...
8. Deploying Ingress...
9. Deploying Horizontal Pod Autoscalers...

✓ Deployment complete!
```

### Step 4: Verify Deployment
```bash
# Check all pods are running
kubectl get pods -n ecommerce

# Wait for all pods to be ready (may take 5-10 minutes)
kubectl wait --for=condition=ready pod --all -n ecommerce --timeout=10m

# Check ingress
kubectl get ingress -n ecommerce

# Check services
kubectl get svc -n ecommerce
```

### Step 5: Access Services

**Frontend:**
- URL: http://ecommerce.local
- Should show React frontend

**API Health Checks:**
```bash
curl http://api.ecommerce.local/api/v1/auth/health
curl http://api.ecommerce.local/api/v1/products/health
curl http://api.ecommerce.local/api/v1/cart/health
```

**Expected Response:**
```json
{"status":"ok","service":"auth-service","timestamp":"2026-06-17T..."}
```

---

## Troubleshooting

### Issue: "Cannot resolve ecommerce.local"
**Solution**: Add hosts file entries (see above)

### Issue: Pods stuck in "Pending"
```bash
# Check PVCs
kubectl get pvc -n ecommerce

# Describe pod to see issue
kubectl describe pod <pod-name> -n ecommerce
```
**Solution**: Usually resolves automatically as PVCs bind (wait 2-3 minutes)

### Issue: Pods in "CrashLoopBackOff"
```bash
# Check logs
kubectl logs <pod-name> -n ecommerce
```
**Common causes**:
- Database not ready (wait for all DB pods to be Running)
- Wrong secret values (check app-secrets)

### Issue: "ImagePullBackOff"
```bash
# Verify images loaded
docker exec -it ecommerce-local-control-plane crictl images | grep ecommerce
```
**Solution**: If missing, rebuild:
```bash
./scripts/build-and-load-images.sh
```

### Issue: Cannot access via browser
```bash
# Check ingress controller
kubectl get pods -n ingress-nginx

# Should see controller pod Running
```
**Workaround**: Port forward directly
```bash
kubectl port-forward -n ecommerce svc/frontend-service 8080:80
# Access: http://localhost:8080
```

---

## Quick Reference

### Useful Commands

```bash
# View all resources
kubectl get all -n ecommerce

# Watch pods
kubectl get pods -n ecommerce -w

# View logs
kubectl logs -f <pod-name> -n ecommerce

# Restart deployment
kubectl rollout restart deployment <service-name> -n ecommerce

# Delete and redeploy
cd k8s
./undeploy.sh
./deploy.sh

# Connect to database
kubectl exec -it auth-db-0 -n ecommerce -- psql -U postgres auth_db

# Connect to Redis
kubectl exec -it redis-0 -n ecommerce -- redis-cli

# Port forward to service
kubectl port-forward -n ecommerce svc/auth-service 8081:8081
```

### Resource Usage

```bash
# Node resources
kubectl top nodes

# Pod resources
kubectl top pods -n ecommerce

# Check HPA status
kubectl get hpa -n ecommerce
```

### Cleanup

```bash
# Delete application (keep cluster)
cd k8s
./undeploy.sh

# Delete entire cluster
kind delete cluster --name ecommerce-local

# Delete local images
docker images | grep ecommerce | awk '{print $3}' | xargs docker rmi -f
```

---

## Resource Requirements

### Minimum:
- **CPU**: 4 cores
- **RAM**: 8 GB
- **Disk**: 30 GB

### Recommended:
- **CPU**: 6 cores
- **RAM**: 12 GB
- **Disk**: 50 GB

### Expected Usage:
- **KIND cluster**: ~2 GB RAM
- **Databases (9)**: ~4.5 GB RAM
- **Redis**: ~256 MB RAM
- **Kafka**: ~1 GB RAM
- **Services (9)**: ~2-3 GB RAM
- **Frontend**: ~256 MB RAM
- **Total**: ~10-11 GB RAM

---

## Next Steps After Testing

1. ✅ **Test all features**:
   - User registration/login
   - Product browsing
   - Cart operations
   - Order creation
   - Payment processing

2. ✅ **Monitor performance**:
   - Check pod resource usage
   - Verify HPA scaling
   - Monitor logs

3. ✅ **Deploy monitoring** (optional):
   ```bash
   cd k8s/monitoring
   kubectl apply -f prometheus-config.yaml
   kubectl apply -f prometheus.yaml
   kubectl apply -f grafana.yaml
   ```
   - Access Grafana: http://localhost:3001

4. ✅ **Prepare for EKS**:
   - Review [eks-migration-guide.md](eks-migration-guide.md)
   - Test backup/restore scripts
   - Document any issues found

---

## Summary

✅ **All configurations optimized for KIND**  
✅ **Registry port fixed**  
✅ **Deployments match new ConfigMap/Secret names**  
✅ **Ready to test**  

**Estimated total time**: ~30 minutes  
**Resource savings vs production**: ~2GB RAM (Kafka optimization)  

---

## Support Files

- **Analysis**: [KIND_CHANGES_ANALYSIS.md](KIND_CHANGES_ANALYSIS.md)
- **Quick Start**: [KIND_QUICKSTART.md](KIND_QUICKSTART.md)
- **Scripts Analysis**: [SCRIPTS_ANALYSIS.md](SCRIPTS_ANALYSIS.md)
- **EKS Migration**: [eks-migration-guide.md](eks-migration-guide.md)

---

**Ready to start?** Run: `./kind-setup.sh` 🚀
