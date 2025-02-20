#!/bin/bash
set -eo pipefail

# Initialize kubeconfig
aws eks update-kubeconfig \
  --name ${KUBE_CLUSTER_NAME} \
  --region ${AWS_DEFAULT_REGION} \
  --kubeconfig /root/.kube/config

# Verify cluster access
kubectl cluster-info --request-timeout=5s
