# K3d Cluster Configuration Reference

This document explains the k3d cluster configuration for MacBook local development.

## Cluster Sizing Decision

### Configuration: 1 Server + 2 Agents

```yaml
servers: 1
agents: 2
```

### Why This Configuration?

#### ✅ Advantages

1. **Resource Efficient**
   - Total: 3 nodes (manageable on MacBook)
   - RAM usage: ~8-10 GB total
   - CPU usage: ~4-6 cores
   - Leaves resources for macOS and other apps

2. **Service Distribution**
   ```
   Agent 1: 
   - auth-service
   - product-service
   - cart-service
   - wishlist-service
   
   Agent 2:
   - order-service
   - payment-service
   - inventory-service
   - notification-service
   - reporting-service
   ```

3. **High Availability**
   - Simulates multi-node production cluster
   - Provides pod scheduling flexibility
   - Enables node failure testing

4. **Learning Benefits**
   - Experience real pod scheduling
   - Test node affinity/anti-affinity
   - Understand distributed systems

#### ❌ Alternative Configurations Considered

| Config | Pros | Cons |
|--------|------|------|
| **1 Server + 0 Agents** | Minimal resources | No HA, not production-like |
| **1 Server + 1 Agent** | Low resource usage | Limited pod distribution |
| **1 Server + 2 Agents** | ✅ **CHOSEN** | Balanced |
| **1 Server + 3 Agents** | More distribution | Higher RAM (~12GB) |
| **1 Server + 4+ Agents** | Maximum HA | Too heavy for laptop |

## Resource Allocation

### MacBook Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| RAM | 16 GB | 32 GB |
| CPU Cores | 4 | 8 |
| Storage | 50 GB | 100 GB |

### Cluster Resource Usage

```
Component              Nodes    Memory    CPU
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
k3d Control Plane      1        1.5 GB    0.5
k3d Agent 1            1        3.0 GB    2.0
k3d Agent 2            1        3.0 GB    2.0
System Overhead        -        1.5 GB    0.5
Docker Desktop         -        1.0 GB    0.5
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
TOTAL                  3        10 GB     5.5 cores
```

## Port Mappings Explained

### Application Ports

```yaml
ports:
  - port: 80:80          # HTTP traffic to Ingress
  - port: 443:443        # HTTPS traffic to Ingress
  - port: 5432:5432      # PostgreSQL external access
  - port: 9092:9092      # Kafka external access
  - port: 6379:6379      # Redis external access
  - port: 9090:9090      # Prometheus UI
  - port: 3001:3000      # Grafana UI (mapped to 3001)
```

### Why These Ports?

- **80/443:** Standard web traffic through NGINX Ingress
- **5432:** Direct database access for tools (pgAdmin, DBeaver)
- **9092:** Kafka for external producers/consumers
- **6379:** Redis CLI access from host
- **9090:** Prometheus metrics browsing
- **3001:** Grafana dashboards (3000 often in use)

## Volume Mounts

```yaml
volumes:
  - volume: ./data:/data              # Data persistence
  - volume: ./k8s:/k8s               # Easy manifest access
```

### Benefits

- **Persistence:** Data survives cluster restarts
- **Easy Access:** Edit manifests from host
- **Backups:** Simple to backup ./data directory

## Registry Configuration

```yaml
registries:
  create:
    name: registry.localhost
    host: "0.0.0.0"
    hostPort: "5001"
```

### Usage

```bash
# Tag image
docker tag myservice:latest registry.localhost:5001/myservice:latest

# Push to local registry
docker push registry.localhost:5001/myservice:latest

# Use in Kubernetes
image: registry.localhost:5001/myservice:latest
```

### Benefits

- ⚡ Fast image loading (no remote pull)
- 🔒 Private images
- 💾 Reduces bandwidth usage

## K3s Optimizations

### Disabled Components

```yaml
- arg: --disable=traefik           # Using NGINX Ingress instead
- arg: --disable=metrics-server    # Using Prometheus instead
```

### Why Disable?

- **Traefik:** NGINX Ingress is more common in production
- **Metrics Server:** Prometheus provides richer metrics

### Enabled Optimizations

```yaml
- arg: --kube-apiserver-arg=--max-requests-inflight=400
```

Increases API server throughput for better performance.

## Scaling Guidelines

### When to Scale Up (More Agents)

Add a 3rd agent if you experience:
- ❌ Pods in Pending state (insufficient resources)
- ❌ Node CPU/Memory pressure
- ❌ Slow pod scheduling
- ❌ Want to test complex scheduling scenarios

```bash
# Recreate with 3 agents
k3d cluster delete ecommerce-cluster
# Edit k3d-config.yaml: agents: 3
k3d cluster create --config k3d-config.yaml
```

### When to Scale Down (Fewer Agents)

Use 1 agent if:
- ✅ MacBook has only 16GB RAM
- ✅ Running many other applications
- ✅ Developing only 1-2 services
- ✅ Don't need HA testing

```bash
# Edit k3d-config.yaml: agents: 1
```

## Cluster Management Commands

### Basic Operations

```bash
# Create cluster
k3d cluster create --config k3d-config.yaml

# Start/Stop (preserves data)
k3d cluster start ecommerce-cluster
k3d cluster stop ecommerce-cluster

# Delete (removes everything)
k3d cluster delete ecommerce-cluster

# List clusters
k3d cluster list

# Get kubeconfig
k3d kubeconfig get ecommerce-cluster
```

### Node Operations

```bash
# List nodes
kubectl get nodes

# Describe node
kubectl describe node k3d-ecommerce-cluster-agent-0

# Drain node (for maintenance)
kubectl drain k3d-ecommerce-cluster-agent-0 --ignore-daemonsets

# Uncordon node
kubectl uncordon k3d-ecommerce-cluster-agent-0
```

### Registry Operations

```bash
# Import local image
k3d image import myimage:latest -c ecommerce-cluster

# List images in cluster
docker exec k3d-ecommerce-cluster-server-0 crictl images
```

## Best Practices

### 1. Resource Limits

Always set resource limits in deployments:

```yaml
resources:
  requests:
    cpu: 100m
    memory: 256Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

### 2. Node Selectors

Distribute workloads intelligently:

```yaml
nodeSelector:
  node-role.kubernetes.io/agent: "true"
```

### 3. Health Checks

Always include health checks:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
```

### 4. Cleanup

Regularly clean unused resources:

```bash
# Remove unused images
docker system prune -a

# Delete completed pods
kubectl delete pods --field-selector status.phase=Succeeded -A
```

## Monitoring Cluster Health

```bash
# Node status
kubectl get nodes

# Node resources
kubectl top nodes

# Pod distribution
kubectl get pods -A -o wide

# Events
kubectl get events -A --sort-by='.lastTimestamp'

# Cluster info
kubectl cluster-info
kubectl cluster-info dump
```

## Troubleshooting

### Cluster Won't Start

```bash
# Check Docker
docker ps

# Check available ports
lsof -i :80 -i :443 -i :5432 -i :6379 -i :9092

# View k3d logs
k3d cluster list
docker logs k3d-ecommerce-cluster-server-0
```

### Out of Resources

```bash
# Check resource usage
kubectl top nodes
kubectl describe nodes

# Scale down non-essential services
kubectl scale deployment -n monitoring --replicas=0 --all

# Increase Docker Desktop memory allocation
```

### Network Issues

```bash
# Restart cluster
k3d cluster stop ecommerce-cluster
k3d cluster start ecommerce-cluster

# Recreate cluster
k3d cluster delete ecommerce-cluster
k3d cluster create --config k3d-config.yaml
```

## Summary

The **1 Server + 2 Agents** configuration provides:

- ✅ Optimal balance for MacBook development
- ✅ Production-like multi-node environment
- ✅ Sufficient resources for 9 microservices
- ✅ Room for databases, Kafka, monitoring
- ✅ Good learning experience
- ✅ Manageable resource consumption

Perfect for developing and testing cloud-native applications locally! 🚀
