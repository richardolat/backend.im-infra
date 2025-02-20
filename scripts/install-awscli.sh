#!/bin/bash
set -eo pipefail

# Install AWS CLI v2 using pip for Alpine compatibility
pip3 install --upgrade pip
pip3 install --no-cache-dir awscli

# Verify installation
aws --version || { echo "AWS CLI installation failed"; exit 1; }
