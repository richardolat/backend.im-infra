#!/usr/bin/env python3
import argparse
import json
import subprocess
import sys
import time
from typing import Dict, Any

class TestRunner:
    def __init__(self, namespace: str, repo_url: str, commit: str, test_cmd: str):
        self.namespace = namespace
        self.repo_url = repo_url
        self.commit = commit
        self.test_cmd = test_cmd
        self.pod_name = ""
        self.result: Dict[str, Any] = {
            "commit": commit,
            "status": "unknown",
            "success": False,
            "duration": 0.0,
            "output": {
                "stdout": "",
                "stderr": "",
                "kubectl_errors": []
            }
        }

    def run_kubectl(self, command: List[str]) -> subprocess.CompletedProcess:
        """Execute kubectl command and capture output"""
        full_command = ["kubectl", "-n", self.namespace] + command
        try:
            result = subprocess.run(
                full_command,
                capture_output=True,
                text=True,
                check=True
            )
            return result
        except subprocess.CalledProcessError as e:
            self.result["output"]["kubectl_errors"].append({
                "command": ' '.join(full_command),
                "error": e.stderr.strip()
            })
            raise

    def execute_test_run(self) -> Dict[str, Any]:
        """Main execution flow returning JSON results"""
        start_time = time.time()
        
        try:
            # Setup pod
            self.result["steps"] = []
            self._add_step("Deploying pod", "setup")
            self.run_kubectl(["apply", "-f", "deployments/templates/test-pod.yaml"])
            self.run_kubectl(["wait", "--for=condition=Ready", "pod", "-l", "app=test-pod", "--timeout=120s"])
            
            # Get pod name
            result = self.run_kubectl(["get", "pod", "-l", "app=test-pod", "-o", "jsonpath={.items[0].metadata.name}"])
            self.pod_name = result.stdout.strip()
            self._add_step("Pod ready", "setup_complete")

            # Clone and checkout
            self._add_step("Cloning repository", "source_setup")
            self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                            f"rm -rf /app/repo && git clone {self.repo_url} /app/repo"])
            self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                            f"cd /app/repo && git checkout {self.commit}"])

            # Install dependencies
            self._add_step("Installing dependencies", "dependencies")
            self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                            "cd /app/repo && pip install --quiet -r requirements.txt"])

            # Run tests
            self._add_step("Executing tests", "test_execution")
            test_result = self.run_kubectl(["exec", self.pod_name, "--", "sh", "-c", 
                                          f"cd /app/repo && {self.test_cmd}"])
            
            # Capture results
            self.result["output"]["stdout"] = test_result.stdout
            self.result["output"]["stderr"] = test_result.stderr
            self.result["success"] = True
            self.result["status"] = "passed"

        except subprocess.CalledProcessError as e:
            self.result["status"] = "failed"
            self.result["success"] = False
            if not self.result["output"]["stdout"]:
                self.result["output"]["stdout"] = e.stdout
            if not self.result["output"]["stderr"]:
                self.result["output"]["stderr"] = e.stderr
        except Exception as e:
            self.result["status"] = "error"
            self.result["output"]["stderr"] = str(e)
        finally:
            # Cleanup
            self._add_step("Cleaning up", "cleanup")
            self.run_kubectl(["delete", "-f", "deployments/templates/test-pod.yaml"], check=False)
            self.result["duration"] = round(time.time() - start_time, 2)
        
        return self.result

    def _add_step(self, description: str, step_id: str):
        """Track execution steps"""
        self.result["steps"].append({
            "step": step_id,
            "description": description,
            "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        })

def main():
    parser = argparse.ArgumentParser(description='Single Commit Test Runner')
    parser.add_argument('-n', '--namespace', required=True, help='Target namespace')
    parser.add_argument('-r', '--repo-url', required=True, help='Git repository URL')
    parser.add_argument('-c', '--commit', required=True, help='Commit hash to test')
    parser.add_argument('-t', '--test-cmd', default='pytest tests/', help='Test command')
    
    args = parser.parse_args()
    
    runner = TestRunner(
        namespace=args.namespace,
        repo_url=args.repo_url,
        commit=args.commit,
        test_cmd=args.test_cmd
    )
    
    result = runner.execute_test_run()
    print(json.dumps(result, indent=2))

if __name__ == '__main__':
    main()
