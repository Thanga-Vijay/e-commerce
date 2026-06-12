# Kubernetes Manifests

Kubernetes YAML manifests for deploying all services.

## Structure

```
kubernetes/
├── namespace.yaml
├── configmaps/
├── secrets/
├── deployments/
├── services/
├── ingress/
├── hpa/
└── pvcs/
```

## Deployment Order

1. Namespace
2. Secrets and ConfigMaps
3. PersistentVolumeClaims
4. Infrastructure (PostgreSQL, Redis, Kafka)
5. Application Services
6. Ingress
7. HPA (Horizontal Pod Autoscaler)

## Apply All Manifests

```bash
kubectl apply -f namespace.yaml
kubectl apply -f secrets/
kubectl apply -f configmaps/
kubectl apply -f pvcs/
kubectl apply -f deployments/
kubectl apply -f services/
kubectl apply -f ingress/
kubectl apply -f hpa/
```

## Verify Deployment

```bash
kubectl get all -n ecommerce
kubectl get pods -n ecommerce
kubectl get services -n ecommerce
```

## Common Operations

Scale deployment:
```bash
kubectl scale deployment auth-service -n ecommerce --replicas=3
```

View logs:
```bash
kubectl logs -f deployment/auth-service -n ecommerce
```

Port forward:
```bash
kubectl port-forward svc/api-gateway 8080:8080 -n ecommerce
```
