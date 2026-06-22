#!/bin/bash

set -e

echo "🚀 Setting up Floci for Cloud Services Emulation (AWS/GCP/Azure)"

# Check if cluster is running
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ k3d cluster is not running. Run ./k3d-setup.sh first"
    exit 1
fi

# Create Floci namespace if it doesn't exist
kubectl create namespace floci --dry-run=client -o yaml | kubectl apply -f -

# Deploy Floci
echo "📦 Deploying Floci..."
kubectl apply -f k8s/floci/

# Wait for Floci to be ready
echo "⏳ Waiting for Floci to be ready..."
kubectl wait --namespace floci \
  --for=condition=ready pod \
  --selector=app=floci \
  --timeout=300s

# Get Floci endpoint
FLOCI_ENDPOINT="http://$(kubectl get svc -n floci floci -o jsonpath='{.spec.clusterIP}'):4566"

echo "✅ Floci deployed successfully!"
echo ""
echo "📌 Floci Info:"
echo "   Endpoint: $FLOCI_ENDPOINT"
echo "   AWS Compatible Endpoint: http://localhost:4566"
echo "   Health Check: http://localhost:4566/health"
echo ""
echo "🔧 Services Available:"
echo "   - S3 (Object Storage)"
echo "   - SQS (Message Queue)"
echo "   - SNS (Notifications)"
echo "   - DynamoDB (NoSQL Database)"
echo "   - Lambda (Serverless Functions)"
echo "   - Secrets Manager"
echo "   - And more..."
echo ""
echo "🌐 Access Floci from pods:"
echo "   Set environment variable: AWS_ENDPOINT_URL=http://floci.floci.svc.cluster.local:4566"
echo ""
echo "📝 Configure AWS CLI:"
echo "   aws configure set aws_access_key_id test"
echo "   aws configure set aws_secret_access_key test"
echo "   aws configure set region us-east-1"
echo "   aws --endpoint-url=http://localhost:4566 s3 ls"
