#!/usr/bin/env python3
import argparse
import subprocess
import sys
import time
from typing import List, Tuple

class TestRunner:
    def __init__(self, namespace: str, repo_url: str, commits: List[str], test_cmd: str):
        self.namespace = namespace
        self.repo_url = repo_url
        self.commits = commits
        self.test_cmd = test_cmd
        self.pod_name = ""

    def run_kubectl(self, command: List[str], check: bool = True) -> subprocess.CompletedProcess:
        """Execute a kubectl command with error handling"""
        result = subprocess.run(
            ["kubectl", "-n", self.namespace] + command,
            capture_output=True,
            text=True
        )
        if check and result.returncode != 0:
            print(f"⚠️ Command failed: {' '.join(command)}")
            print(f"Stderr: {result.stderr.strip()}")
            sys.exit(1)
        return result

    def setup_pod(self):
        """Deploy and wait for test pod"""
        print(f"🚀 Deploying test pod to namespace {self.namespace}")
        self.run_kubectl(["apply", "-f", "deployments/templates/test-pod.yaml"])
        
        print("⏳ Waiting for pod to be ready...")
        self.run_kubectl(["wait", "--for=condition=Ready", "pod", "-l", "app=test-pod", "--timeout=120s"])
        
        # Get pod name
        result = self.run_kubectl(["get", "pod", "-l", "app=test-pod", "-o", "jsonpath={.items[0].metadata.name}"])
        self.pod_name = result.stdout.strip()

    def cleanup(self):
        """Clean up Kubernetes resources"""
        print("\n🧹 Cleaning up test resources...")
        self.run_kubectl(["delete", "-f", "deployments/templates/test-pod.yaml"], check=False)

    def run_tests(self, commit: str) -> Tuple[int, float]:
        """Execute test sequence for a single commit"""
        print(f"\n🔍 Testing commit {commit}")
        
        # Clean and clone repo
        self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                        f"rm -rf /app/repo && git clone {self.repo_url} /app/repo"])
        
        # Checkout specific commit
        self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                        f"cd /app/repo && git checkout {commit}"])
        
        # Install dependencies
        self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                        "cd /app/repo && pip install -r requirements.txt"])
        
        # Execute tests
        print(f"⚙️ Executing: {self.test_cmd}")
        start_time = time.time()
        result = self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                                 f"cd /app/repo && {self.test_cmd}"], check=False)
        duration = time.time() - start_time
        
        # Print results
        if result.returncode == 0:
            print(f"✅ Passed in {duration:.1f}s")
            print(result.stdout)
        else:
            print(f"❌ Failed after {duration:.1f}s")
            print(result.stderr)
        
        return (result.returncode, duration)

def main():
    parser = argparse.ArgumentParser(description='Kubernetes Test Runner')
    parser.add_argument('-n', '--namespace', required=True, help='Existing namespace for testing')
    parser.add_argument('-r', '--repo-url', required=True, help='Git repository URL to test')
    parser.add_argument('-c', '--commits', required=True, 
                      help='Comma-separated list of commit hashes to test')
    parser.add_argument('-t', '--test-cmd', default='pytest tests/',
                      help='Test command to execute (default: pytest tests/)')
    
    args = parser.parse_args()
    commits = args.commits.split(',')
    
    runner = TestRunner(
        namespace=args.namespace,
        repo_url=args.repo_url,
        commits=commits,
        test_cmd=args.test_cmd
    )
    
    try:
        runner.setup_pod()
        
        results = {}
        for commit in commits:
            exit_code, duration = runner.run_tests(commit)
            results[commit] = {
                'status': 'passed' if exit_code == 0 else 'failed',
                'duration': duration
            }
            print("=" * 60)
        
        # Print summary
        print("\n📊 Test Summary:")
        for commit, data in results.items():
            print(f"{commit[:7]}: {data['status'].upper()} ({data['duration']:.1f}s)")
        
    finally:
        runner.cleanup()

if __name__ == '__main__':
    main()
