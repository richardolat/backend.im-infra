package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "websocket-git/internal/handlers"
    "websocket-git/internal/services"
)

func main() {
    // Initialize services
    gitService := services.NewGitService()

    // Initialize handlers
    wsHandler := handlers.NewWebSocketHandler(gitService)

    // Setup router
    r := gin.Default()
    r.GET("/ws", wsHandler.HandleConnection)
    
    log.Printf("WebSocket server starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
