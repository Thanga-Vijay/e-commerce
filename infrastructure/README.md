# Infrastructure

Infrastructure as Code for the E-Commerce Platform.

## Structure

```
infrastructure/
├── docker/              # Docker configurations
├── kind/               # KIND cluster setup
├── kubernetes/         # Kubernetes manifests
├── monitoring/         # Prometheus, Grafana, Loki
├── kafka/              # Kafka configurations
├── redis/              # Redis configurations
└── github-actions/     # CI/CD workflows
```

## Components

### Docker
- Dockerfiles for all services
- Docker Compose for local development
- Multi-stage builds for optimization

### Kubernetes (KIND)
- Namespace configurations
- Deployments
- Services
- ConfigMaps
- Secrets
- Ingress
- HPA (Horizontal Pod Autoscaler)
- PersistentVolumeClaims

### Monitoring
- Prometheus for metrics
- Grafana dashboards
- Loki for log aggregation
- OpenTelemetry for tracing

### Kafka
- Kafka cluster setup
- Zookeeper configuration
- Topic creation scripts
- Schema registry

### Redis
- Redis cluster configuration
- Sentinel for high availability

### CI/CD
- GitHub Actions workflows
- Build and test pipelines
- Docker image building
- Kubernetes deployment automation

## Quick Start

### Local Development with Docker Compose
```bash
cd infrastructure/docker
docker-compose up -d
```

### KIND Cluster Setup
```bash
cd infrastructure/kind
./setup-cluster.sh
```

### Deploy to Kubernetes
```bash
cd infrastructure/kubernetes
kubectl apply -f namespace.yaml
kubectl apply -f .
```

### Setup Monitoring
```bash
cd infrastructure/monitoring
./setup-monitoring.sh
```

Test