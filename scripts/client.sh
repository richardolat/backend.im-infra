#!/bin/bash

# WebSocket client test using wscat
# Requires wscat: npm install -g wscat

echo "Connecting to WebSocket server..."

# Prepare JSON message
read -r -d '' JSON_MESSAGE <<'EOF'
{
    "userId": "someId",
    "chatId": "chatId",
    "repoURL": "https://github.com/obiMadu/hng12-stage2",
    "commitHash": "2acf0f9a74b83bc881aa2f06235b8c927892d28a"
}
EOF

# Send JSON message and keep connection open interactively
echo "Starting WebSocket session (Ctrl+C to exit)..."
wscat -c ws://localhost:8080/ws <<EOF
$JSON_MESSAGE
EOF
