# Tenet SaaS Backend — Cloud‑Native Deployment on AWS EKS

> Production‑style SaaS backend deployed on Kubernetes (AWS EKS) with Postgres, Redis, autoscaling, and observability foundations.

---

# 🏯 Project Overview

**Tenet** is a multi‑tenant SaaS backend implemented in Go and deployed as a cloud‑native system on AWS using Kubernetes (EKS).

This repository demonstrates:

- Containerized backend service
- Kubernetes deployment (EKS)
- Managed Postgres + Redis inside cluster
- LoadBalancer + Ingress exposure
- Horizontal Pod Autoscaling (HPA)
- Metrics Server
- Terraform‑managed AWS resources (ECR)
- Observability foundation (Grafana)

---

# 🧱 Architecture

```
Client → AWS ALB → K8s Service → Backend Pods
                       ↓
                 Redis / Postgres
```

---

# 📦 Repository Structure

```
saas-platform/
  backend/              Go API server
  infra/
    k8s/                Kubernetes manifests
    terraform/          AWS infrastructure
```

---

# ⚙️ Prerequisites

Install locally:

- Docker
- kubectl
- eksctl
- awscli
- helm
- terraform
- jq

AWS setup:

```bash
aws configure
```

---

# 🐳 Stage 1 — Build Backend Image

From backend folder:

```bash
docker build -t tenet-saas-backend .
```

---

# 📦 Stage 2 — Create AWS ECR Repository (Terraform)

```bash
cd saas-platform/infra/terraform
terraform init
terraform apply
```

Output gives ECR URI.

---

# ☁️ Stage 3 — Push Image to ECR

```bash
ECR_URI=$(aws ecr describe-repositories \
  --repository-names tenet/saas-backend \
  --region us-east-1 \
  --query 'repositories[0].repositoryUri' \
  --output text)

aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin $ECR_URI

docker tag tenet-saas-backend:latest $ECR_URI:latest
docker push $ECR_URI:latest
```

---

# ☸️ Stage 4 — Create EKS Cluster

```bash
eksctl create cluster \
  --name tenet-cluster \
  --region us-east-1 \
  --nodes 1
```

Verify:

```bash
kubectl get nodes
```

---

# 🗄️ Stage 5 — Deploy Databases

### Redis

```bash
kubectl apply -f redis.yaml
```

### Postgres

```bash
kubectl apply -f postgres.yaml
```

Create DB schema:

```bash
POD=$(kubectl get pod -l app=postgres -o jsonpath='{.items[0].metadata.name}')
kubectl cp ../../backend/schema.sql $POD:/schema.sql
kubectl exec -it $POD -- psql -U postgres -d saas -f /schema.sql
```

---

# 🚀 Stage 6 — Deploy Backend

```bash
kubectl apply -f backend-deployment.yaml
kubectl apply -f backend-service.yaml
```

Get external URL:

```bash
kubectl get svc saas-backend
```

Example:

```
http://<EXTERNAL-IP>
```

---

# 🔑 API Test

```bash
SERVICE_URL=http://<EXTERNAL-IP>

TENANT_ID=$(curl -s -X POST $SERVICE_URL/tenants \
  -H "Content-Type: application/json" \
  -d '{"name":"AWSCo"}' | jq -r '.id')

API_KEY=$(curl -s -X POST $SERVICE_URL/tenants/$TENANT_ID/keys | jq -r '.api_key')

curl $SERVICE_URL/protected -H "X-API-Key: $API_KEY"
```

---

# 📈 Stage 7A — Metrics Server (Autoscaling prerequisite)

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

kubectl -n kube-system patch deploy metrics-server --type='json' -p='[
  {"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"},
  {"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname"}
]'

kubectl -n kube-system rollout status deploy/metrics-server
```

Verify:

```bash
kubectl top nodes
```

---

# 📊 Stage 7B — Horizontal Pod Autoscaler

```bash
kubectl apply -f hpa.yaml
kubectl get hpa
```

---

# 📊 Stage 7C — Grafana (Lightweight Observability)

Add repo:

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

Install:

```bash
helm install tenet-grafana grafana/grafana \
  -n monitoring \
  --create-namespace \
  --set resources.requests.memory=128Mi
```

Get password:

```bash
kubectl -n monitoring get secret tenet-grafana \
  -o jsonpath="{.data.admin-password}" | base64 --decode
```

Port‑forward:

```bash
kubectl -n monitoring port-forward svc/tenet-grafana 3000:80
```

Open:

```
http://localhost:3000
```

---

# 📈 Current System Capabilities

- Multi‑tenant SaaS backend
- API key auth
- Postgres persistence
- Redis caching
- Kubernetes deployment
- External LoadBalancer
- Autoscaling via HPA
- Metrics Server
- Grafana dashboards (ready)
- Terraform ECR provisioning
- AWS EKS cluster

---

# 🧪 Example Endpoints

```
POST /tenants
POST /tenants/{id}/keys
GET  /protected
```

---

# 🔭 Future Scope

## Infrastructure

- ALB Ingress Controller (fully IAM‑wired)
- TLS via ACM + HTTPS
- Multi‑AZ nodegroups
- Cluster autoscaler

## Observability

- Prometheus Operator (full stack)
- ServiceMonitor for backend
- RED metrics dashboards
- Alertmanager rules

## Backend

- Rate limiting (Redis token bucket)
- Tenant quotas
- Usage metering
- Billing integration

## DevOps

- GitHub Actions CI/CD
- Helm charts for Tenet
- Blue‑Green deployment
- Canary rollout

---

# 🏁 Resume Value

This project demonstrates production‑style cloud backend engineering:

- Kubernetes orchestration
- AWS infrastructure
- Stateful services
- Autoscaling
- Observability
- Infrastructure‑as‑Code

Suitable for:

- Backend Engineer roles
- Cloud / Platform roles
- Distributed Systems roles


