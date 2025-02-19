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

# Create a named pipe
FIFO=$(mktemp -u)
mkfifo "$FIFO"

# Start wscat in background reading from FIFO
wscat -c ws://localhost:8080/ws < "$FIFO" &

# Send message and keep connection open
echo "Starting WebSocket session (Ctrl+C to exit)..."
echo "$JSON_MESSAGE" > "$FIFO"
 
# Keep script running and reading responses
sleep 1  # Give time for initial response
tail -f "$FIFO"  # Keep pipe open for reading responses

# Cleanup
rm "$FIFO"
