#!/usr/bin/env python3
import json
import subprocess
import sys
from datetime import datetime

def main():
    try:
        if len(sys.argv) != 3:
            print(json.dumps({
                "status": "error",
                "message": "Exactly 2 arguments required: chatID and userID"
            }))
            sys.exit(1)

        chat_id = sys.argv[1]
        user_id = sys.argv[2]
        namespace = f"im-{chat_id}-{user_id}".lower()
        timestamp = datetime.utcnow().isoformat() + "Z"

        # Check namespace existence
        check = subprocess.run(
            ["kubectl", "get", "namespace", namespace],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )

        # Add kubectl error logging for check command
        if check.returncode != 0:
            sys.stderr.write(f"kubectl check error: {check.stderr.strip()}\n")

        if check.returncode == 0:
            print(json.dumps({
                "status": "exists",
                "namespace": namespace,
                "timestamp": timestamp
            }))
            sys.exit(0)

        # Create namespace
        create = subprocess.run(
            ["kubectl", "create", "namespace", namespace],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )

        if create.returncode == 0:
            print(json.dumps({
                "status": "created",
                "namespace": namespace,
                "timestamp": timestamp
            }))
            sys.exit(0)

        print(json.dumps({
            "status": "error", 
            "message": create.stderr.strip(),
            "kubectl_error": create.stderr.strip(),  # Explicit error field
            "namespace": namespace
        }))
        sys.exit(1)

    except Exception as e:
        print(json.dumps({
            "status": "error",
            "message": f"Script execution failed: {str(e)}",
            "python_error": str(e)
        }))
        sys.exit(1)

if __name__ == "__main__":
    main()
