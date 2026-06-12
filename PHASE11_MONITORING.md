# E-Commerce Platform - Phase 11: Monitoring & Observability

Complete monitoring, logging, and tracing infrastructure for production-grade observability.

## Overview

Phase 11 implements a comprehensive observability stack:
- **Prometheus** - Metrics collection and alerting
- **Grafana** - Visualization and dashboards
- **Loki** - Log aggregation and querying
- **Promtail** - Log shipping agent
- **Jaeger** - Distributed tracing
- **Alertmanager** - Alert routing and management

## Architecture

```
┌─────────────────────────────────────────────────────┐
│            Observability Stack                      │
│                                                     │
│  ┌─────────────┐    ┌──────────────┐              │
│  │  Grafana    │◄───┤  Prometheus  │              │
│  │  (Dashboards│    │  (Metrics)   │              │
│  └──────┬──────┘    └───────▲──────┘              │
│         │                   │                       │
│         │           ┌───────┴────────┐             │
│         │           │                │             │
│    ┌────▼─────┐  ┌─▼───────┐  ┌────▼─────┐       │
│    │   Loki   │  │Services │  │AlertMgr  │       │
│    │  (Logs)  │  │ /metrics│  │(Alerts)  │       │
│    └────▲─────┘  └─────────┘  └──────────┘       │
│         │                                          │
│    ┌────┴─────┐         ┌──────────────┐         │
│    │ Promtail │         │    Jaeger    │         │
│    │(Shipper) │         │  (Tracing)   │         │
│    └──────────┘         └──────────────┘         │
└─────────────────────────────────────────────────────┘
           ▲                        ▲
           │                        │
    ┌──────┴────────────────────────┴──────┐
    │      E-Commerce Microservices         │
    │   (Auth, Product, Cart, Order, etc.)  │
    └───────────────────────────────────────┘
```

## Components

### 1. Prometheus (Metrics)

**Purpose:** Scrapes metrics from services, stores time-series data, evaluates alert rules

**Features:**
- 15s scrape interval
- 30-day retention
- Auto-discovery via Kubernetes service discovery
- Custom scrape configs for all services
- 50Gi persistent storage

**Access:** `http://<prometheus-ip>:9090`

**Key Metrics Collected:**
- HTTP request rate, latency, errors
- CPU and memory usage
- Pod availability and restart rate
- Database connection pools
- Kafka consumer lag
- Redis operations

### 2. Grafana (Visualization)

**Purpose:** Visualizes metrics, logs, and traces in unified dashboards

**Features:**
- Pre-configured datasources (Prometheus, Loki, Jaeger)
- Pre-built dashboard for e-commerce platform
- 10Gi persistent storage
- LoadBalancer service for external access

**Access:** `http://<grafana-ip>:3000`
**Credentials:** admin / (set in secrets)

**Dashboards Included:**
- E-Commerce Platform Overview
- Service-level metrics
- Infrastructure metrics
- Database performance
- Kafka monitoring

### 3. Loki + Promtail (Logs)

**Purpose:** Aggregates logs from all pods, provides log search and correlation

**Features:**
- DaemonSet deployment (runs on every node)
- 30-day log retention
- 20Gi storage
- Label-based log indexing
- Integration with Grafana

**Promtail:** Collects logs from `/var/log/pods` and ships to Loki

**Query Examples:**
```logql
# All logs from auth service
{app="auth-service"}

# Error logs across all services
{namespace="ecommerce"} |= "error"

# Logs with response time > 1s
{app="product-service"} | json | duration > 1s
```

### 4. Jaeger (Distributed Tracing)

**Purpose:** Traces requests across microservices, identifies bottlenecks

**Features:**
- All-in-one deployment (collector, query, UI)
- 20Gi persistent storage (Badger DB)
- OpenTelemetry support
- LoadBalancer for UI access

**Access:** `http://<jaeger-ip>:16686`

**Ports:**
- 16686: UI
- 14268: HTTP collector
- 14250: gRPC collector
- 4317: OTLP gRPC
- 4318: OTLP HTTP
- 9411: Zipkin compatibility

### 5. Alertmanager (Alerting)

**Purpose:** Routes alerts from Prometheus to notification channels

**Features:**
- Email notifications
- Slack integration (configure webhook)
- Alert grouping and deduplication
- Inhibition rules (suppress lower-severity alerts)
- Critical/Warning routing

**Access:** `http://<cluster-ip>:9093` (ClusterIP only)

## Alert Rules

### Critical Alerts

1. **PodDown** - Pod unavailable for 5+ minutes
2. **DatabaseConnectionPoolExhausted** - >80% connections used
3. **HighErrorRate** - Error rate >5%
4. **RedisDown** - Redis unavailable

### Warning Alerts

1. **HighCPUUsage** - CPU >80% for 10 minutes
2. **HighMemoryUsage** - Memory >90% for 5 minutes
3. **HighPodRestartRate** - Frequent restarts
4. **HighAPILatency** - P95 latency >2s
5. **KafkaConsumerLag** - Lag >1000 messages
6. **LowDiskSpace** - <10% disk space

## Deployment

### Quick Start

```bash
cd k8s/monitoring

# Deploy entire monitoring stack
chmod +x deploy-monitoring.sh
./deploy-monitoring.sh
```

### Manual Deployment

```bash
# 1. Deploy Prometheus
kubectl apply -f prometheus-config.yaml
kubectl apply -f prometheus-rules.yaml
kubectl apply -f prometheus.yaml

# 2. Deploy Alertmanager
kubectl apply -f alertmanager.yaml

# 3. Deploy Grafana
kubectl apply -f grafana.yaml

# 4. Deploy Loki
kubectl apply -f loki.yaml

# 5. Deploy Jaeger
kubectl apply -f jaeger.yaml

# 6. Deploy ServiceMonitors (if using Prometheus Operator)
kubectl apply -f servicemonitors.yaml
```

### Verify Deployment

```bash
# Check all monitoring pods
kubectl get pods -n ecommerce | grep -E "prometheus|grafana|loki|jaeger|promtail|alertmanager"

# Check services
kubectl get svc -n ecommerce | grep -E "prometheus|grafana|loki|jaeger|alertmanager"

# Check PVCs
kubectl get pvc -n ecommerce | grep -E "prometheus|grafana|loki|jaeger"
```

## Configuration

### Update Alertmanager Notifications

Edit `alertmanager.yaml`:

```yaml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'your-email@company.com'
  smtp_auth_username: 'your-email@company.com'
  smtp_auth_password: 'your-app-password'

receivers:
- name: 'team-ops-critical'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
    channel: '#alerts-critical'
```

Then apply:
```bash
kubectl apply -f alertmanager.yaml
kubectl rollout restart deployment/alertmanager -n ecommerce
```

### Add Grafana Admin Password

```bash
# Create/update secret
kubectl create secret generic app-secrets \
  --from-literal=GRAFANA_ADMIN_PASSWORD='your-secure-password' \
  -n ecommerce \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart Grafana
kubectl rollout restart deployment/grafana -n ecommerce
```

### Add Prometheus Scrape Annotations

Services already have annotations. For custom pods:

```yaml
metadata:
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
    prometheus.io/path: "/metrics"
```

## Access Services

### Port Forwarding (Local Development)

```bash
# Grafana
kubectl port-forward -n ecommerce svc/grafana-service 3000:3000

# Prometheus
kubectl port-forward -n ecommerce svc/prometheus-service 9090:9090

# Jaeger UI
kubectl port-forward -n ecommerce svc/jaeger-query-service 16686:16686

# Alertmanager
kubectl port-forward -n ecommerce svc/alertmanager-service 9093:9093
```

Then access:
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Jaeger: http://localhost:16686
- Alertmanager: http://localhost:9093

### LoadBalancer (Cloud Clusters)

```bash
# Get Grafana IP
kubectl get svc grafana-service -n ecommerce

# Get Jaeger IP
kubectl get svc jaeger-query-service -n ecommerce
```

## Using Grafana

### Import Pre-built Dashboards

1. Access Grafana UI
2. Navigate to **Dashboards** → **Import**
3. Import by ID:
   - **Kubernetes Cluster Monitoring**: 7249
   - **Kubernetes Pod Monitoring**: 6417
   - **PostgreSQL Database**: 9628
   - **Redis**: 11835
   - **Kafka Overview**: 7589
   - **Go Runtime Metrics**: 10826

### Create Custom Dashboard

1. Go to **Dashboards** → **New Dashboard**
2. Add Panel
3. Select **Prometheus** datasource
4. Use PromQL queries:

```promql
# Request rate by service
sum(rate(http_requests_total{namespace="ecommerce"}[5m])) by (service)

# P95 latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{namespace="ecommerce"}[5m])) by (service, le))

# Error rate
sum(rate(http_requests_total{status=~"5..", namespace="ecommerce"}[5m])) by (service) / sum(rate(http_requests_total{namespace="ecommerce"}[5m])) by (service)

# Pod CPU usage
rate(container_cpu_usage_seconds_total{namespace="ecommerce"}[5m])

# Pod memory usage
container_memory_usage_bytes{namespace="ecommerce"}
```

## Using Loki (Logs)

### Query Logs in Grafana

1. Go to **Explore**
2. Select **Loki** datasource
3. Use LogQL:

```logql
# All logs from a service
{app="auth-service"}

# Filter by log level
{app="auth-service"} |= "ERROR"

# JSON parsing
{app="product-service"} | json | line_format "{{.message}}"

# Metric from logs (error count)
sum(rate({namespace="ecommerce"} |= "error" [5m])) by (app)

# Slow queries
{app="order-service"} | json | duration > 1000ms
```

### Common Queries

```logql
# All 5xx errors
{namespace="ecommerce"} |~ "status\":(5..)"

# Failed database connections
{namespace="ecommerce"} |~ "database.*connection.*failed"

# Authentication failures
{app="auth-service"} |= "authentication failed"

# Payment errors
{app="payment-service"} |= "payment" |= "error"
```

## Using Jaeger (Tracing)

### Instrument Services

Services need Jaeger client library. For Go:

```go
import (
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
)

func initJaeger(serviceName string) io.Closer {
    cfg := config.Configuration{
        ServiceName: serviceName,
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans:          true,
            LocalAgentHostPort: "jaeger-agent-service.ecommerce:6831",
        },
    }
    tracer, closer, err := cfg.NewTracer()
    if err != nil {
        panic(err)
    }
    opentracing.SetGlobalTracer(tracer)
    return closer
}
```

### Search Traces

1. Access Jaeger UI
2. Select service (e.g., "order-service")
3. Click "Find Traces"
4. Filter by:
   - Operation
   - Tags (e.g., `http.status_code=500`)
   - Duration (e.g., `>2s`)
5. Click trace to see span details

### Analyze Performance

- **Span duration**: Time spent in each service
- **Service dependencies**: Visualize call graph
- **Critical path**: Identify slowest operations
- **Error traces**: Filter by `error=true` tag

## Metrics Endpoints

### Add Metrics to Go Services

Use Prometheus client library:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

// Expose metrics endpoint
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Instrument Middleware

```go
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}

router.Use(PrometheusMiddleware())
```

## Troubleshooting

### Prometheus Not Scraping

```bash
# Check Prometheus targets
kubectl port-forward -n ecommerce svc/prometheus-service 9090:9090
# Open http://localhost:9090/targets

# Check pod annotations
kubectl get pod <pod-name> -n ecommerce -o yaml | grep -A 3 annotations
```

### Grafana Datasource Not Working

```bash
# Check datasource config
kubectl get cm grafana-datasources -n ecommerce -o yaml

# Check Prometheus service
kubectl get svc prometheus-service -n ecommerce

# Test connectivity from Grafana pod
kubectl exec -it <grafana-pod> -n ecommerce -- wget -O- http://prometheus-service:9090/api/v1/query?query=up
```

### Loki Not Receiving Logs

```bash
# Check Promtail pods
kubectl get pods -n ecommerce -l app=promtail

# Check Promtail logs
kubectl logs -n ecommerce -l app=promtail

# Check Loki service
kubectl get svc loki-service -n ecommerce

# Test log ingestion
kubectl logs <app-pod> -n ecommerce
```

### Jaeger Not Showing Traces

```bash
# Check Jaeger pod
kubectl get pod -l app=jaeger -n ecommerce

# Check collector service
kubectl get svc jaeger-collector-service -n ecommerce

# Verify service instrumentation (check for Jaeger env vars)
kubectl exec -it <service-pod> -n ecommerce -- env | grep JAEGER
```

## Best Practices

### 1. Metric Naming
- Use consistent prefixes: `http_`, `db_`, `kafka_`
- Include units in names: `_seconds`, `_bytes`, `_total`
- Use labels for dimensions, not metric names

### 2. Alert Tuning
- Start with longer evaluation periods, tune down
- Use `for` clause to reduce flapping
- Group related alerts
- Set appropriate severity levels

### 3. Dashboard Design
- One dashboard per service
- Include RED metrics (Rate, Errors, Duration)
- Add resource usage panels
- Use templating for multi-service dashboards

### 4. Log Retention
- 30 days for most logs
- 90+ days for audit logs
- Archive to S3/GCS for long-term storage

### 5. Trace Sampling
- Use adaptive sampling in production
- 100% sampling in development
- Sample by service importance

## Cost Optimization

### Reduce Storage

```yaml
# Prometheus retention (default 30d)
--storage.tsdb.retention.time=15d

# Loki retention
retention_period: 360h  # 15 days
```

### Reduce Scrape Frequency

```yaml
# Less critical services
scrape_interval: 30s  # Instead of 15s
```

### Use Recording Rules

Pre-compute expensive queries:

```yaml
groups:
- name: recording_rules
  interval: 30s
  rules:
  - record: job:http_requests:rate5m
    expr: sum(rate(http_requests_total[5m])) by (job)
```

## Integration with CI/CD

### Export Grafana Dashboards

```bash
# Export dashboard JSON
kubectl exec -it <grafana-pod> -n ecommerce -- \
  curl -u admin:password \
  http://localhost:3000/api/dashboards/uid/ecommerce-overview

# Store in Git for version control
```

### Alert on Deployments

```yaml
- alert: DeploymentFailed
  expr: kube_deployment_status_replicas_unavailable{namespace="ecommerce"} > 0
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "Deployment has unavailable replicas"
```

## Summary

Phase 11 delivers production-grade observability:
- ✅ **Prometheus** - 50Gi storage, 30-day retention, 10+ alert rules
- ✅ **Grafana** - Pre-configured datasources, custom dashboard
- ✅ **Loki + Promtail** - Centralized logging, 30-day retention
- ✅ **Jaeger** - Distributed tracing, OTLP support
- ✅ **Alertmanager** - Email + Slack notifications
- ✅ **ServiceMonitors** - Auto-discovery for Prometheus Operator
- ✅ **Deployment script** - One-command setup
- ✅ **Comprehensive documentation** - Setup, usage, troubleshooting

**Next Steps:**
- Instrument services with Prometheus metrics
- Configure Jaeger tracing in microservices
- Set up notification channels in Alertmanager
- Import pre-built Grafana dashboards
- Create SLO/SLI dashboards

Your platform now has full observability! 🎯📊
