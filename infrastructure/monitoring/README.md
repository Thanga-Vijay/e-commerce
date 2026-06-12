# Monitoring Stack

Prometheus, Grafana, and Loki configuration for observability.

## Components

### Prometheus
- Metrics collection
- Service discovery
- Alert rules
- Recording rules

### Grafana
- Dashboards
- Data source configuration
- Alert visualization

### Loki
- Log aggregation
- Log querying
- Integration with Grafana

### OpenTelemetry
- Distributed tracing
- Service mesh observability

## Dashboards

### Application Dashboards
- API Gateway Metrics
- Service Health Overview
- Order Flow Dashboard
- Payment Processing Dashboard
- Inventory Dashboard

### Infrastructure Dashboards
- Kubernetes Cluster Overview
- Node Metrics
- Pod Metrics
- Kafka Metrics
- Redis Metrics
- PostgreSQL Metrics

### Business Dashboards
- Revenue Analytics
- Order Analytics
- Customer Analytics
- Product Performance

## Access

Grafana:
```bash
kubectl port-forward svc/grafana 3000:3000 -n monitoring
```
URL: http://localhost:3000
Default credentials: admin / admin

Prometheus:
```bash
kubectl port-forward svc/prometheus 9090:9090 -n monitoring
```
URL: http://localhost:9090

## Alerts

### Critical Alerts
- Service Down
- High Error Rate (>5%)
- High Latency (p99 >1s)
- Database Connection Pool Exhausted
- Kafka Consumer Lag >1000

### Warning Alerts
- High Memory Usage (>80%)
- High CPU Usage (>80%)
- Disk Space Low (<20%)
- Payment Failures (>10/min)

## Setup

```bash
./setup-monitoring.sh
```

This will:
1. Create monitoring namespace
2. Deploy Prometheus
3. Deploy Grafana
4. Deploy Loki
5. Import dashboards
6. Configure data sources
