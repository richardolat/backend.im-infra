#!/bin/bash
set -eo pipefail

# Install AWS CLI v2 with proper verification
TMPDIR=$(mktemp -d)
cd ${TMPDIR}

curl --fail "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip -q awscliv2.zip
./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli --update

# Verify installation
aws --version || { echo "AWS CLI installation failed"; exit 1; }

# Cleanup
cd -
rm -rf ${TMPDIR}
