#!/bin/bash
set -eo pipefail

# Validate environment variables
required_vars=(AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_REGION KUBE_CLUSTER_NAME)
for var in "${required_vars[@]}"; do
    if [[ -z "${!var}" ]]; then
        echo "ERROR: Missing required environment variable $var" >&2
        exit 1
    fi
done

# Configure AWS CLI credentials
aws configure set aws_access_key_id ${AWS_ACCESS_KEY_ID}
aws configure set aws_secret_access_key ${AWS_SECRET_ACCESS_KEY}
aws configure set default.region ${AWS_REGION}

# Generate kubeconfig
aws eks update-kubeconfig \
    --name ${KUBE_CLUSTER_NAME} \
    --region ${AWS_REGION} \
    --kubeconfig ${KUBECONFIG} \
    --alias automated-cluster

# Verify cluster access with timeout
if ! kubectl cluster-info --request-timeout=10s; then
    echo "Failed to connect to Kubernetes cluster" >&2
    exit 1
fi

# Verify kubectl version
echo "Kubectl Version:"
kubectl version --client -o json | jq -r '.clientVersion.gitVersion'

# Verify AWS credentials
echo "AWS Identity:"
aws sts get-caller-identity
