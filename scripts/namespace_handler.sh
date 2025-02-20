#!/bin/bash
set -eo pipefail

# Validate input parameter
if [ -z "$1" ]; then
  echo '{"status": "error", "message": "Namespace name parameter is required"}'
  exit 1
fi

# Generate namespace name (lowercase)
NAMESPACE_NAME=$(echo "im-$1" | tr '[:upper:]' '[:lower:]')

# Check if namespace exists
if kubectl get namespace "$NAMESPACE_NAME" >/dev/null 2>&1; then
  echo '{"status": "exists", "namespace": "'"$NAMESPACE_NAME"'", "timestamp": "'$(date -u +%FT%TZ)'"}'
else
  # Create namespace if it doesn't exist
  if kubectl create namespace "$NAMESPACE_NAME" >/dev/null 2>&1; then
    echo '{"status": "created", "namespace": "'"$NAMESPACE_NAME"'", "timestamp": "'$(date -u +%FT%TZ)'"}'
  else
    echo '{"status": "error", "message": "Failed to create namespace", "namespace": "'"$NAMESPACE_NAME"'"}'
    exit 1
  fi
fi
