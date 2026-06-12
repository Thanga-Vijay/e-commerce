# E-Commerce Platform - Complete Deployment Guide

**🎉 Production-Ready Enterprise E-Commerce Platform**

This comprehensive guide covers deployment, configuration, and operations for a complete cloud-native e-commerce platform with 14 phases implemented.

## 📖 Table of Contents

1. [Quick Start](#quick-start)
2. [Phase-by-Phase Deployment](#phase-by-phase-deployment)
3. [Production Deployment](#production-deployment)
4. [Security Configuration](#security-configuration)
5. [Monitoring & Observability](#monitoring--observability)
6. [CI/CD Setup](#cicd-setup)
7. [Disaster Recovery](#disaster-recovery)
8. [Operations & Maintenance](#operations--maintenance)
9. [Troubleshooting](#troubleshooting)

## Quick Start

### Prerequisites

**Required:**
- Kubernetes cluster (1.25+)
- kubectl configured
- Docker & Docker Compose
- Git

**Optional:**
- Helm 3+
- AWS CLI (for S3 backups)
- GitHub CLI

### 5-Minute Local Development

```bash
# 1. Clone repository
git clone <repository-url>
cd e-commerce

# 2. Start infrastructure
make infra-up

# 3. Start all services
make services-up

# 4. Start frontend
make frontend-up

# 5. Access application
# Frontend: http://localhost:3000
# API: http://localhost:8081-8089
```

### 15-Minute Kubernetes Deployment

```bash
# 1. Deploy to Kubernetes
cd k8s
./deploy.sh

# 2. Deploy monitoring
cd monitoring
./deploy-monitoring.sh

# 3. Apply security
cd ../security
kubectl apply -f .

# 4. Set up backups
cd ../backup
kubectl apply -f backup-cronjobs.yaml

# 5. Access services
kubectl get svc -n ecommerce
kubectl port-forward -n ecommerce svc/frontend-service 3000:80
```

## Phase-by-Phase Deployment

### Phase 1-9: Application Development (Complete)

**Components:**
- 9 microservices (Go + Gin)
- 9 PostgreSQL databases
- Redis cache
- Apache Kafka
- React frontend

**Already implemented** - Services ready for deployment

### Phase 10: Kubernetes Deployment

**Deploy infrastructure:**
```bash
cd k8s

# 1. Create namespace and configs
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml

# 2. Deploy databases (StatefulSets)
kubectl apply -f databases/

# 3. Deploy Redis and Kafka
kubectl apply -f redis/
kubectl apply -f kafka/

# 4. Deploy microservices
kubectl apply -f services/

# 5. Deploy frontend
kubectl apply -f frontend/

# 6. Configure ingress
kubectl apply -f ingress/

# 7. Enable auto-scaling
kubectl apply -f hpa/
```

**Verify:**
```bash
kubectl get all -n ecommerce
kubectl get pvc -n ecommerce
kubectl get ingress -n ecommerce
```

**Documentation:** [PHASE10_KUBERNETES.md](PHASE10_KUBERNETES.md)

### Phase 11: Monitoring & Observability

**Deploy monitoring stack:**
```bash
cd k8s/monitoring

# Automated deployment
./deploy-monitoring.sh

# Or manual
kubectl apply -f prometheus-config.yaml
kubectl apply -f prometheus-rules.yaml
kubectl apply -f prometheus.yaml
kubectl apply -f alertmanager.yaml
kubectl apply -f grafana.yaml
kubectl apply -f loki.yaml
kubectl apply -f jaeger.yaml
```

**Access dashboards:**
```bash
# Grafana
kubectl port-forward -n ecommerce svc/grafana-service 3000:3000
# http://localhost:3000 (admin/admin)

# Prometheus
kubectl port-forward -n ecommerce svc/prometheus-service 9090:9090

# Jaeger
kubectl port-forward -n ecommerce svc/jaeger-query-service 16686:16686
```

**Documentation:** [PHASE11_MONITORING.md](PHASE11_MONITORING.md)

### Phase 12: CI/CD Pipeline

**Setup GitHub Actions:**

1. **Add secrets to GitHub repository:**
   ```
   KUBECONFIG_STAGING
   KUBECONFIG_PRODUCTION
   SNYK_TOKEN
   SLACK_WEBHOOK
   ```

2. **Workflows auto-run on:**
   - Pull requests → CI pipeline (lint, test, scan)
   - Push to main → CD to staging
   - Git tags → CD to production (with approval)

3. **Manual trigger:**
   ```bash
   gh workflow run cd.yml
   gh run list
   ```

**Workflows:**
- `.github/workflows/ci.yml` - Testing & validation
- `.github/workflows/cd.yml` - Build & deployment
- `.github/workflows/security.yml` - Security scanning

**Documentation:** [PHASE12_CICD.md](PHASE12_CICD.md)

### Phase 13: Security & Hardening

**Apply security policies:**
```bash
cd k8s/security

# 1. Network policies (zero-trust)
kubectl apply -f network-policies.yaml

# 2. RBAC
kubectl apply -f rbac.yaml

# 3. Pod security + PDBs
kubectl apply -f pod-security.yaml

# 4. Update deployments with ServiceAccounts
for svc in auth product cart wishlist order payment inventory notification reporting; do
  kubectl patch deployment ${svc}-service -n ecommerce \
    -p '{"spec":{"template":{"spec":{"serviceAccountName":"ecommerce-app"}}}}'
done
```

**Enable Pod Security Standards:**
```bash
kubectl label namespace ecommerce \
  pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/audit=restricted \
  pod-security.kubernetes.io/warn=restricted
```

**Verify:**
```bash
kubectl get networkpolicies -n ecommerce
kubectl get serviceaccounts -n ecommerce
kubectl get pdb -n ecommerce
kubectl describe ns ecommerce | grep pod-security
```

**Documentation:** [PHASE13_SECURITY.md](PHASE13_SECURITY.md)

### Phase 14: Disaster Recovery

**Set up automated backups:**
```bash
cd k8s/backup
kubectl apply -f backup-cronjobs.yaml

# Verify CronJobs
kubectl get cronjobs -n ecommerce
```

**Manual backup:**
```bash
# Full DR backup
./scripts/disaster-recovery.sh

# Database only
./scripts/backup-databases.sh

# With S3 upload
export BACKUP_S3_BUCKET=my-bucket
./scripts/backup-databases.sh
```

**Test restore:**
```bash
# Restore single database
./scripts/restore-database.sh \
  backups/20240612/auth_db_20240612_020000.sql.gz \
  auth_db

# DR drill (monthly)
# See PHASE14_DISASTER_RECOVERY.md
```

**Documentation:** [PHASE14_DISASTER_RECOVERY.md](PHASE14_DISASTER_RECOVERY.md)

## Production Deployment

### Pre-Deployment Checklist

- [ ] **Kubernetes cluster ready** (3+ nodes, 16GB+ RAM)
- [ ] **Domain configured** (DNS records)
- [ ] **TLS certificates** (cert-manager installed)
- [ ] **Secrets updated** (passwords, API keys)
- [ ] **Storage provisioned** (PersistentVolumes)
- [ ] **Monitoring configured** (Prometheus, Grafana)
- [ ] **Backup storage** (S3 bucket created)
- [ ] **CI/CD configured** (GitHub Actions secrets)
- [ ] **Network policies** (firewall rules)
- [ ] **Disaster recovery tested** (DR drill completed)

### Production Configuration

**1. Update secrets:**
```bash
cd k8s
nano secrets.yaml

# Update ALL secrets:
# - POSTGRES_PASSWORD
# - REDIS_PASSWORD
# - JWT_SECRET
# - STRIPE_SECRET_KEY
# - STRIPE_WEBHOOK_SECRET
# - SMTP credentials
# - GRAFANA_ADMIN_PASSWORD

kubectl apply -f secrets.yaml
```

**2. Configure domain and TLS:**
```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Create ClusterIssuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@yourdomain.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

# Update ingress with your domain
nano k8s/ingress/ingress.yaml
# Change ecommerce.yourdomain.com to your domain
kubectl apply -f k8s/ingress/ingress.yaml
```

**3. Configure resource limits for production:**
```bash
# Adjust based on your workload
nano k8s/services/auth-service.yaml

resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

**4. Deploy:**
```bash
# One-command deployment
cd k8s
./deploy.sh

# Or push to main branch for automated deployment
git push origin main
```

### Post-Deployment Validation

```bash
# 1. Check all pods running
kubectl get pods -n ecommerce
# All should be Running

# 2. Check services
kubectl get svc -n ecommerce

# 3. Check ingress
kubectl get ingress -n ecommerce
kubectl describe ingress ecommerce-ingress -n ecommerce

# 4. Test endpoints
curl https://your-domain.com/api/v1/health

# 5. Check monitoring
kubectl port-forward -n ecommerce svc/grafana-service 3000:3000
# Open http://localhost:3000

# 6. Verify backups
kubectl get cronjobs -n ecommerce
kubectl logs -n ecommerce -l app=backup-job
```

## Security Configuration

### Immediate Security Actions

```bash
# 1. Rotate all secrets
kubectl create secret generic app-secrets-new \
  --from-literal=... \
  --dry-run=client -o yaml | kubectl apply -f -

# 2. Enable network policies
kubectl apply -f k8s/security/network-policies.yaml

# 3. Restrict pod permissions
kubectl apply -f k8s/security/pod-security.yaml

# 4. Audit cluster
kubectl get pods --all-namespaces -o json | \
  jq '.items[] | select(.spec.securityContext.runAsUser == 0)'
# Should return empty (no root containers)
```

### Security Scanning

```bash
# Run Trivy on all images
trivy image ghcr.io/your-org/ecommerce/auth-service:latest

# Check for CVEs
kubectl get images -A -o json | jq -r '.items[].status.containerStatuses[].image' | sort -u | xargs -I {} trivy image {}

# Security audit
kube-bench run --targets master,node
```

## Monitoring & Observability

### Key Dashboards

**Grafana (port 3000):**
- E-Commerce Overview
- Kubernetes Cluster
- PostgreSQL Performance
- Redis Metrics
- Kafka Monitoring

**Import dashboards:**
```bash
# In Grafana UI: Dashboards → Import
- Kubernetes Cluster: 7249
- PostgreSQL: 9628
- Redis: 11835
- Kafka: 7589
- Go Runtime: 10826
```

### Critical Alerts

**Configure Alertmanager:**
```bash
nano k8s/monitoring/alertmanager.yaml

# Add your Slack webhook
slack_configs:
- api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'

kubectl apply -f k8s/monitoring/alertmanager.yaml
```

**Alert rules:**
- PodDown (5min)
- HighCPUUsage (10min, >80%)
- HighMemoryUsage (5min, >90%)
- DatabaseConnectionPoolExhausted (5min, >80%)
- HighErrorRate (5min, >5%)
- KafkaConsumerLag (10min, >1000)

### Log Queries (Loki)

```logql
# All errors
{namespace="ecommerce"} |= "error"

# Slow queries
{app="order-service"} | json | duration > 1000ms

# Failed payments
{app="payment-service"} |= "payment" |= "failed"

# Authentication failures
{app="auth-service"} |= "authentication failed"
```

## CI/CD Setup

### GitHub Repository Setup

**1. Add secrets:**
```bash
# Navigate to: Settings → Secrets and variables → Actions

# Add these secrets:
KUBECONFIG_STAGING       # Base64 encoded kubeconfig
KUBECONFIG_PRODUCTION    # Base64 encoded kubeconfig
SNYK_TOKEN              # From snyk.io
SLACK_WEBHOOK           # Slack incoming webhook
AWS_ACCESS_KEY_ID       # For S3 backups (optional)
AWS_SECRET_ACCESS_KEY   # For S3 backups (optional)
```

**2. Set up environments:**
```bash
# Settings → Environments → New environment

# Staging:
# - No protection rules
# - Auto-deploy on main branch

# Production:
# - Required reviewers: 2
# - Only deploy on tags (v*.*.*)
```

**3. Branch protection:**
```bash
# Settings → Branches → Add rule

# For main:
# ✓ Require pull request reviews (1 approver)
# ✓ Require status checks to pass (CI pipeline)
# ✓ Require branches to be up to date
```

### Deployment Workflow

```bash
# 1. Feature development
git checkout -b feature/new-feature
# ... make changes ...
git commit -m "feat: add new feature"
git push origin feature/new-feature

# 2. Create PR → CI runs
# 3. Merge to develop → CI runs
# 4. Merge to main → Deploys to staging

# 5. Create release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# → Triggers production deployment (with approval)

# 6. Monitor deployment
gh run list --workflow=cd.yml
gh run view <run-id> --log
```

## Disaster Recovery

### Backup Schedule

- **Databases**: Daily at 2 AM (CronJob)
- **Redis**: Every 6 hours (CronJob)
- **Full DR**: Manual / weekly
- **Retention**: 30 days (PVC), 12 months (S3)

### Recovery Scenarios

**1. Single pod failure:**
```bash
# Kubernetes automatically restarts
kubectl get pods -n ecommerce -w
```

**2. Database corruption:**
```bash
./scripts/restore-database.sh \
  backups/latest/order_db_*.sql.gz \
  order_db
# RTO: 15 minutes, RPO: 24 hours
```

**3. Complete cluster failure:**
```bash
# See PHASE14_DISASTER_RECOVERY.md
# Full restoration procedure
# RTO: 2 hours, RPO: 24 hours
```

### Monthly DR Drill

**First Saturday of each month:**
```bash
# 1. Create test namespace
kubectl create namespace ecommerce-dr-test

# 2. Restore to test namespace
# (modify scripts to use test namespace)

# 3. Validate data integrity
./scripts/validate-backup.sh

# 4. Test API functionality
curl http://test-api/health

# 5. Document results

# 6. Clean up
kubectl delete namespace ecommerce-dr-test
```

## Operations & Maintenance

### Daily Operations

```bash
# Check cluster health
kubectl get nodes
kubectl top nodes
kubectl get pods -n ecommerce

# Check HPA status
kubectl get hpa -n ecommerce

# View recent logs
kubectl logs -n ecommerce -l app=auth-service --tail=100

# Check backup status
kubectl get cronjobs -n ecommerce
kubectl logs -n ecommerce job/database-backup-xxx
```

### Weekly Operations

```bash
# Review metrics in Grafana
# Review alerts in Alertmanager
# Check backup integrity
sha256sum -c backups/latest/checksums.txt

# Update container images
kubectl set image deployment/auth-service \
  auth-service=ghcr.io/your-org/ecommerce/auth-service:v1.0.1

# Review resource usage
kubectl top pods -n ecommerce --sort-by=memory
```

### Monthly Operations

```bash
# Run DR drill
# Update dependencies
# Review and rotate secrets
# Capacity planning review
# Security audit
kube-bench run
# Review access logs
```

## Troubleshooting

### Common Issues

**Pods not starting:**
```bash
kubectl describe pod <pod-name> -n ecommerce
kubectl logs <pod-name> -n ecommerce
kubectl get events -n ecommerce --sort-by='.lastTimestamp'
```

**Service unreachable:**
```bash
kubectl get svc <service-name> -n ecommerce
kubectl get endpoints <service-name> -n ecommerce
kubectl run test --image=busybox --rm -it -n ecommerce -- wget -O- http://service:port/health
```

**Database connection issues:**
```bash
kubectl exec -it auth-db-0 -n ecommerce -- psql -U postgres -l
kubectl logs auth-service-xxx -n ecommerce | grep -i "database"
```

**Ingress not working:**
```bash
kubectl get ingress -n ecommerce
kubectl describe ingress ecommerce-ingress -n ecommerce
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller
```

**Monitoring issues:**
```bash
# Check Prometheus targets
kubectl port-forward -n ecommerce svc/prometheus-service 9090:9090
# Visit http://localhost:9090/targets

# Check Grafana datasources
kubectl logs -n ecommerce -l app=grafana
```

### Debug Commands

```bash
# Enter pod shell
kubectl exec -it <pod-name> -n ecommerce -- /bin/sh

# Port forward to service
kubectl port-forward -n ecommerce svc/auth-service 8081:8081

# View resource usage
kubectl top pods -n ecommerce
kubectl top nodes

# Check persistent volumes
kubectl get pv
kubectl get pvc -n ecommerce

# Describe resource
kubectl describe <resource-type> <name> -n ecommerce

# View events
kubectl get events -n ecommerce --sort-by='.lastTimestamp'
```

## Resources

### Documentation
- [Architecture Overview](docs/PHASE1_ARCHITECTURE.md)
- [Docker Setup](PHASE8_SETUP.md)
- [Kafka Events](PHASE9_KAFKA.md)
- [Kubernetes Deployment](PHASE10_KUBERNETES.md)
- [Monitoring & Observability](PHASE11_MONITORING.md)
- [CI/CD Pipeline](PHASE12_CICD.md)
- [Security & Hardening](PHASE13_SECURITY.md)
- [Disaster Recovery](PHASE14_DISASTER_RECOVERY.md)

### External Resources
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [Go Documentation](https://go.dev/doc/)
- [React Documentation](https://react.dev/)

### Support
- GitHub Issues
- Documentation Wiki
- Team Slack Channel
- Email: support@example.com

---

**Built with ❤️ using Go, React, Kubernetes, and Cloud Native technologies**

**Production-Ready · Secure · Scalable · Observable · Resilient**
