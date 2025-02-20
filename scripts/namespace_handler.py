#!/usr/bin/env python3
import json
import subprocess
import sys
from datetime import datetime

def main():
    try:
        if len(sys.argv) < 2:
            print(json.dumps({
                "status": "error",
                "message": "Namespace identifier required"
            }))
            sys.exit(1)

        identifier = sys.argv[1]
        namespace = f"im-{identifier}".lower()
        timestamp = datetime.utcnow().isoformat() + "Z"

        # Check namespace existence
        check = subprocess.run(
            ["kubectl", "get", "namespace", namespace],
            capture_output=True,
            text=True
        )

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
            capture_output=True,
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
            "namespace": namespace
        }))
        sys.exit(1)

    except Exception as e:
        print(json.dumps({
            "status": "error",
            "message": f"Handler failure: {str(e)}"
        }))
        sys.exit(1)

if __name__ == "__main__":
    main()
