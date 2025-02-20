# Backendim Brain

## Requirements
- Docker 20.10+
- AWS credentials with EKS access

## Setup
1. Copy `.env.example` to `.env`
2. Fill in your credentials
3. Run `make build && make up`

## Environment Variables
| Variable | Description |
|----------|-------------|
| AWS_ACCESS_KEY | AWS IAM Access Key |
| AWS_SECRET | AWS IAM Secret Key |
| AWS_REGION | AWS Region |
| DOMAIN | Domain for Traefik |
| KUBE_CLUSTER_NAME | EKS Cluster Name |
