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

# Extract values from JSON
REPO_URL=$(echo "$JSON_MESSAGE" | jq -r '.repoURL')
COMMIT_HASH=$(echo "$JSON_MESSAGE" | jq -r '.commitHash')

# Handle git operations
./git_handler.sh "$REPO_URL" "$COMMIT_HASH"

# Send JSON message using wscat
echo "Sending message: $JSON_MESSAGE"
echo "$JSON_MESSAGE" | wscat -c ws://localhost:8080/ws
