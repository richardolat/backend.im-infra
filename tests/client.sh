#!/bin/bash

# WebSocket client test using wscat
# Requires wscat: npm install -g wscat

echo "Connecting to WebSocket server..."

# Prepare JSON message
JSON_MESSAGE='{"type":"client_message","payload":{"message":"Hello Server!","data":{"test":true}}}'

# Send JSON message using wscat
echo "Sending message: $JSON_MESSAGE"
echo $JSON_MESSAGE | wscat -c ws://localhost:8080/ws
