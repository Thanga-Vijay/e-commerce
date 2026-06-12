# Utility Scripts

This directory contains utility scripts for development, deployment, and maintenance.

## Available Scripts

### Development
- `setup-dev.sh` - Set up local development environment
- `generate-mocks.sh` - Generate test mocks
- `run-tests.sh` - Run all service tests

### Docker
- `build-all.sh` - Build all Docker images
- `push-all.sh` - Push all images to registry
- `docker-compose-up.sh` - Start all services with Docker Compose

### Kubernetes
- `kind-setup.sh` - Set up KIND cluster
- `deploy-all.sh` - Deploy all services to K8s
- `destroy-cluster.sh` - Destroy KIND cluster

### Database
- `run-migrations.sh` - Run database migrations
- `seed-data.sh` - Seed initial data
- `backup-db.sh` - Backup databases

### Monitoring
- `setup-monitoring.sh` - Set up Prometheus, Grafana, Loki
- `port-forward-grafana.sh` - Port forward to Grafana

### Kafka
- `setup-kafka.sh` - Set up Kafka cluster
- `create-topics.sh` - Create Kafka topics

## Usage

Make scripts executable:
```bash
chmod +x scripts/*.sh
```

Run a script:
```bash
./scripts/setup-dev.sh
```
