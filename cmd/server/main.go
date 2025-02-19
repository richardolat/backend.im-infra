package main

import (
	"log"
	"websocket-git/internal/handlers"
	"websocket-git/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to debug mode
	gin.SetMode(gin.DebugMode)

	// Initialize services
	gitService := services.NewGitService()

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(gitService)

	// Setup router with logger and recovery middleware
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ws", wsHandler.HandleConnection)

	log.Printf("WebSocket server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
