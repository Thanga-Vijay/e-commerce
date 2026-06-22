#!/bin/bash

set -e

echo "🔧 Initializing LocalStack AWS Resources"

# LocalStack endpoint
ENDPOINT="http://localhost:30456"
AWS_REGION="us-east-1"

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=$AWS_REGION

# Wait for LocalStack to be ready
echo "⏳ Waiting for LocalStack..."
until curl -s "$ENDPOINT/_localstack/health" | grep -q "running"; do
    echo "Waiting for LocalStack to be ready..."
    sleep 2
done

echo "✅ LocalStack is ready!"

# Create S3 buckets
echo "📦 Creating S3 buckets..."
aws --endpoint-url=$ENDPOINT s3 mb s3://ecommerce-images || true
aws --endpoint-url=$ENDPOINT s3 mb s3://ecommerce-backups || true
aws --endpoint-url=$ENDPOINT s3 mb s3://ecommerce-logs || true

# Set bucket policies (public read for images)
aws --endpoint-url=$ENDPOINT s3api put-bucket-policy \
  --bucket ecommerce-images \
  --policy '{
    "Version": "2012-10-17",
    "Statement": [{
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::ecommerce-images/*"
    }]
  }'

# Create ECR repositories
echo "🐳 Creating ECR repositories..."
for service in auth-service product-service cart-service wishlist-service order-service payment-service inventory-service notification-service reporting-service frontend; do
    aws --endpoint-url=$ENDPOINT ecr create-repository --repository-name ecommerce/$service || true
done

# Create Secrets Manager secrets
echo "🔐 Creating secrets..."
aws --endpoint-url=$ENDPOINT secretsmanager create-secret \
  --name ecommerce/database \
  --secret-string '{
    "username": "postgres",
    "password": "postgres123",
    "host": "postgres.ecommerce.svc.cluster.local",
    "port": "5432"
  }' || true

aws --endpoint-url=$ENDPOINT secretsmanager create-secret \
  --name ecommerce/jwt \
  --secret-string '{
    "secret": "your-super-secret-jwt-key-change-in-production",
    "expiry": "24h"
  }' || true

aws --endpoint-url=$ENDPOINT secretsmanager create-secret \
  --name ecommerce/kafka \
  --secret-string '{
    "brokers": "kafka.ecommerce.svc.cluster.local:9092"
  }' || true

# Create SSM parameters
echo "📝 Creating SSM parameters..."
aws --endpoint-url=$ENDPOINT ssm put-parameter \
  --name /ecommerce/config/log-level \
  --value "info" \
  --type String || true

aws --endpoint-url=$ENDPOINT ssm put-parameter \
  --name /ecommerce/config/feature-flags \
  --value '{"new_checkout": true, "recommendations": true}' \
  --type String || true

# Create CloudWatch Log Groups
echo "📊 Creating CloudWatch log groups..."
for service in auth product cart wishlist order payment inventory notification reporting frontend; do
    aws --endpoint-url=$ENDPOINT logs create-log-group \
      --log-group-name /ecommerce/$service || true
done

# Create IAM role for services
echo "👤 Creating IAM roles..."
aws --endpoint-url=$ENDPOINT iam create-role \
  --role-name ecommerce-service-role \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Principal": {"Service": "ec2.amazonaws.com"},
      "Action": "sts:AssumeRole"
    }]
  }' || true

# Attach policies
aws --endpoint-url=$ENDPOINT iam attach-role-policy \
  --role-name ecommerce-service-role \
  --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess || true

echo ""
echo "✅ LocalStack initialization complete!"
echo ""
echo "📌 Created Resources:"
echo "   S3 Buckets:"
echo "     - ecommerce-images"
echo "     - ecommerce-backups"
echo "     - ecommerce-logs"
echo ""
echo "   ECR Repositories: 9 service repositories + frontend"
echo "   Secrets: database, jwt, kafka"
echo "   SSM Parameters: log-level, feature-flags"
echo "   CloudWatch Log Groups: 10 groups"
echo ""
echo "🔍 Verify resources:"
echo "   aws --endpoint-url=$ENDPOINT s3 ls"
echo "   aws --endpoint-url=$ENDPOINT ecr describe-repositories"
echo "   aws --endpoint-url=$ENDPOINT secretsmanager list-secrets"
