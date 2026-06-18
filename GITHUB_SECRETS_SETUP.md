# GitHub Secrets Setup Guide

Required secrets for the optimized CI/CD workflows.

## Required Secrets

### 1. AWS Credentials (Required for ECR)

| Secret Name | Description | How to Get |
|-------------|-------------|------------|
| `AWS_ACCESS_KEY_ID` | AWS IAM access key | Create IAM user with ECR permissions |
| `AWS_SECRET_ACCESS_KEY` | AWS IAM secret key | Get from IAM user creation |
| `AWS_ACCOUNT_ID` | Your AWS account ID | Find in AWS Console (top right) |

### 2. Optional Secrets

| Secret Name | Description | Required For |
|-------------|-------------|--------------|
| `SNYK_TOKEN` | Snyk API token | Enhanced dependency scanning |
| `SLACK_WEBHOOK` | Slack webhook URL | Deployment notifications |

---

## Setup Instructions

### Step 1: Create IAM User for GitHub Actions

```bash
# 1. Go to AWS IAM Console
https://console.aws.amazon.com/iam/

# 2. Create new user
Name: github-actions-ecr
Access type: Programmatic access

# 3. Attach policy (create inline policy)
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ECRFullAccess",
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetRepositoryPolicy",
        "ecr:DescribeRepositories",
        "ecr:ListImages",
        "ecr:DescribeImages",
        "ecr:BatchGetImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
        "ecr:PutImage"
      ],
      "Resource": "*"
    }
  ]
}

# 4. Save credentials
AWS_ACCESS_KEY_ID: AKIA...
AWS_SECRET_ACCESS_KEY: ...
```

### Step 2: Get AWS Account ID

```bash
# Option 1: AWS CLI
aws sts get-caller-identity --query Account --output text

# Option 2: AWS Console
# Top right → Click account name → Account ID shown

# Example: 123456789012
```

### Step 3: Add Secrets to GitHub

```bash
# 1. Go to your repository
https://github.com/yourusername/e-commerce

# 2. Navigate to Settings
Settings → Secrets and variables → Actions

# 3. Click "New repository secret"

# 4. Add each secret:
Name: AWS_ACCESS_KEY_ID
Value: AKIA...

Name: AWS_SECRET_ACCESS_KEY
Value: <paste secret key>

Name: AWS_ACCOUNT_ID
Value: 123456789012
```

---

## Validation

### Test AWS Credentials

```bash
# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Configure credentials
aws configure
# Paste your ACCESS_KEY_ID and SECRET_ACCESS_KEY

# Test ECR access
aws ecr describe-repositories --region us-east-1

# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin \
  <your-account-id>.dkr.ecr.us-east-1.amazonaws.com
```

### Test GitHub Actions

```bash
# Push a small change
git checkout -b test/secrets
echo "# Test" >> README.md
git add .
git commit -m "Test GitHub Actions secrets"
git push origin test/secrets

# Create PR and watch Actions tab
# Should see:
# ✓ CI jobs running
# ✓ Building Docker images
# (Will fail at ECR push if secrets not configured)
```

---

## Optional: Snyk Setup

### Get Snyk Token

```bash
# 1. Sign up at https://snyk.io
# 2. Go to Account Settings
# 3. Click "Generate" under API Token
# 4. Copy token

# 5. Add to GitHub
Name: SNYK_TOKEN
Value: <paste token>
```

### Benefits
- Detailed dependency vulnerability reports
- License compliance checking
- Fix recommendations
- Integration with GitHub Security

---

## Optional: Slack Notifications

### Create Webhook

```bash
# 1. Go to your Slack workspace
# 2. Add "Incoming Webhooks" app
# 3. Choose channel (e.g., #deployments)
# 4. Copy webhook URL

# 5. Add to GitHub
Name: SLACK_WEBHOOK
Value: https://hooks.slack.com/services/...
```

### Notification Format

```
✅ Deployment Successful
Repository: ecommerce
Branch: main
Commit: abc123d
Environment: production
Time: 2 minutes
```

---

## Security Best Practices

### IAM User Permissions

✅ **Do:**
- Use dedicated IAM user for GitHub Actions
- Grant minimum required permissions (ECR only)
- Enable MFA for IAM user (console access)
- Rotate credentials every 90 days
- Use IAM roles instead of keys (if possible with OIDC)

❌ **Don't:**
- Use root account credentials
- Grant AdministratorAccess
- Share credentials across multiple projects
- Commit credentials to Git (even accidentally)

### GitHub Secrets

✅ **Do:**
- Use environment-specific secrets
- Review secret access logs
- Rotate secrets regularly
- Use GitHub environments for prod secrets

❌ **Don't:**
- Print secrets in workflow logs
- Store secrets in code or configs
- Use same secrets for dev and prod

---

## Troubleshooting

### Issue: "Unable to locate credentials"

**Cause:** AWS secrets not configured

**Solution:**
```bash
# Verify secrets exist
Settings → Secrets → Check AWS_ACCESS_KEY_ID exists

# Check workflow file uses correct secret names
- uses: aws-actions/configure-aws-credentials@v4
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}  # Must match exactly
```

### Issue: "no basic auth credentials"

**Cause:** ECR login failed

**Solution:**
```bash
# Check AWS account ID is correct
echo ${{ secrets.AWS_ACCOUNT_ID }}

# Verify ECR registry format
ECR_REGISTRY: ${{ secrets.AWS_ACCOUNT_ID }}.dkr.ecr.us-east-1.amazonaws.com
```

### Issue: "AccessDenied" when pushing to ECR

**Cause:** IAM user lacks permissions

**Solution:**
```bash
# Add ECR push permissions to IAM user
{
  "Effect": "Allow",
  "Action": [
    "ecr:PutImage",
    "ecr:InitiateLayerUpload",
    "ecr:UploadLayerPart",
    "ecr:CompleteLayerUpload"
  ],
  "Resource": "*"
}
```

### Issue: Snyk scan fails

**Cause:** Invalid token or quota exceeded

**Solution:**
```bash
# Check token is valid
https://app.snyk.io/account → Generate new token

# Check if free tier limit reached (100 scans/month)
# Upgrade to paid plan or disable Snyk
```

---

## Alternative: OIDC (Recommended for Production)

Instead of long-lived access keys, use OpenID Connect (OIDC) for better security.

### Setup OIDC

```bash
# 1. Create OIDC provider in AWS
https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services

# 2. Create IAM role with trust policy
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
          "token.actions.githubusercontent.com:sub": "repo:OWNER/REPO:ref:refs/heads/main"
        }
      }
    }
  ]
}

# 3. Update workflow to use OIDC
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v4
  with:
    role-to-assume: arn:aws:iam::ACCOUNT_ID:role/GitHubActionsRole
    aws-region: us-east-1

# 4. Remove access key secrets (no longer needed)
```

### Benefits of OIDC
- ✅ No long-lived credentials
- ✅ Automatic credential rotation
- ✅ Better security
- ✅ Audit trail in CloudTrail
- ✅ Fine-grained permissions per branch

---

## Checklist

Before running workflows:

- [ ] AWS IAM user created with ECR permissions
- [ ] AWS credentials added to GitHub Secrets
  - [ ] AWS_ACCESS_KEY_ID
  - [ ] AWS_SECRET_ACCESS_KEY
  - [ ] AWS_ACCOUNT_ID
- [ ] Credentials tested locally
- [ ] ECR repositories created
- [ ] Optional: Snyk token added
- [ ] Optional: Slack webhook added
- [ ] Secrets visible in Settings → Secrets
- [ ] Test workflow run successful

---

## Quick Reference

### Secrets Location
```
Repository → Settings → Secrets and variables → Actions → Repository secrets
```

### AWS CLI Commands
```bash
# Get account ID
aws sts get-caller-identity --query Account --output text

# List ECR repos
aws ecr describe-repositories

# Login to ECR
aws ecr get-login-password | docker login --username AWS --password-stdin <account>.dkr.ecr.us-east-1.amazonaws.com
```

### Workflow Secret Usage
```yaml
env:
  AWS_REGION: us-east-1
  ECR_REGISTRY: ${{ secrets.AWS_ACCOUNT_ID }}.dkr.ecr.us-east-1.amazonaws.com

steps:
  - uses: aws-actions/configure-aws-credentials@v4
    with:
      aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
      aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
      aws-region: ${{ env.AWS_REGION }}
```

---

**Updated:** 2026-06-17
**Version:** 1.0
