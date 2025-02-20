#!/bin/bash
POD_NAME=$(kubectl get pod -l app=test-pod -o jsonpath="{.items[0].metadata.name}")

run_tests() {
    local commit_hash=$1
    local test_command=${2:-"pytest tests/"}
    
    echo "=== Testing $commit_hash ==="
    
    kubectl exec $POD_NAME -- rm -rf /app/repo
    kubectl exec $POD_NAME -- git clone https://github.com/yourusername/fastapi-app.git /app/repo
    kubectl exec $POD_NAME -- sh -c "cd /app/repo && git checkout $commit_hash"
    
    kubectl exec $POD_NAME -- sh -c "cd /app/repo && pip install -r requirements.txt"
    
    start=$(date +%s)
    kubectl exec $POD_NAME -- sh -c "cd /app/repo && $test_command"
    status=$?
    end=$(date +%s)
    
    runtime=$((end-start))
    [ $status -eq 0 ] && echo "✅ Passed in ${runtime}s" || echo "❌ Failed after ${runtime}s"
    return $status
}

while true; do
    read -p "Commit hash (q to quit): " commit
    [[ "$commit" == "q" ]] && break
    
    read -p "Test command [pytest tests/]: " cmd
    run_tests "$commit" "${cmd:-pytest tests/}"
    
    read -p "Press Enter to continue..."
done

kubectl delete -f ../deployments/templates/test-pod.yaml
