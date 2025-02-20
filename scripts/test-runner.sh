#!/bin/bash
set -eo pipefail

# Parse command-line arguments
usage() {
    echo "Usage: $0 -n NAMESPACE -r REPO_URL -c COMMIT_HASHES [-t TEST_CMD]"
    echo "Example: $0 -n my-tests -r https://github.com/user/repo.git -c abc123,def456 -t 'pytest tests/'"
    exit 1
}

declare -a COMMITS
NAMESPACE=""
REPO_URL=""
TEST_CMD="pytest tests/"

while getopts ":n:r:c:t:" opt; do
    case $opt in
        n) NAMESPACE="$OPTARG" ;;
        r) REPO_URL="$OPTARG" ;;
        c) IFS=',' read -ra COMMITS <<< "$OPTARG" ;;
        t) TEST_CMD="$OPTARG" ;;
        *) usage ;;
    esac
done

[[ -z "$NAMESPACE" || -z "$REPO_URL" || ${#COMMITS[@]} -eq 0 ]] && usage

# Deploy test pod in specified namespace
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f - >/dev/null
kubectl apply -n $NAMESPACE -f ../deployments/templates/test-pod.yaml >/dev/null

# Wait for pod to be ready
echo "Waiting for test pod to initialize..."
kubectl wait -n $NAMESPACE --for=condition=Ready pod -l app=test-pod --timeout=120s

# Get pod name
POD_NAME=$(kubectl get pod -n $NAMESPACE -l app=test-pod -o jsonpath="{.items[0].metadata.name}")

run_tests() {
    local commit_hash=$1
    echo "=== Testing commit $commit_hash in namespace $NAMESPACE ==="
    
    # Clean and clone repo
    kubectl exec -n $NAMESPACE $POD_NAME -- sh -c "rm -rf /app/repo && git clone $REPO_URL /app/repo" >/dev/null
    
    # Checkout specific commit
    kubectl exec -n $NAMESPACE $POD_NAME -- sh -c "cd /app/repo && git checkout $commit_hash" >/dev/null
    
    # Install dependencies
    kubectl exec -n $NAMESPACE $POD_NAME -- sh -c "cd /app/repo && pip install -r requirements.txt" >/dev/null
    
    # Execute tests with timing
    start=$(date +%s)
    echo "Executing: $TEST_CMD"
    kubectl exec -n $NAMESPACE $POD_NAME -- sh -c "cd /app/repo && $TEST_CMD"
    status=$?
    end=$(date +%s)
    
    # Report results
    runtime=$((end-start))
    [ $status -eq 0 ] && echo "✅ Passed in ${runtime}s" || echo "❌ Failed after ${runtime}s"
    return $status
}

# Process all commits
for commit in "${COMMITS[@]}"; do
    run_tests "$commit" || true
    echo "========================================="
done

# Cleanup
echo "Cleaning up test resources..."
kubectl delete -n $NAMESPACE -f ../deployments/templates/test-pod.yaml >/dev/null
