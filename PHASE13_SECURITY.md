# E-Commerce Platform - Phase 13: Security & Hardening

Comprehensive security implementation with network policies, RBAC, pod security, and compliance measures.

## Overview

Phase 13 implements defense-in-depth security:
- **Network Policies**: Zero-trust networking, pod-to-pod restrictions
- **RBAC**: Role-based access control, least privilege
- **Pod Security**: Security contexts, admission control
- **Resource Management**: Quotas, limits, disruption budgets
- **Compliance**: Security best practices, CIS benchmarks

## Security Architecture

```
┌─────────────────────────────────────────────────┐
│            Security Layers                      │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │  1. Network Policies (L4 Firewall)      │  │
│  │     - Default deny all                   │  │
│  │     - Explicit allow rules               │  │
│  └──────────────────────────────────────────┘  │
│                      │                          │
│  ┌──────────────────▼──────────────────────┐  │
│  │  2. RBAC (Authentication & Authorization)│  │
│  │     - ServiceAccounts                    │  │
│  │     - Roles & RoleBindings              │  │
│  └──────────────────────────────────────────┘  │
│                      │                          │
│  ┌──────────────────▼──────────────────────┐  │
│  │  3. Pod Security (Runtime Security)     │  │
│  │     - SecurityContext                    │  │
│  │     - PodDisruptionBudgets              │  │
│  └──────────────────────────────────────────┘  │
│                      │                          │
│  ┌──────────────────▼──────────────────────┐  │
│  │  4. Resource Management                  │  │
│  │     - ResourceQuotas                     │  │
│  │     - LimitRanges                       │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

## Components

### 1. Network Policies

**File**: `k8s/security/network-policies.yaml`

#### Default Deny
Blocks all ingress traffic by default (zero-trust):
```yaml
kind: NetworkPolicy
metadata:
  name: default-deny-ingress
spec:
  podSelector: {}  # Applies to all pods
  policyTypes:
  - Ingress
```

#### Allowed Traffic Flows

| Source | Destination | Port | Purpose |
|--------|-------------|------|---------|
| Frontend | All Backend Services | 8081-8089 | API calls |
| Ingress Controller | Frontend | 80 | HTTP traffic |
| Backend Services | PostgreSQL | 5432 | Database queries |
| Backend Services | Redis | 6379 | Cache access |
| Backend Services | Kafka | 9092 | Event streaming |
| Prometheus | All Services | 8080-8089 | Metrics scraping |
| Promtail | Loki | 3100 | Log shipping |
| All Services | Jaeger | 14268, 6831 | Trace collection |
| All Pods | kube-dns | 53 | DNS resolution |
| Payment/Notification | External | 443, 587 | Stripe, SMTP |

#### Benefits
- ✅ Prevents lateral movement
- ✅ Limits blast radius of compromised pod
- ✅ Compliance with zero-trust principles
- ✅ Network segmentation within cluster

### 2. RBAC (Role-Based Access Control)

**File**: `k8s/security/rbac.yaml`

#### Service Accounts

**1. ecommerce-app** (Application Pods)
- **Permissions**: Read-only access to ConfigMaps, Secrets, Pods
- **Use**: Runtime service operations
- **Principle**: Least privilege

**2. cicd-deployer** (CI/CD Pipeline)
- **Permissions**: Full CRUD on deployments, services, configs
- **Use**: Automated deployments from GitHub Actions
- **Principle**: Deployment-only scope

**3. backup-job** (Backup Operations)
- **Permissions**: Read PVCs, exec into pods
- **Use**: Database backup jobs
- **Principle**: Backup-specific access

**4. prometheus** (Monitoring)
- **Permissions**: Read-only cluster-wide
- **Use**: Metrics collection
- **Principle**: Observer role

#### Applying RBAC

```bash
# Apply RBAC configuration
kubectl apply -f k8s/security/rbac.yaml

# Update deployments to use ServiceAccounts
kubectl patch deployment auth-service \
  -n ecommerce \
  -p '{"spec":{"template":{"spec":{"serviceAccountName":"ecommerce-app"}}}}'

# Verify ServiceAccount
kubectl get sa -n ecommerce
kubectl describe sa ecommerce-app -n ecommerce
```

### 3. Pod Security

**File**: `k8s/security/pod-security.yaml`

#### Pod Security Standards

Label namespace with security level:
```bash
# Enforce restricted policy (most secure)
kubectl label namespace ecommerce \
  pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/audit=restricted \
  pod-security.kubernetes.io/warn=restricted
```

**Restricted policy requires:**
- ✅ Run as non-root
- ✅ Drop all capabilities
- ✅ No privilege escalation
- ✅ Read-only root filesystem
- ✅ Seccomp profile

#### SecurityContext Example

```yaml
# Pod-level
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 2000
  seccompProfile:
    type: RuntimeDefault

# Container-level
securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000
```

#### PodDisruptionBudgets

Ensures high availability during maintenance:

| Service | Min Available | Purpose |
|---------|--------------|---------|
| auth-service | 1 | Always available |
| product-service | 2 | High traffic tolerance |
| order-service | 1 | Transaction integrity |
| payment-service | 1 | Critical service |
| frontend | 2 | User experience |
| databases | 0 unavailable | Data integrity |
| kafka | 2 | Cluster quorum |

**Apply PDB:**
```bash
kubectl apply -f k8s/security/pod-security.yaml

# Verify
kubectl get pdb -n ecommerce

# Test during drain
kubectl drain <node-name> --ignore-daemonsets
# PDBs prevent too many pods from terminating
```

### 4. Resource Management

#### ResourceQuota (Namespace-level)

Prevents resource exhaustion:
```yaml
spec:
  hard:
    requests.cpu: "50"          # 50 CPU cores total
    requests.memory: 100Gi      # 100 GB memory total
    limits.cpu: "100"           # Max 100 CPU cores
    limits.memory: 200Gi        # Max 200 GB memory
    persistentvolumeclaims: "50" # Max 50 PVCs
    services.loadbalancers: "5"  # Max 5 LoadBalancers
```

#### LimitRange (Pod/Container-level)

Enforces default and max limits:
```yaml
spec:
  limits:
  - max:
      cpu: "4"
      memory: "8Gi"
    min:
      cpu: "50m"
      memory: "64Mi"
    default:
      cpu: "500m"
      memory: "512Mi"
    defaultRequest:
      cpu: "100m"
      memory: "128Mi"
    type: Container
```

**Benefits:**
- Prevents single pod from consuming all resources
- Auto-applies limits to pods without specs
- Budget enforcement

## Deployment

### Quick Deploy All Security

```bash
cd k8s/security

# 1. Apply Network Policies
kubectl apply -f network-policies.yaml

# 2. Apply RBAC
kubectl apply -f rbac.yaml

# 3. Apply Pod Security
kubectl apply -f pod-security.yaml

# 4. Label namespace for Pod Security Standards
kubectl label namespace ecommerce \
  pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/audit=restricted \
  pod-security.kubernetes.io/warn=restricted

# 5. Update deployments with ServiceAccounts
for service in auth product cart wishlist order payment inventory notification reporting; do
  kubectl patch deployment ${service}-service -n ecommerce \
    -p '{"spec":{"template":{"spec":{"serviceAccountName":"ecommerce-app"}}}}'
done
```

### Verify Security Configuration

```bash
# Check Network Policies
kubectl get networkpolicies -n ecommerce
kubectl describe networkpolicy default-deny-ingress -n ecommerce

# Check RBAC
kubectl get serviceaccounts -n ecommerce
kubectl get roles -n ecommerce
kubectl get rolebindings -n ecommerce

# Check Pod Security
kubectl get pdb -n ecommerce
kubectl describe pdb auth-service-pdb -n ecommerce

# Check Resource Quotas
kubectl get resourcequota -n ecommerce
kubectl describe resourcequota ecommerce-quota -n ecommerce

# Check LimitRange
kubectl get limitrange -n ecommerce
kubectl describe limitrange ecommerce-limits -n ecommerce
```

## Security Best Practices

### 1. Secret Management

**External Secrets (Recommended)**
```bash
# Install External Secrets Operator
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets \
  external-secrets/external-secrets \
  -n external-secrets-system \
  --create-namespace

# Use AWS Secrets Manager, Azure Key Vault, or Vault
```

**Sealed Secrets**
```bash
# Install Sealed Secrets
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Seal secrets before committing to Git
kubeseal --format yaml < secret.yaml > sealed-secret.yaml
```

### 2. Image Security

**Image Pull Policy**
```yaml
spec:
  containers:
  - name: auth-service
    imagePullPolicy: Always  # Always pull for latest:tag
```

**Private Registry**
```bash
# Create image pull secret
kubectl create secret docker-registry regcred \
  --docker-server=ghcr.io \
  --docker-username=<username> \
  --docker-password=<token> \
  -n ecommerce

# Use in deployment
spec:
  imagePullSecrets:
  - name: regcred
```

### 3. TLS/SSL Everywhere

**Internal TLS (Service Mesh)**
```bash
# Install Istio for mTLS
istioctl install --set profile=production

# Enable mTLS
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: ecommerce
spec:
  mtls:
    mode: STRICT
EOF
```

### 4. Audit Logging

Enable Kubernetes audit logging:
```yaml
# /etc/kubernetes/audit-policy.yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
  namespaces: ["ecommerce"]
```

### 5. Security Scanning

**Continuous Scanning**
```bash
# Install Falco for runtime security
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm install falco falcosecurity/falco \
  --namespace falco \
  --create-namespace

# Monitor suspicious activity
kubectl logs -n falco -l app=falco
```

## Compliance

### CIS Kubernetes Benchmark

**Check compliance:**
```bash
# Install kube-bench
kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job.yaml

# View results
kubectl logs job/kube-bench

# Remediate issues
```

### OWASP Top 10

| Risk | Mitigation |
|------|-----------|
| Injection | Parameterized queries, input validation |
| Broken Auth | JWT with short expiry, MFA |
| Sensitive Data | Encryption at rest and transit |
| XML External Entities | Disable XML parsing |
| Broken Access Control | RBAC, Network Policies |
| Security Misconfiguration | Pod Security Standards |
| XSS | Content Security Policy |
| Insecure Deserialization | Input validation |
| Using Components with Known Vulnerabilities | Trivy scanning |
| Insufficient Logging | Loki + Prometheus |

### SOC 2 / ISO 27001

**Evidence collection:**
- Audit logs (Loki)
- Access logs (RBAC)
- Change logs (Git)
- Security scans (Trivy)
- Monitoring (Prometheus)

## Incident Response

### Security Incident Runbook

**1. Detect**
```bash
# Check Falco alerts
kubectl logs -n falco -l app=falco | grep "Warning"

# Check audit logs
kubectl logs -n kube-system -l component=kube-apiserver
```

**2. Isolate**
```bash
# Quarantine compromised pod
kubectl label pod <pod-name> quarantine=true

# Update NetworkPolicy to isolate
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: quarantine-policy
spec:
  podSelector:
    matchLabels:
      quarantine: "true"
  policyTypes:
  - Ingress
  - Egress
EOF
```

**3. Investigate**
```bash
# Get pod details
kubectl describe pod <pod-name> -n ecommerce

# Check logs
kubectl logs <pod-name> -n ecommerce --previous

# Export for forensics
kubectl cp <pod-name>:/app/logs /tmp/forensics/
```

**4. Remediate**
```bash
# Delete compromised pod
kubectl delete pod <pod-name> -n ecommerce

# Update deployment with patched image
kubectl set image deployment/<service> \
  <container>=<patched-image>:tag
```

**5. Recover**
```bash
# Restore from backup
./scripts/restore-database.sh <backup-file> <db-name>

# Verify system integrity
kubectl get pods -n ecommerce
```

## Testing Security

### Penetration Testing

```bash
# Install kube-hunter
kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-hunter/main/job.yaml

# View findings
kubectl logs job/kube-hunter
```

### Network Policy Testing

```bash
# Test from outside namespace (should fail)
kubectl run test-pod --image=busybox --rm -it -- \
  wget -O- http://auth-service.ecommerce:8081/health

# Test from inside namespace (should succeed if allowed)
kubectl run test-pod --image=busybox --rm -it -n ecommerce -- \
  wget -O- http://auth-service:8081/health
```

## Summary

Phase 13 delivers enterprise security:
- ✅ **Network Policies** - Zero-trust, 10+ policies, ingress/egress rules
- ✅ **RBAC** - 4 ServiceAccounts, least privilege access
- ✅ **Pod Security** - Restricted standards, 8 PodDisruptionBudgets
- ✅ **Resource Management** - ResourceQuota, LimitRange enforcement
- ✅ **Security Contexts** - Non-root, capability dropping, read-only FS
- ✅ **Compliance** - CIS Benchmark, OWASP Top 10 coverage
- ✅ **Incident Response** - Runbooks and procedures
- ✅ **Security Scanning** - Integrated with CI/CD

Your platform is now hardened for production! 🔒
