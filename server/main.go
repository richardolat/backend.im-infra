package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // Allow all origins for testing
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    defer conn.Close()

    log.Printf("New WebSocket connection established")

    for {
        // Read message from client
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }
        
        log.Printf("Received message: %s", message)

        // Echo the message back
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            log.Printf("Error writing message: %v", err)
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleWebSocket)
    
    port := ":8080"
    fmt.Printf("WebSocket server starting on %s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
