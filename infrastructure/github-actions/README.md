# CI/CD with GitHub Actions

Automated build, test, and deployment pipelines.

## Workflows

### 1. Build and Test (`build-test.yml`)
Triggered on: Push to main, Pull requests

Steps:
- Checkout code
- Setup Go
- Run unit tests
- Run integration tests
- Code coverage report
- Lint code

### 2. Build Docker Images (`build-images.yml`)
Triggered on: Push to main

Steps:
- Build Docker images for all services
- Tag images with commit SHA and latest
- Push to container registry
- Scan images for vulnerabilities

### 3. Deploy to Development (`deploy-dev.yml`)
Triggered on: Push to main

Steps:
- Build and push images
- Deploy to KIND cluster (dev environment)
- Run smoke tests
- Notify on Slack

### 4. Deploy to Production (`deploy-prod.yml`)
Triggered on: Manual trigger, Release tag

Steps:
- Build and push images with release tag
- Deploy to production Kubernetes cluster
- Run smoke tests
- Rollback on failure
- Notify on Slack

### 5. Database Migrations (`migrations.yml`)
Triggered on: Manual trigger

Steps:
- Run database migrations
- Validate schema
- Rollback on failure

## Secrets Required

GitHub Repository Secrets:
- `DOCKER_USERNAME`
- `DOCKER_PASSWORD`
- `KUBE_CONFIG_DEV`
- `KUBE_CONFIG_PROD`
- `SLACK_WEBHOOK`
- `STRIPE_SECRET_KEY`
- `JWT_SECRET`

## Branch Strategy

- `main` - Main development branch, auto-deploys to dev
- `staging` - Staging environment
- `production` - Production environment
- Feature branches: `feature/*`
- Bugfix branches: `bugfix/*`

## Release Process

1. Create release branch from main
2. Run tests and quality checks
3. Create release tag (e.g., v1.0.0)
4. Trigger production deployment
5. Monitor deployment
6. Merge back to main

## Notifications

Slack notifications sent for:
- Build failures
- Test failures
- Deployment success/failure
- Production deployments

## Example Workflow

```yaml
name: Build and Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run Tests
        run: |
          cd auth-service
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```
