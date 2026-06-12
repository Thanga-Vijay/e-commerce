# KIND Kubernetes Cluster

KIND (Kubernetes IN Docker) configuration for local development.

## Files

- `kind-config.yaml` - KIND cluster configuration
- `setup-cluster.sh` - Script to create cluster
- `destroy-cluster.sh` - Script to destroy cluster
- `load-images.sh` - Script to load Docker images into KIND

## Cluster Configuration

- 1 Control plane node
- 3 Worker nodes
- Ingress enabled
- Port mappings for services

## Setup

Create cluster:
```bash
./setup-cluster.sh
```

This will:
1. Create KIND cluster
2. Install ingress controller
3. Set up load balancer
4. Configure port forwarding

Destroy cluster:
```bash
./destroy-cluster.sh
```

Load images:
```bash
./load-images.sh
```

## Access Services

After deployment:
- Frontend: http://localhost
- API Gateway: http://localhost/api
- Grafana: http://localhost/grafana
- Prometheus: http://localhost/prometheus
