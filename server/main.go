package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type Message struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func handleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    defer conn.Close()

    log.Printf("New WebSocket connection established")

    for {
        // Read JSON message from client
        _, rawMessage, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        var clientMsg Message
        if err := json.Unmarshal(rawMessage, &clientMsg); err != nil {
            log.Printf("Error parsing JSON: %v", err)
            continue
        }

        log.Printf("Received message: %+v", clientMsg)

        // Prepare server response
        response := Message{
            Type:    "server_response",
            Payload: map[string]interface{}{
                "received": clientMsg.Payload,
                "status":  "ok",
            },
        }

        // Send JSON response back
        if err := conn.WriteJSON(response); err != nil {
            log.Printf("Error writing response: %v", err)
            break
        }
    }
}

func main() {
    r := gin.Default()
    r.GET("/ws", handleWebSocket)
    
    log.Printf("WebSocket server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
