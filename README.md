# BackendIM Brain: Kubernetes-GitOps WebSocket Gateway

![Architecture](https://img.shields.io/badge/Architecture-Microservices-blue)
![WebSocket](https://img.shields.io/badge/Protocol-WebSocket-lightgrey)
![EKS](https://img.shields.io/badge/Cloud-AWS-orange)

A production-grade system for managing Git repositories and Kubernetes clusters through secure WebSocket connections.

## 🌐 System Overview

### Core Components

1. **WebSocket Server** (Go/Gin) - Handles real-time client connections
2. **Git Service** - Manages repository operations via `git_handler.sh`
3. **Kubernetes Integrator** - Handles cluster auth via `kube-init.sh`
4. **Traefik Proxy** - Manages TLS and routing
5. **AWS CLI Runtime** - Isolated execution environment

```mermaid
sequenceDiagram
    participant C as Client
    participant T as Traefik
    participant B as BackendIM-Brain
    participant K as Kubernetes
    participant G as Git
    
    C->>T: wss://DOMAIN/ws
    T->>B: Forward request
    B->>G: Clone/Fetch/Checkout
    B->>K: Apply configurations
    K-->>B: Cluster status
    B-->>C: JSON response
```

## 🚀 Getting Started

### Prerequisites

- AWS Account with EKS cluster
- Domain name with DNS pointing to your server
- Docker 20.10+ and docker-compose 1.29+

### Installation

```bash
git clone https://github.com/obiMadu/backend.im-infra.git
cd backendim-brain
cp .env.example .env
nano .env  # Configure credentials
docker compose build --no-cache
docker compose up -d
```

## 🔧 Configuration

### Environment Variables (.env)

| Variable                  | Description                                | Example                     |
|---------------------------|--------------------------------------------|-----------------------------|
| `AWS_ACCESS_KEY_ID`       | IAM Access Key                             | AKIAXXXXXXXXXXXXXXXX        |
| `AWS_SECRET_ACCESS_KEY`   | IAM Secret Key                             | XXXXXXXXXXXXXXXXXXXXXXXXXX  |
| `AWS_DEFAULT_REGION`      | AWS Region                                 | us-west-2                   |
| `KUBE_CLUSTER_NAME`       | EKS Cluster Name                           | production-cluster          |
| `DOMAIN`                  | Public FQDN                                | gitops.example.com          |
| `ACME_EMAIL`              | Let's Encrypt Email                        | <admin@example.com>           |

### File Structure

```
.
├── cmd/                  # Go application entrypoints
├── internal/             # Core application logic
│   ├── handlers/         # WebSocket connection handlers
│   ├── models/          # Data structures
│   └── services/        # Business logic services
├── scripts/             # Operational scripts
│   ├── git_handler.sh   # Git repository management
│   ├── kube-init.sh     # Kubernetes cluster auth
│   └── install-*.sh     # Dependency installers
└── docker-compose.yml   # Service orchestration
```

## 📡 API Specification

### WebSocket Endpoint

`wss://<DOMAIN>/ws`

### Message Format

```json
{
  "userId": "unique-user-id",
  "chatId": "session-identifier",
  "repoURL": "https://github.com/user/repo.git",
  "commitHash": "a1b2c3d4e5"
}
```

### Response Types

| Type       | Format                                  | Example Use Case              |
|------------|-----------------------------------------|--------------------------------|
| Success    | `{type: "success", payload: {}}`       | Completed Git operation        |
| Error      | `{type: "error", message: "text"}`     | Invalid commit hash            |
| Progress   | `{type: "progress", stage: "text"}`    | Clone in progress              |

## 🛠️ Operations

### Key Commands

```bash
# View live logs
docker compose logs -f

# Test WebSocket connection
python3 scripts/client.py

# Force certificate renewal 
docker compose exec traefik traefik cert renew

# Scale backend service
docker compose up -d --scale backendim-brain=3

# Access Kubernetes cluster
docker compose exec backendim-brain kubectl get nodes
```

### Monitoring Endpoints

| Endpoint                  | Description                       | Access                        |
|---------------------------|-----------------------------------|-------------------------------|
| `:8080/health`            | Service health status             | Internal only                 |
| `:8080/metrics`           | Prometheus metrics                | Monitoring systems            |
| Traefik Dashboard         | Reverse proxy metrics             | `http://traefik.DOMAIN:8080`  |

## 🔒 Security Architecture

### Protection Layers

1. **Network Security**
   - TLS 1.3 via Let's Encrypt
   - Isolated Docker networks
   - Port 80 auto-redirect to 443

2. **Runtime Security**

   ```dockerfile
   USER backenduser  # Non-root execution
   RUN chmod 0700 /home/backenduser/.kube
   ```

3. **Secret Management**

   ```bash
   # .gitignore protected patterns
   *.key
   *.pem
   .env
   ```

### Audit Controls

```bash
# Verify file permissions
docker compose exec backendim-brain \
  ls -l /home/backenduser/.kube/config

# Check user context
docker compose exec backendim-brain \
  whoami && id  # Should show backenduser
```

## 🚨 Troubleshooting Guide

### Common Issues

**WebSocket Connection Fails**

```bash
# Verify DNS records
dig +short $DOMAIN

# Check certificate status
openssl s_client -connect $DOMAIN:443 -servername $DOMAIN
```

**Kubernetes Auth Failure**

```bash
# Test cluster access
docker compose exec backendim-brain \
  kubectl auth can-i create deployments
```

**Git Operations Hanging**

```bash
# Check repository storage
docker compose exec backendim-brain \
  du -sh /app/repos

# Verify Git version
docker compose exec backendim-brain git --version
```

## 📈 Performance Tuning

### Recommended Deployment

```yaml
# docker-compose.override.yml
services:
  backendim-brain:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
    healthcheck:
      interval: 15s
      timeout: 5s
```

### Scaling Strategy

1. **Horizontal Scaling**

   ```bash
   docker compose up -d --scale backendim-brain=5
   ```

2. **Resource Limits**

   ```yaml
   reservations:
     cpus: '0.5'
     memory: 512M
   ```

## 🤝 Contributing

### Development Workflow

1. Create feature branch

   ```bash
   git checkout -b feature/new-handler
   ```

2. Run local stack

   ```bash
   docker compose -f docker-compose.dev.yml up
   ```

3. Submit PR with:
   - Updated tests
   - Documentation changes
   - Versioned migration scripts

---

**Maintainers**: Obi Madu  
**Status**: Production Ready (v1.3.0)  
**Last Updated**: 2024-02-20
