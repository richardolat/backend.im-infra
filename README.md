# Automated Commit Testing System for Backend.im

![Progress Workflow Diagram](./workflow.png)

A distributed system for testing Git commits in isolated Kubernetes environments using WebSocket communication. Built to catch regressions before they hit production.

## ⚡️ Features

- **Real-time Test Pipeline** - WebSocket-based test orchestration
- **K8s Isolation** - Dedicated namespaces per test session
- **Commit-by-Commit Testing** - Test multiple commits in sequence
- **Custom Test Commands** - Override default test behavior
- **Failure Forensics** - Detailed error logging and output capture

## 🏗 Architecture

```mermaid
sequenceDiagram
    Client->>Server: Connect via WS
    Client->>Server: Send commit+test command
    Server->>K8s: Create namespace
    Server->>K8s: Deploy test pod
    K8s->>TestPod: Clone repo
    K8s->>TestPod: Checkout commit
    K8s->>TestPod: Run tests
    TestPod->>Server: JSON results
    Server->>Client: Return results
```

## 🛠 Prerequisites

- Kubernetes cluster (Docker Desktop K8s works)
- `kubectl` configured and working
- Python 3.9+ (client side)
- Go 1.18+ (server side)

## 🚀 Quick Start

### 1. Clone and Prepare

```bash
git clone https://git.obimadu.pro/your-repo.git
cd automated-testing
```

### 2. Server Setup

```bash
# Install dependencies
make deps

# Start server (port 8080)
make run
```

### 3. Client Setup

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r app/requirements.txt

# Copy and edit config
cp config.example.json config.json
nano config.json  # Set your repo and commits
```

### 4. Run Tests

```bash
python3 scripts/client.py
```

## ⚙️ Configuration

**config.json**

```json
{
  "ws_url": "ws://your-server:8080/ws",
  "repo_url": "https://github.com/your/repo.git",
  "user_id": "your-user-id",
  "chat_id": "your-chat-id",
  "project_type": "fastapi",
  "test_command": "pytest -v tests/network/",
  "commits": ["HEAD~2", "HEAD~1", "HEAD"]
}
```

## 🔥 Custom Test Commands

Override default test behavior in `config.json`:

```json
{
  "test_command": "mvn test -Dtest=SmokeTestSuite"
}
```

Supported formats:

- Simple commands: `npm test`
- Chained commands: `make test && coverage report`
- Script paths: `bash tests/e2e.sh`

## 🚨 Troubleshooting

### "Missing deployment template" Error

Create required K8s templates:

```bash
mkdir -p deployments/templates/fastapi/
# Add your test-pod.yaml here
```

### WebSocket Connection Issues

```bash
# Verify server is running
curl -I http://localhost:8080/health

# Check firewall rules
sudo ufw allow 8080/tcp
```

### K8s Pod Failing

Check pod logs:

```bash
kubectl logs -n im-chat-14-user-14 test-pod
```

## 📋 Test Flow Breakdown

1. Client connects via WebSocket
2. Server creates isolated K8s namespace
3. Test pod deploys with project-specific template
4. Repo is cloned/updated in pod
5. Target commit is checked out
6. Custom test command executes
7. Results stream back via WebSocket

## ⁉️ FAQ

**Q: How are resources cleaned up?**  
Namespaces persist for debugging. Manual cleanup:

```bash
kubectl delete namespace im-chat-14-user-14
```

**Q: Can I modify the test pod?**  
Edit templates in `deployments/templates/<project-type>/test-pod.yaml`

**Q: What's the performance overhead?**  
Each test run uses ~300MB RAM. Scale namespaces as needed.

## 📜 License

AGPL-3.0 - Fuck around and find out edition
