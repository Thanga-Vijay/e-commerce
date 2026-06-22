#!/bin/bash

set -e

echo "🚀 Setting up LocalStack for AWS Services Emulation"

# Check if cluster is running
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ k3d cluster is not running. Run ./k3d-setup.sh first"
    exit 1
fi

# Deploy LocalStack
echo "📦 Deploying LocalStack..."
kubectl apply -f k8s/localstack/

# Wait for LocalStack to be ready
echo "⏳ Waiting for LocalStack to be ready..."
kubectl wait --namespace localstack \
  --for=condition=ready pod \
  --selector=app=localstack \
  --timeout=300s

# Get LocalStack endpoint
LOCALSTACK_ENDPOINT="http://$(kubectl get svc -n localstack localstack -o jsonpath='{.spec.clusterIP}'):4566"

echo "✅ LocalStack deployed successfully!"
echo ""
echo "📌 LocalStack Info:"
echo "   Endpoint: $LOCALSTACK_ENDPOINT"
echo "   Dashboard: http://localhost:4566/_localstack/health"
echo ""
echo "🔧 Initialize AWS Resources:"
echo "   Run: ./scripts/init-localstack.sh"
echo ""
echo "🌐 Access LocalStack from pods:"
echo "   AWS_ENDPOINT_URL: http://localstack.localstack.svc.cluster.local:4566"
