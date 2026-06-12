#!/bin/bash

# Deploy Monitoring Stack to Kubernetes
set -e

NAMESPACE="ecommerce"
BLUE='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

echo -e "${BLUE}E-Commerce Platform - Monitoring Stack Deployment${RESET}"
echo "=================================================="
echo ""

# Check prerequisites
if ! command -v kubectl &> /dev/null; then
    echo -e "${YELLOW}kubectl not found. Please install kubectl first.${RESET}"
    exit 1
fi

echo -e "${GREEN}✓ kubectl installed${RESET}"
echo ""

# 1. Deploy Prometheus
echo -e "${BLUE}1. Deploying Prometheus...${RESET}"
kubectl apply -f prometheus-config.yaml
kubectl apply -f prometheus-rules.yaml
kubectl apply -f prometheus.yaml
echo "Waiting for Prometheus to be ready..."
kubectl wait --for=condition=ready pod -l app=prometheus -n $NAMESPACE --timeout=180s || true
echo ""

# 2. Deploy Alertmanager
echo -e "${BLUE}2. Deploying Alertmanager...${RESET}"
echo -e "${YELLOW}⚠ Update alertmanager.yaml with your SMTP and Slack credentials!${RESET}"
kubectl apply -f alertmanager.yaml
echo ""

# 3. Deploy Grafana
echo -e "${BLUE}3. Deploying Grafana...${RESET}"
kubectl apply -f grafana.yaml
echo "Waiting for Grafana to be ready..."
kubectl wait --for=condition=ready pod -l app=grafana -n $NAMESPACE --timeout=180s || true
echo ""

# 4. Deploy Loki (Log Aggregation)
echo -e "${BLUE}4. Deploying Loki and Promtail...${RESET}"
kubectl apply -f loki.yaml
echo "Waiting for Loki to be ready..."
kubectl wait --for=condition=ready pod -l app=loki -n $NAMESPACE --timeout=180s || true
echo ""

# 5. Deploy Jaeger (Distributed Tracing)
echo -e "${BLUE}5. Deploying Jaeger...${RESET}"
kubectl apply -f jaeger.yaml
echo "Waiting for Jaeger to be ready..."
kubectl wait --for=condition=ready pod -l app=jaeger -n $NAMESPACE --timeout=180s || true
echo ""

# 6. Deploy ServiceMonitors (if using Prometheus Operator)
echo -e "${BLUE}6. Deploying ServiceMonitors (optional)...${RESET}"
kubectl apply -f servicemonitors.yaml 2>/dev/null || echo "Skipping ServiceMonitors (Prometheus Operator not installed)"
echo ""

# Display status
echo -e "${GREEN}✓ Monitoring stack deployment complete!${RESET}"
echo ""
echo -e "${BLUE}Checking deployment status...${RESET}"
kubectl get pods -n $NAMESPACE | grep -E "prometheus|grafana|loki|jaeger|promtail|alertmanager"
echo ""

# Get service URLs
echo -e "${GREEN}Access URLs:${RESET}"
echo ""

GRAFANA_IP=$(kubectl get svc grafana-service -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "pending")
PROMETHEUS_PORT=$(kubectl get svc prometheus-service -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || echo "N/A")
JAEGER_IP=$(kubectl get svc jaeger-query-service -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "pending")

echo "Grafana:    http://$GRAFANA_IP:3000 (admin/admin)"
echo "Prometheus: http://<node-ip>:$PROMETHEUS_PORT or use port-forward"
echo "Jaeger:     http://$JAEGER_IP:16686"
echo ""

echo -e "${YELLOW}Port-forward commands if LoadBalancer IPs are pending:${RESET}"
echo "kubectl port-forward -n $NAMESPACE svc/grafana-service 3000:3000"
echo "kubectl port-forward -n $NAMESPACE svc/prometheus-service 9090:9090"
echo "kubectl port-forward -n $NAMESPACE svc/jaeger-query-service 16686:16686"
echo ""

echo -e "${GREEN}Next steps:${RESET}"
echo "1. Access Grafana and configure dashboards"
echo "2. Update Alertmanager config with your notification channels"
echo "3. Add Prometheus annotations to service pods for metrics scraping"
echo "4. Configure Jaeger client libraries in microservices"
echo "5. Import pre-built Grafana dashboards for Kubernetes, Go, PostgreSQL"
