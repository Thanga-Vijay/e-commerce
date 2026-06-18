# AWS EKS Production Readiness Guide

Comprehensive answers to your pre-deployment questions.

---

## 1. ✅ Monorepo vs Multiple Repos

### Recommendation: **Keep Monorepo**

**Why it works for you:**
- ✅ All 9 microservices are tightly coupled (auth, products, cart, orders, etc.)
- ✅ Shared contracts and data models
- ✅ Atomic cross-service changes
- ✅ Easier to maintain consistency
- ✅ Simpler CI/CD (one pipeline, multiple jobs)
- ✅ Better for teams under 50 developers

**How to optimize monorepo for EKS:**
```yaml
# Use path filters to only build changed services
on:
  push:
    paths:
      - 'auth-service/**'
      - '.github/workflows/**'
```

**When to split:**
- Team grows to 50+ developers
- Services have completely independent lifecycles
- Different teams need separate repos
- Services deployed to different regions/clouds

**Real-world examples using monorepo:**
- Google (single repo with billions of lines)
- Microsoft (Windows)
- Facebook/Meta (React, React Native)

**Verdict: Keep your monorepo, optimize CI/CD** ✅

---

## 2. 🔧 Fix CI/CD Workflows

### Issues Found:

**Problem 1: Wrong Directory Paths**
```yaml
# WRONG (current):
working-directory: services/${{ matrix.service }}

# CORRECT (your structure):
working-directory: ${{ matrix.service }}-service
```

**Problem 2: Wrong Service Names**
```yaml
# WRONG:
service: [auth, product, cart, ...]

# CORRECT:
service: [auth-service, product-service, cart-service, ...]
```

**Problem 3: Registry Not Configured**
- Current: Uses GHCR (GitHub Container Registry)
- Need: AWS ECR for EKS production

### Fixed Workflows Ready

I'll create corrected workflows that:
- ✅ Use correct directory paths
- ✅ Support both GHCR (testing) and ECR (production)
- ✅ Only build changed services
- ✅ Work with your monorepo structure
- ✅ Support ArgoCD deployment

---

## 3. 📦 Unnecessary Files for EKS

### Files to REMOVE/Archive:

#### Docker Compose Files (Not needed for EKS):
```
❌ docker-compose.yml
❌ docker-compose.kafka.yml
❌ docker-compose.monitoring.yml
❌ docker-compose.prod.yml
❌ docker-compose.override.yml
❌ docker-compose.secrets.yml
```

#### Local Testing Scripts:
```
❌ scripts/backup.sh (Docker Compose backup)
❌ scripts/restore.sh (Docker Compose restore)
```

#### KIND-Specific Files:
```
❌ kind-config.yaml
❌ kind-setup.sh
❌ scripts/build-and-load-images.sh
```

#### Documentation for Local:
```
❌ KIND_QUICKSTART.md
❌ KIND_CHANGES_ANALYSIS.md
❌ KIND_ACTION_CHECKLIST.md
⚠️ PHASE2_SETUP.md through PHASE9_KAFKA.md (keep if historical reference needed)
```

#### Optional (depends on workflow):
```
⚠️ .env.prod.example (use AWS Secrets Manager instead)
⚠️ Makefile (if not used)
```

### Files to KEEP:

```
✅ k8s/ (all Kubernetes manifests)
✅ .github/workflows/ (CI/CD - will fix)
✅ scripts/backup-databases.sh (works with K8s)
✅ scripts/restore-database.sh (works with K8s)
✅ scripts/disaster-recovery.sh (works with K8s)
✅ scripts/health-check.sh (adapt for EKS)
✅ DEPLOYMENT_GUIDE.md
✅ eks-migration-guide.md
✅ PHASE10_KUBERNETES.md through PHASE14_DISASTER_RECOVERY.md
✅ README.md
✅ All service code directories
```

### Cleanup Script:

```bash
# Create backup first
mkdir ../ecommerce-backup
cp -r . ../ecommerce-backup/

# Remove Docker Compose files
rm docker-compose*.yml

# Remove KIND-specific files
rm kind-config.yaml kind-setup.sh
rm -f KIND_*.md

# Remove local scripts
rm scripts/backup.sh scripts/restore.sh scripts/build-and-load-images.sh

# Optional: Clean old phase docs (if not needed)
# rm PHASE2_SETUP.md PHASE3_SETUP.md ... PHASE9_KAFKA.md

echo "Cleaned for EKS production"
```

---

## 4. 🏗️ Terraform vs Manual EKS

### Recommendation: **Use Terraform** ✅

**Why Terraform:**

1. ✅ **Infrastructure as Code**
   - Version controlled
   - Reviewable changes
   - Reproducible environments

2. ✅ **Automation**
   - One command to create entire infrastructure
   - Easy to tear down and recreate
   - Consistent across environments (dev/staging/prod)

3. ✅ **State Management**
   - Tracks what exists
   - Prevents drift
   - Supports team collaboration

4. ✅ **AWS Best Practices Built-in**
   - Uses AWS modules maintained by HashiCorp
   - Security defaults
   - High availability patterns

**Terraform Modules to Use:**

```hcl
# terraform/main.tf
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"
  
  name = "ecommerce-vpc"
  cidr = "10.0.0.0/16"
  
  azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  
  enable_nat_gateway = true
  enable_vpn_gateway = false
  
  tags = {
    Environment = "production"
    Project     = "ecommerce"
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.16.0"
  
  cluster_name    = "ecommerce-prod"
  cluster_version = "1.30"
  
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets
  
  # Managed node groups
  eks_managed_node_groups = {
    application = {
      min_size     = 2
      max_size     = 6
      desired_size = 3
      
      instance_types = ["t3.large"]
      capacity_type  = "ON_DEMAND"
      
      labels = {
        tier = "application"
      }
    }
    
    database = {
      min_size     = 3
      max_size     = 5
      desired_size = 3
      
      instance_types = ["t3.xlarge"]
      capacity_type  = "ON_DEMAND"
      
      labels = {
        tier = "database"
      }
      
      taints = {
        dedicated = {
          key    = "database"
          value  = "true"
          effect = "NO_SCHEDULE"
        }
      }
    }
  }
  
  # Add-ons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
    aws-ebs-csi-driver = {
      most_recent = true
    }
  }
}

module "ecr" {
  source = "terraform-aws-modules/ecr/aws"
  
  for_each = toset([
    "auth-service",
    "product-service",
    "cart-service",
    "wishlist-service",
    "order-service",
    "payment-service",
    "inventory-service",
    "notification-service",
    "reporting-service",
    "frontend"
  ])
  
  repository_name = "ecommerce/${each.key}"
  
  repository_image_scan_on_push = true
  repository_lifecycle_policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 10 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 10
      }
      action = {
        type = "expire"
      }
    }]
  })
}
```

**Terraform Structure:**

```
terraform/
├── main.tf              # Main configuration
├── variables.tf         # Input variables
├── outputs.tf           # Output values (cluster info)
├── versions.tf          # Provider versions
├── backend.tf           # S3 backend for state
├── environments/
│   ├── dev.tfvars
│   ├── staging.tfvars
│   └── prod.tfvars
└── modules/
    ├── eks/
    ├── ecr/
    ├── rds/             # If using RDS instead of K8s DBs
    └── monitoring/
```

**Terraform vs Manual Comparison:**

| Feature | Terraform | Manual (eksctl/Console) |
|---------|-----------|------------------------|
| **Reproducibility** | ✅ Perfect | ❌ Hard to repeat |
| **Version Control** | ✅ Git tracked | ❌ Not tracked |
| **Multi-Environment** | ✅ Easy | ⚠️ Manual work |
| **Team Collaboration** | ✅ Easy | ❌ Difficult |
| **Cost Estimation** | ✅ terraform plan | ❌ Manual calc |
| **Drift Detection** | ✅ Built-in | ❌ Manual check |
| **Learning Curve** | ⚠️ Moderate | ✅ Easy (Console) |
| **Initial Setup** | ⚠️ 1-2 days | ✅ 2-3 hours |

**My Recommendation:**
- **For Production: Use Terraform** ✅
- **For Quick Testing: Use eksctl** (faster initial setup)

**Terraform Workflow:**

```bash
# 1. Initialize
cd terraform
terraform init

# 2. Plan (see what will be created)
terraform plan -var-file=environments/prod.tfvars

# 3. Apply (create infrastructure)
terraform apply -var-file=environments/prod.tfvars

# 4. Get cluster config
aws eks update-kubeconfig --name ecommerce-prod --region us-east-1

# 5. Verify
kubectl get nodes
```

**Cost Estimate:**
```
EKS Control Plane: $73/month
EC2 Nodes (6 × t3.large): ~$300/month
EBS Volumes: ~$50/month
Load Balancers: ~$20/month
Data Transfer: ~$30/month
Total: ~$470/month
```

---

## 5. 🚀 ArgoCD for GitOps

### Recommendation: **Yes, Use ArgoCD** ✅✅✅

**Why ArgoCD is Perfect for You:**

1. ✅ **GitOps Workflow**
   - Git is source of truth
   - Automatic sync from Git to cluster
   - Easy rollbacks (just revert Git commit)

2. ✅ **Multi-Environment**
   - Manage dev, staging, prod from one place
   - Different Git branches/folders per environment
   - Preview changes before applying

3. ✅ **Security**
   - No need to expose cluster credentials to CI/CD
   - Pull-based deployment (cluster pulls from Git)
   - RBAC integration

4. ✅ **Visibility**
   - Beautiful UI showing all deployments
   - Health status of all resources
   - Deployment history and diffs

5. ✅ **Automation**
   - Auto-sync on Git changes
   - Automatic pruning of deleted resources
   - Self-healing (restores if someone changes cluster directly)

**ArgoCD vs Traditional CD:**

| Feature | ArgoCD (GitOps) | GitHub Actions CD |
|---------|-----------------|-------------------|
| **Security** | ✅ Pull-based | ❌ Push needs cluster access |
| **Visibility** | ✅ Great UI | ⚠️ Logs only |
| **Multi-Cluster** | ✅ Easy | ❌ Complex |
| **Rollback** | ✅ Git revert | ⚠️ Manual/script |
| **Drift Detection** | ✅ Automatic | ❌ None |
| **Learning Curve** | ⚠️ Moderate | ✅ Easy |

**Recommended Architecture:**

```
GitHub (Git Repo)
    │
    ├─── Application Code
    │    └─── Triggers CI (build/test/push to ECR)
    │
    └─── K8s Manifests (k8s/)
         └─── ArgoCD watches this
              └─── Auto-deploys to EKS
```

**Workflow with ArgoCD:**

```bash
# Developer workflow:
1. git push code changes
2. GitHub Actions: Build → Test → Push to ECR
3. Update k8s manifests with new image tag
4. git push manifest changes
5. ArgoCD: Detects change → Deploys to EKS

# No kubectl needed! ArgoCD does it all.
```

**ArgoCD Setup:**

```bash
# 1. Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 2. Access ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443

# 3. Get admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# 4. Create application
argocd app create ecommerce \
  --repo https://github.com/yourusername/e-commerce.git \
  --path k8s \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace ecommerce \
  --sync-policy automated \
  --auto-prune \
  --self-heal
```

**ArgoCD Project Structure:**

```yaml
# argocd/application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: ecommerce-prod
  namespace: argocd
spec:
  project: default
  
  source:
    repoURL: https://github.com/yourusername/e-commerce.git
    targetRevision: main
    path: k8s
  
  destination:
    server: https://kubernetes.default.svc
    namespace: ecommerce
  
  syncPolicy:
    automated:
      prune: true      # Delete resources not in Git
      selfHeal: true   # Revert manual changes
      allowEmpty: false
    
    syncOptions:
    - CreateNamespace=true
    
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

**Benefits for Your Use Case:**

1. ✅ **9 Microservices**: ArgoCD shows health of all at once
2. ✅ **Multiple Databases**: Easy to see which are up/down
3. ✅ **Kafka/Redis**: Health monitoring included
4. ✅ **Monitoring Stack**: Manage Prometheus/Grafana with ArgoCD
5. ✅ **Rollbacks**: Just revert Git commit, ArgoCD redeploys

**ArgoCD + GitHub Actions Combo:**

```yaml
# .github/workflows/cd.yml (simplified with ArgoCD)
name: CD - Build and Update Manifests

on:
  push:
    branches: [main]

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - name: Build and push to ECR
        # ... build steps ...
      
      - name: Update K8s manifest
        run: |
          # Update image tag in deployment
          sed -i "s|image: .*auth-service:.*|image: $ECR_REGISTRY/auth-service:${{ github.sha }}|g" \
            k8s/services/auth-service.yaml
          
          # Commit and push
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add k8s/services/auth-service.yaml
          git commit -m "Update auth-service image to ${{ github.sha }}"
          git push
      
      # That's it! ArgoCD takes over from here.
```

---

## 🎯 Recommended Architecture

### Phase 1: Infrastructure (Week 1)
- ✅ Use **Terraform** to create EKS cluster
- ✅ Install **ArgoCD** on cluster
- ✅ Create **ECR repositories**
- ✅ Set up **VPC, subnets, security groups**

### Phase 2: CI/CD (Week 2)
- ✅ Fix **GitHub Actions** workflows (I'll provide corrected versions)
- ✅ Configure **ECR push** in CI
- ✅ Set up **ArgoCD applications**
- ✅ Test deployment pipeline

### Phase 3: Migration (Week 3)
- ✅ Deploy databases to EKS
- ✅ Deploy services one by one
- ✅ Set up monitoring (Prometheus/Grafana)
- ✅ Configure ingress/load balancers

### Phase 4: Production (Week 4)
- ✅ DNS configuration
- ✅ SSL/TLS certificates
- ✅ Backup automation
- ✅ Disaster recovery testing
- ✅ Go live!

---

## 📋 Action Items

### Immediate (Before AWS):
1. ✅ Fix CI/CD workflows (I'll create corrected versions)
2. ✅ Remove unnecessary files (Docker Compose, KIND files)
3. ✅ Test GitHub Actions (build, lint, security scan)

### Infrastructure Setup:
1. ✅ Learn Terraform basics (2-3 days)
2. ✅ Create Terraform configs for EKS (I can provide templates)
3. ✅ Provision EKS cluster with Terraform
4. ✅ Install ArgoCD on cluster

### Deployment:
1. ✅ Configure ArgoCD applications
2. ✅ Set up ECR push in CI
3. ✅ Deploy to EKS via ArgoCD
4. ✅ Monitor and optimize

---

## 🚦 Next Steps

**What do you want me to do next?**

A. **Fix CI/CD Workflows** - Create corrected ci.yml, cd.yml, security.yml that work with your structure

B. **Create Terraform Configs** - Full Terraform setup for EKS + ECR + VPC + all AWS resources

C. **ArgoCD Setup Guide** - Step-by-step guide to install and configure ArgoCD

D. **Clean Repository** - Remove unnecessary files and organize for production

E. **All of the Above** - Complete production-ready setup

Let me know which you want first, and I'll create it! 🚀

---

**My Recommendation:** Start with **A + D** (fix workflows + clean repo), then **B** (Terraform), then **C** (ArgoCD).
