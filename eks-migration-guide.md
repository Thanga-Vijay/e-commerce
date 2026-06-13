# AWS EKS Migration Guide

Complete guide for migrating from KIND local testing to AWS EKS production deployment.

## Prerequisites

### AWS CLI & Tools
```bash
# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Install eksctl
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# Install kubectl (if not already installed)
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### Configure AWS Credentials
```bash
aws configure
# Enter:
# - AWS Access Key ID
# - AWS Secret Access Key
# - Default region: us-east-1
# - Default output format: json

# Verify
aws sts get-caller-identity
```

## Step 1: Create EKS Cluster

### Option A: Using eksctl (Recommended)

**Create cluster configuration file:**

```bash
cat > eks-cluster.yaml <<EOF
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: ecommerce-prod
  region: us-east-1
  version: "1.28"

# VPC Configuration
vpc:
  cidr: 10.0.0.0/16
  nat:
    gateway: HighlyAvailable

# IAM Configuration
iam:
  withOIDC: true

# Managed Node Groups
managedNodeGroups:
  # Application nodes
  - name: app-nodes
    instanceType: t3.large
    desiredCapacity: 3
    minSize: 2
    maxSize: 6
    volumeSize: 50
    labels:
      tier: application
      node-type: compute
    tags:
      Name: ecommerce-app-node
      Environment: production
    iam:
      withAddonPolicies:
        autoScaler: true
        cloudWatch: true
        albIngress: true

  # Database nodes (for StatefulSets)
  - name: db-nodes
    instanceType: t3.xlarge
    desiredCapacity: 3
    minSize: 3
    maxSize: 5
    volumeSize: 100
    labels:
      tier: database
      node-type: storage
    tags:
      Name: ecommerce-db-node
      Environment: production
    iam:
      withAddonPolicies:
        ebs: true

# Add-ons
addons:
  - name: vpc-cni
    version: latest
  - name: coredns
    version: latest
  - name: kube-proxy
    version: latest
  - name: aws-ebs-csi-driver
    version: latest

# CloudWatch Logging
cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
EOF

# Create cluster (takes 15-20 minutes)
eksctl create cluster -f eks-cluster.yaml
```

### Option B: Using AWS Console
1. Go to EKS Console
2. Create cluster
3. Configure networking, IAM, add-ons
4. Create node groups

### Verify Cluster
```bash
# Update kubeconfig
aws eks update-kubeconfig --region us-east-1 --name ecommerce-prod

# Verify connection
kubectl cluster-info
kubectl get nodes
```

## Step 2: Set Up Container Registry (ECR)

### Create ECR Repositories
```bash
# Create repositories for each service
SERVICES=(
  "auth-service"
  "product-service"
  "cart-service"
  "wishlist-service"
  "order-service"
  "payment-service"
  "inventory-service"
  "notification-service"
  "reporting-service"
  "frontend"
)

AWS_REGION="us-east-1"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

for service in "${SERVICES[@]}"; do
  echo "Creating ECR repository for ${service}..."
  aws ecr create-repository \
    --repository-name ecommerce/${service} \
    --region ${AWS_REGION} \
    --image-scanning-configuration scanOnPush=true \
    --encryption-configuration encryptionType=AES256 || echo "Repository already exists"
done
```

### Build and Push Images
```bash
# Login to ECR
aws ecr get-login-password --region ${AWS_REGION} | \
  docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

# Build and push each service
VERSION="v1.0.0"

for service in "${SERVICES[@]}"; do
  echo "Building and pushing ${service}..."
  
  # Build image
  if [ "$service" = "frontend" ]; then
    docker build -t ecommerce/${service}:${VERSION} frontend/
  else
    docker build -t ecommerce/${service}:${VERSION} services/${service}/
  fi
  
  # Tag for ECR
  docker tag ecommerce/${service}:${VERSION} \
    ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/ecommerce/${service}:${VERSION}
  
  # Push to ECR
  docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/ecommerce/${service}:${VERSION}
  
  echo "✓ ${service} pushed to ECR"
done
```

## Step 3: Install Add-ons

### AWS Load Balancer Controller (for Ingress)
```bash
# Create IAM policy
curl -o iam_policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.6.0/docs/install/iam_policy.json

aws iam create-policy \
  --policy-name AWSLoadBalancerControllerIAMPolicy \
  --policy-document file://iam_policy.json

# Create service account
eksctl create iamserviceaccount \
  --cluster=ecommerce-prod \
  --namespace=kube-system \
  --name=aws-load-balancer-controller \
  --attach-policy-arn=arn:aws:iam::${AWS_ACCOUNT_ID}:policy/AWSLoadBalancerControllerIAMPolicy \
  --approve

# Install controller using Helm
helm repo add eks https://aws.github.io/eks-charts
helm repo update

helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=ecommerce-prod \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller

# Verify
kubectl get deployment -n kube-system aws-load-balancer-controller
```

### EBS CSI Driver (for Persistent Volumes)
```bash
# Create IAM role
eksctl create iamserviceaccount \
  --name ebs-csi-controller-sa \
  --namespace kube-system \
  --cluster ecommerce-prod \
  --attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy \
  --approve \
  --role-only \
  --role-name AmazonEKS_EBS_CSI_DriverRole

# Install add-on
eksctl create addon \
  --name aws-ebs-csi-driver \
  --cluster ecommerce-prod \
  --service-account-role-arn arn:aws:iam::${AWS_ACCOUNT_ID}:role/AmazonEKS_EBS_CSI_DriverRole \
  --force

# Verify
kubectl get pods -n kube-system | grep ebs-csi
```

### External DNS (Optional)
```bash
# For automatic DNS management with Route53
helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/

helm install external-dns external-dns/external-dns \
  --namespace kube-system \
  --set provider=aws \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=ecommerce-prod
```

### Cert-Manager (for TLS)
```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wait for cert-manager to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager --timeout=90s

# Create ClusterIssuer for Let's Encrypt
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@yourdomain.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: alb
EOF
```

## Step 4: Update Kubernetes Manifests for EKS

### Update Image References
```bash
# Update all deployment files to use ECR images
cd k8s/services

for file in *.yaml; do
  sed -i "s|image: ecommerce/|image: ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/ecommerce/|g" $file
  sed -i "s|:latest|:v1.0.0|g" $file
done

cd ../frontend
sed -i "s|image: ecommerce/|image: ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/ecommerce/|g" frontend.yaml
sed -i "s|:latest|:v1.0.0|g" frontend.yaml
```

### Update Ingress for AWS ALB
```bash
# Create ALB Ingress
cat > k8s/ingress/ingress-alb.yaml <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ecommerce-ingress
  namespace: ecommerce
  annotations:
    # AWS Load Balancer Controller
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    
    # SSL/TLS
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS": 443}]'
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:REGION:ACCOUNT:certificate/CERT_ID
    
    # Health check
    alb.ingress.kubernetes.io/healthcheck-path: /health
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    
    # Other settings
    alb.ingress.kubernetes.io/load-balancer-name: ecommerce-alb
    alb.ingress.kubernetes.io/group.name: ecommerce
spec:
  rules:
  - host: ecommerce.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
  
  - host: api.ecommerce.yourdomain.com
    http:
      paths:
      - path: /api/v1/auth
        pathType: Prefix
        backend:
          service:
            name: auth-service
            port:
              number: 8081
      # Add other service paths...
EOF
```

### Update StorageClass for EBS
```bash
cat > k8s/storage-class.yaml <<EOF
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: ebs-gp3
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: ebs.csi.aws.com
parameters:
  type: gp3
  iops: "3000"
  throughput: "125"
  encrypted: "true"
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
reclaimPolicy: Delete
EOF

kubectl apply -f k8s/storage-class.yaml
```

## Step 5: Deploy to EKS

### Deploy Application
```bash
# Deploy namespace and configs
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# Deploy databases
kubectl apply -f k8s/databases/

# Deploy Redis and Kafka
kubectl apply -f k8s/redis/
kubectl apply -f k8s/kafka/

# Wait for infrastructure
kubectl wait --for=condition=ready pod -l tier=database -n ecommerce --timeout=10m

# Deploy microservices
kubectl apply -f k8s/services/

# Deploy frontend
kubectl apply -f k8s/frontend/

# Deploy ingress
kubectl apply -f k8s/ingress/ingress-alb.yaml

# Deploy HPA
kubectl apply -f k8s/hpa/
```

### Deploy Monitoring
```bash
cd k8s/monitoring
./deploy-monitoring.sh
```

### Deploy Security
```bash
kubectl apply -f k8s/security/
```

## Step 6: Configure DNS

### Get ALB DNS Name
```bash
ALB_DNS=$(kubectl get ingress -n ecommerce ecommerce-ingress -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
echo "ALB DNS: $ALB_DNS"
```

### Update Route53
```bash
# Create Route53 hosted zone (if not exists)
aws route53 create-hosted-zone --name yourdomain.com --caller-reference $(date +%s)

# Create A record pointing to ALB
# (Use AWS Console or CLI)
```

## Step 7: Set Up Monitoring & Logging

### CloudWatch Container Insights
```bash
# Install Container Insights
curl https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluentd-quickstart.yaml | \
sed "s/{{cluster_name}}/ecommerce-prod/;s/{{region_name}}/${AWS_REGION}/" | \
kubectl apply -f -
```

### AWS X-Ray (for tracing)
```bash
# Install X-Ray DaemonSet
kubectl apply -f https://eksworkshop.com/intermediate/245_x-ray/daemonset.files/xray-k8s-daemonset.yaml
```

## Step 8: Set Up Backup

### Velero for K8s Backup
```bash
# Install Velero CLI
wget https://github.com/vmware-tanzu/velero/releases/download/v1.12.0/velero-v1.12.0-linux-amd64.tar.gz
tar -xvf velero-v1.12.0-linux-amd64.tar.gz
sudo mv velero-v1.12.0-linux-amd64/velero /usr/local/bin/

# Install Velero in cluster
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws:v1.8.0 \
  --bucket ecommerce-velero-backups \
  --backup-location-config region=${AWS_REGION} \
  --snapshot-location-config region=${AWS_REGION} \
  --secret-file ./credentials-velero

# Create backup schedule
velero schedule create daily-backup --schedule="0 2 * * *"
```

## Cost Optimization

### Use Spot Instances
```bash
# Add spot instance node group
eksctl create nodegroup \
  --cluster=ecommerce-prod \
  --name=spot-nodes \
  --node-type=t3.large \
  --nodes=2 \
  --nodes-min=1 \
  --nodes-max=5 \
  --spot
```

### Set Up Cluster Autoscaler
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/autoscaler/master/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml
```

## Differences: KIND vs EKS

| Feature | KIND | EKS |
|---------|------|-----|
| LoadBalancer | NodePort/HostPort | AWS ALB/NLB |
| Storage | hostPath | EBS (gp3) |
| Ingress | NGINX | AWS ALB Controller |
| DNS | localhost | Route53 |
| TLS | Self-signed | ACM/Let's Encrypt |
| Monitoring | Prometheus | CloudWatch + Prometheus |
| Cost | Free | ~$200-500/month |
| Scalability | Limited | Auto-scaling |

## Testing Checklist

- [ ] All pods running
- [ ] Services accessible via ALB
- [ ] Database connections working
- [ ] Kafka topics created
- [ ] Monitoring dashboards available
- [ ] Logs flowing to CloudWatch
- [ ] Auto-scaling working
- [ ] Backups configured
- [ ] SSL/TLS working
- [ ] DNS resolving correctly

## Rollback Procedure

```bash
# If issues occur, rollback deployment
kubectl rollout undo deployment/<service-name> -n ecommerce

# Or restore from backup
velero restore create --from-backup daily-backup-<date>
```

## Support & Troubleshooting

```bash
# Check pod status
kubectl get pods -n ecommerce

# View logs
kubectl logs -f <pod-name> -n ecommerce

# Describe resources
kubectl describe pod <pod-name> -n ecommerce

# Check ALB
kubectl describe ingress -n ecommerce

# View EKS cluster info
eksctl get cluster ecommerce-prod
```

---

**Estimated Migration Time:** 2-3 hours  
**Estimated Monthly Cost:** $200-500 (depending on usage)  
**Production Readiness:** Yes ✅
