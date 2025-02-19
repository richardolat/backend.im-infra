#!/bin/bash

# Simple WebSocket client test using wscat
# Requires wscat: npm install -g wscat

echo "Connecting to WebSocket server..."
wscat -c ws://localhost:8080/ws <<EOF
Hello Server!
This is a test message
EOF
