# Kubernetes-Based Development Environment Setup

## 1. Prerequisites
- Kubernetes cluster
- kubectl configured
- Docker registry access

## 2. Build & Deploy

```bash
# Build and push Docker image
docker build -t your-registry/test-env:latest -f deployments/templates/Dockerfile .
docker push your-registry/test-env:latest

# Deploy to Kubernetes
kubectl apply -f deployments/templates/test-pod.yaml
```

## 3. Testing Workflow

```bash
# Make script executable
chmod +x scripts/test-runner.sh

# Start interactive testing
./scripts/test-runner.sh
```

## 4. Monitoring

```bash
# Check pod status
kubectl get pods -l app=test-pod

# View logs
kubectl logs -f deployment/test-pod
```

## 5. Cleanup

```bash
kubectl delete -f deployments/templates/test-pod.yaml
```
