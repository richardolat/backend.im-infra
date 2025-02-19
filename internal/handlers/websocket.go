package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "websocket-git/internal/models"
    "websocket-git/internal/services"
)

type WebSocketHandler struct {
    upgrader   websocket.Upgrader
    gitService *services.GitService
}

func NewWebSocketHandler(gitService *services.GitService) *WebSocketHandler {
    return &WebSocketHandler{
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        },
        gitService: gitService,
    }
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }
    defer conn.Close()

    log.Printf("New WebSocket connection established from %s", c.Request.RemoteAddr)

    for {
        messageType, rawMessage, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message from %s: %v", c.Request.RemoteAddr, err)
            break
        }

        log.Printf("[%s] Raw message received (type: %d): %s", 
            c.Request.RemoteAddr, 
            messageType, 
            string(rawMessage))

        var gitMsg models.GitMessage
        if err := json.Unmarshal(rawMessage, &gitMsg); err != nil {
            log.Printf("[%s] Error parsing JSON: %v", c.Request.RemoteAddr, err)
            log.Printf("[%s] Invalid JSON content: %s", c.Request.RemoteAddr, string(rawMessage))
            sendError(conn, "Invalid message format")
            continue
        }

        log.Printf("[%s] Parsed message: %+v", c.Request.RemoteAddr, gitMsg)

        // Handle git operations
        log.Printf("[%s] Starting git operation for repo: %s at commit: %s", 
            c.Request.RemoteAddr, 
            gitMsg.RepoURL, 
            gitMsg.CommitHash)

        if err := h.gitService.HandleRepository(gitMsg.RepoURL, gitMsg.CommitHash); err != nil {
            log.Printf("[%s] Git operation failed: %v", c.Request.RemoteAddr, err)
            sendError(conn, "Git operation failed")
            continue
        }

        response := models.Response{
            Type: "success",
            Payload: map[string]interface{}{
                "message": "Git operation completed successfully",
                "repoURL": gitMsg.RepoURL,
                "commit":  gitMsg.CommitHash,
            },
        }

        if err := conn.WriteJSON(response); err != nil {
            log.Printf("Error writing response: %v", err)
            break
        }
    }
}

func sendError(conn *websocket.Conn, message string) {
    response := models.Response{
        Type: "error",
        Payload: map[string]interface{}{
            "message": message,
        },
    }
    conn.WriteJSON(response)
}
