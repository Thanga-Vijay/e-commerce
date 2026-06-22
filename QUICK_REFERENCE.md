# 🚀 Quick Reference Card - MacBook k3d Setup

## Essential Commands

### Start/Stop Cluster
```bash
k3d cluster start ecommerce-cluster
k3d cluster stop ecommerce-cluster
k3d cluster delete ecommerce-cluster
```

### Check Status
```bash
kubectl get nodes                    # Check cluster nodes
kubectl get pods -n ecommerce        # Check application pods
kubectl get pods -A                  # Check all pods
kubectl top nodes                    # Resource usage
```

### View Logs
```bash
kubectl logs -n ecommerce -l app=auth-service -f
kubectl logs -n ecommerce <pod-name> -f
kubectl logs -n ecommerce <pod-name> --previous
```

### Port Forward (Access Services)
```bash
kubectl port-forward -n ecommerce svc/auth-service 8081:8081
kubectl port-forward -n ecommerce svc/product-service 8082:8082
```

### Scale Services
```bash
kubectl scale deployment -n ecommerce auth-service --replicas=2
kubectl scale deployment -n ecommerce --replicas=1 --all
```

### Restart Services
```bash
kubectl rollout restart deployment -n ecommerce auth-service
kubectl rollout restart deployment -n ecommerce --all
```

### Debug Pods
```bash
kubectl describe pod -n ecommerce <pod-name>
kubectl exec -it -n ecommerce <pod-name> -- /bin/sh
kubectl get events -n ecommerce --sort-by='.lastTimestamp'
```

### Database Access
```bash
# PostgreSQL
kubectl exec -it -n ecommerce <postgres-pod> -- psql -U postgres -d auth_db

# Redis
kubectl exec -it -n ecommerce <redis-pod> -- redis-cli
```

## Access URLs

| Service | URL |
|---------|-----|
| Frontend | http://ecommerce.local |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3001 |
| Floci | http://localhost:30456 |

## Cluster Info

- **Configuration:** 1 Server + 2 Agents
- **Registry:** registry.localhost:5001
- **Namespaces:** ecommerce, monitoring, floci
- **Ingress:** NGINX Ingress Controller

## Resource Usage

- **RAM:** ~10 GB
- **CPU:** ~5-6 cores
- **Pods:** ~27 total

## One-Line Setup

```bash
./setup-all.sh
```

## Troubleshooting

### Pods Not Starting
```bash
kubectl describe pod -n ecommerce <pod-name>
kubectl get events -n ecommerce
```

### Out of Resources
```bash
# Increase Docker Desktop memory
# Or scale down services
kubectl scale deployment -n monitoring --replicas=0 --all
```

### Reset Everything
```bash
k3d cluster delete ecommerce-cluster
./setup-all.sh
```

## File Structure

```
e-commerce/
├── k3d-config.yaml         # Cluster configuration
├── k3d-setup.sh            # Create cluster
├── floci-setup.sh          # Deploy Floci
├── setup-all.sh            # Full automated setup
├── MACBOOK_SETUP.md        # Detailed guide
├── K3D_CLUSTER_GUIDE.md    # Cluster details
└── k8s/                    # Kubernetes manifests
    ├── databases/
    ├── services/
    ├── kafka/
    ├── monitoring/
    └── ...
```

## Daily Workflow

1. **Start work:**
   ```bash
   k3d cluster start ecommerce-cluster
   kubectl get pods -n ecommerce
   ```

2. **Make changes:**
   - Edit code
   - Build Docker image
   - Push to local registry
   - Update deployment

3. **Test:**
   ```bash
   kubectl logs -n ecommerce -l app=<service> -f
   curl http://ecommerce.local
   ```

4. **Stop work:**
   ```bash
   k3d cluster stop ecommerce-cluster
   ```

## Helpful Aliases

Add to `~/.zshrc`:

```bash
alias k='kubectl'
alias kge='kubectl get pods -n ecommerce'
alias kgm='kubectl get pods -n monitoring'
alias kl='kubectl logs -n ecommerce'
alias kdesc='kubectl describe -n ecommerce'
alias kx='kubectl exec -it -n ecommerce'
```

---

**📖 Full Docs:** MACBOOK_SETUP.md  
**🆘 Support:** Check logs and events first
