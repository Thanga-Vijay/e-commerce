# E-Commerce Platform - Phase 12: CI/CD Pipeline

Complete CI/CD implementation with GitHub Actions for automated testing, building, and deployment.

## Overview

Phase 12 implements production-grade CI/CD pipelines:
- **Continuous Integration**: Automated testing, linting, security scanning
- **Continuous Deployment**: Automated builds, container registry, multi-environment deployment
- **Quality Gates**: Code quality, test coverage, vulnerability scanning
- **Deployment Strategies**: Staging → Production with approval gates

## Architecture

```
┌─────────────────────────────────────────────────┐
│           GitHub Repository (Main)              │
└────────────┬────────────────────────────────────┘
             │
      ┌──────┴───────┐
      │              │
┌─────▼─────┐  ┌────▼──────┐
│PR/Push to │  │Push to    │
│develop    │  │main/tag   │
└─────┬─────┘  └────┬──────┘
      │             │
┌─────▼─────────────▼──────────────────┐
│       GitHub Actions Workflows       │
├──────────────────────────────────────┤
│  CI Pipeline           CD Pipeline    │
│  • Lint                • Build        │
│  • Test                • Push Images  │
│  • Security Scan       • Deploy       │
│  • Build Check         • Rollout      │
└──────┬───────────────────┬───────────┘
       │                   │
       │         ┌─────────┴────────┐
       │         │                  │
       │    ┌────▼──────┐     ┌────▼─────┐
       │    │ Staging   │     │Production│
       │    │Environment│     │Environment│
       │    └───────────┘     └──────────┘
       │
       └─> Notifications (Slack, Email)
```

## GitHub Actions Workflows

### 1. CI Pipeline (.github/workflows/ci.yml)

**Triggers:**
- Pull requests to main/develop
- Pushes to develop branch

**Jobs:**

#### a) Lint Backend
- Runs golangci-lint on all Go services
- Checks code style and common issues
- Enforces coding standards

#### b) Test Services
- Runs unit tests for all 9 microservices
- Generates code coverage reports
- Uploads coverage to Codecov
- Matrix strategy for parallel execution

#### c) Build Services
- Builds Docker images for all services
- Validates Dockerfiles
- Uses layer caching for faster builds
- Does not push (validation only)

#### d) Test Frontend
- Installs npm dependencies
- Runs ESLint
- Executes Jest tests
- Builds production bundle
- Generates coverage report

#### e) Security Scan
- Runs Trivy vulnerability scanner
- Scans for CVEs in dependencies
- Uploads results to GitHub Security tab
- Fails on high severity issues

#### f) Validate Kubernetes
- Validates YAML syntax
- Dry-run applies manifests
- Runs kubeval for schema validation
- Ensures deployability

### 2. CD Pipeline (.github/workflows/cd.yml)

**Triggers:**
- Pushes to main branch
- Git tags (v*.*.*)

**Jobs:**

#### a) Build and Push
- Builds Docker images for all services
- Tags with version, branch, SHA
- Pushes to GitHub Container Registry (ghcr.io)
- Matrix strategy for 10 services (9 backend + frontend)
- Automatic image tagging:
  - `latest` for main branch
  - `v1.2.3` for tags
  - `main-abc123` for commits

#### b) Deploy to Staging
- Auto-deploys on main branch
- Updates Kubernetes deployments
- Waits for rollout completion
- Runs smoke tests
- Environment: staging

#### c) Deploy to Production
- Requires manual approval
- Only for tagged releases (v*.*.*)
- Updates production cluster
- Health checks after deployment
- Sends Slack notification
- Environment: production

#### d) Rollback
- Automatic on deployment failure
- Reverts to previous version
- Triggered when deploy-production fails

### 3. Security Scanning (.github/workflows/security.yml)

**Triggers:**
- Daily schedule (2 AM)
- Pushes to main/develop
- Pull requests

**Jobs:**

#### a) Trivy Container Scan
- Scans all Docker images
- Detects OS and application vulnerabilities
- Uploads SARIF to GitHub Security
- Matrix for all 9 services

#### b) Dependency Scan
- Runs Snyk for Go dependencies
- Runs npm audit for frontend
- Fails on high severity issues
- Checks for outdated packages

#### c) Secret Scanning
- Uses Gitleaks to detect leaked secrets
- Scans commit history
- Prevents credential exposure

#### d) CodeQL Analysis
- Static analysis for Go and JavaScript
- Detects security issues and bugs
- Integrated with GitHub Security

## Setup Instructions

### 1. Configure GitHub Secrets

Add these secrets in GitHub repository settings:

```bash
# Container Registry
GITHUB_TOKEN  # Auto-provided by GitHub

# Kubernetes (Staging)
KUBECONFIG_STAGING  # Base64 encoded kubeconfig

# Kubernetes (Production)
KUBECONFIG_PRODUCTION  # Base64 encoded kubeconfig

# Security Scanning
SNYK_TOKEN  # From snyk.io

# Notifications
SLACK_WEBHOOK  # Slack incoming webhook URL

# AWS (optional, for S3 backups)
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
```

### 2. Encode Kubeconfig

```bash
# Get kubeconfig for staging cluster
kubectl config view --minify --flatten > staging-kubeconfig.yaml

# Base64 encode
cat staging-kubeconfig.yaml | base64 -w 0

# Add to GitHub Secrets as KUBECONFIG_STAGING
```

### 3. Enable GitHub Container Registry

1. Go to repository **Settings** → **Actions** → **General**
2. Under "Workflow permissions", select "Read and write permissions"
3. Check "Allow GitHub Actions to create and approve pull requests"

### 4. Configure Branch Protection

**For main branch:**
- Require pull request reviews (1 approver)
- Require status checks (CI pipeline)
- Require branches to be up to date
- Include administrators

**For develop branch:**
- Require status checks (CI pipeline)
- Allow force pushes for rebase

### 5. Set Up Environments

Create environments in GitHub:

**Staging:**
- No approval required
- Auto-deploy on main branch

**Production:**
- Require reviewers (2 approvers)
- Only for tagged releases
- Wait timer: 0 minutes (immediate after approval)

## Usage

### Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/new-feature

# 2. Make changes
# ... code changes ...

# 3. Commit and push
git add .
git commit -m "feat: add new feature"
git push origin feature/new-feature

# 4. Create pull request
# CI pipeline runs automatically

# 5. After approval, merge to develop
# CI pipeline runs on develop

# 6. Merge develop to main
# CD pipeline deploys to staging

# 7. Create release tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# CD pipeline deploys to production (after approval)
```

### Deployment Commands

```bash
# Manual deployment trigger (if needed)
gh workflow run cd.yml

# Check workflow status
gh workflow list
gh run list --workflow=cd.yml

# View logs
gh run view <run-id> --log

# Cancel running workflow
gh run cancel <run-id>
```

### Image Tags

Images are automatically tagged:
- `latest` - Main branch
- `v1.2.3` - Git tag
- `main-abc123` - Commit SHA
- `develop` - Develop branch

**Pull images:**
```bash
docker pull ghcr.io/your-org/ecommerce/auth-service:latest
docker pull ghcr.io/your-org/ecommerce/auth-service:v1.2.3
docker pull ghcr.io/your-org/ecommerce/frontend-service:latest
```

## CI/CD Metrics

### Build Times
- Lint: ~2 minutes
- Test (per service): ~3 minutes
- Build (per service): ~5 minutes
- Total CI time: ~15 minutes (parallel)
- Total CD time: ~20 minutes (sequential)

### Success Rates
- Target CI success rate: >95%
- Target CD success rate: >98%
- Rollback rate: <2%

## Testing Strategies

### Unit Tests
```go
// services/auth/handlers/auth_test.go
func TestRegisterUser(t *testing.T) {
    // Test implementation
}
```

### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./...
```

### End-to-End Tests
```bash
# Run E2E tests in staging
kubectl exec -it test-runner -- npm run test:e2e
```

## Quality Gates

Pipeline fails if:
- ✗ Lint errors present
- ✗ Any unit test fails
- ✗ Code coverage drops below 70%
- ✗ High severity vulnerabilities found
- ✗ Docker build fails
- ✗ Kubernetes manifests invalid
- ✗ Deployment rollout fails
- ✗ Health checks fail

## Monitoring CI/CD

### GitHub Actions Dashboard
- View workflow runs in **Actions** tab
- Filter by workflow, branch, status
- Download logs and artifacts

### Notifications
**Slack notifications for:**
- Production deployments
- Deployment failures
- Security vulnerabilities found
- Long-running workflows

**Email notifications for:**
- Workflow failures
- Approval requests
- Deployment completions

## Best Practices

### 1. Commit Messages
Follow conventional commits:
```
feat: add new feature
fix: resolve bug
docs: update documentation
test: add tests
refactor: refactor code
chore: maintenance
```

### 2. Branch Strategy
- `main` - Production-ready code
- `develop` - Integration branch
- `feature/*` - Feature branches
- `hotfix/*` - Urgent fixes

### 3. Versioning
Follow semantic versioning (semver):
- `v1.0.0` - Major release
- `v1.1.0` - Minor release (features)
- `v1.1.1` - Patch release (fixes)

### 4. Rollback Procedure
```bash
# Automatic rollback on failure
# Or manual rollback:
kubectl rollout undo deployment/auth-service -n ecommerce
kubectl rollout undo deployment/product-service -n ecommerce
```

### 5. Cache Optimization
- Use Docker layer caching
- Cache Go modules between runs
- Cache npm dependencies
- Reduces build time by 50%

## Troubleshooting

### Build Failures
```bash
# View logs
gh run view <run-id> --log

# Common issues:
# - Missing dependencies: Update go.mod
# - Test failures: Fix failing tests
# - Docker build: Check Dockerfile
```

### Deployment Failures
```bash
# Check pod status
kubectl get pods -n ecommerce

# View deployment events
kubectl describe deployment auth-service -n ecommerce

# Check rollout status
kubectl rollout status deployment/auth-service -n ecommerce
```

### Security Scan Failures
```bash
# View Trivy report
gh run view <run-id> --log | grep "Total:"

# Update dependencies
go get -u
npm update

# Rebuild images
docker build --no-cache
```

## Cost Optimization

### GitHub Actions Minutes
- Free tier: 2,000 minutes/month
- Optimize with caching
- Use self-hosted runners for heavy workloads

### Container Registry Storage
- Delete old image tags
- Implement retention policy
- Use multi-stage builds to reduce size

## Security Considerations

1. **Secret Management**
   - Never commit secrets
   - Use GitHub Secrets
   - Rotate secrets regularly

2. **Image Scanning**
   - Scan before deployment
   - Update base images regularly
   - Monitor CVE databases

3. **Access Control**
   - Limit who can approve deployments
   - Audit workflow changes
   - Use branch protection

## Summary

Phase 12 delivers enterprise-grade CI/CD:
- ✅ **3 GitHub Actions workflows** - CI, CD, Security
- ✅ **Automated testing** - Unit, integration, E2E
- ✅ **Security scanning** - Trivy, Snyk, CodeQL, Gitleaks
- ✅ **Multi-environment deployment** - Staging + Production
- ✅ **Container registry** - GitHub Container Registry
- ✅ **Approval gates** - Manual approval for production
- ✅ **Automatic rollback** - On deployment failure
- ✅ **Slack notifications** - Deployment status updates
- ✅ **Matrix builds** - Parallel execution for speed
- ✅ **Quality gates** - Coverage, linting, security

Your platform now has full CI/CD automation! 🚀
