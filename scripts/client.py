#!/usr/bin/env python3
import json
import signal
import sys
from websocket import WebSocketApp

def main():
    # Configure WebSocket connection
    ws_url = "ws://localhost:8080/ws"
    message = {
        "userId": "someId",
        "chatId": "chatId",
        "repoURL": "https://github.com/obiMadu/hng12-stage2",
        "commitHash": "2acf0f9a74b83bc881aa2f06235b8c927892d28a"
    }

    # Configure graceful shutdown
    def signal_handler(sig, frame):
        print("\nClosing connection...")
        ws.close()
        sys.exit(0)

    signal.signal(signal.SIGINT, signal_handler)

    # WebSocket event handlers
    def on_open(ws):
        print(f"Connected to {ws_url}")
        print("Sending initial message:")
        print(json.dumps(message, indent=2))
        ws.send(json.dumps(message))

    def on_message(ws, message):
        print("\nReceived response:")
        try:
            response = json.loads(message)
            print(json.dumps(response, indent=2))
        except json.JSONDecodeError:
            print(f"Non-JSON response: {message}")

    def on_error(ws, error):
        print(f"\nWebSocket error: {error}")

    def on_close(ws, close_status_code, close_msg):
        print("\nConnection closed")
        if close_status_code or close_msg:
            print(f"Close status code: {close_status_code}")
            print(f"Close message: {close_msg}")

    # Create and run WebSocket client
    print("Starting WebSocket client...")
    ws = WebSocketApp(ws_url,
                      on_open=on_open,
                      on_message=on_message,
                      on_error=on_error,
                      on_close=on_close)

    print("Press Ctrl+C to exit")
    ws.run_forever()

if __name__ == "__main__":
    main()
