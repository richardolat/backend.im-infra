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
    upgrader         websocket.Upgrader
    namespaceService *services.NamespaceService  // Correct field name
}

func NewWebSocketHandler(namespaceService *services.NamespaceService) *WebSocketHandler {
    return &WebSocketHandler{
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true  // Adjust origin validation for production
            },
        },
        namespaceService: namespaceService,  // Remove gitService reference
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

        // Handle namespace operations
        result, err := h.namespaceService.HandleNamespace(gitMsg.ChatID, gitMsg.UserID)
        if err != nil {
            log.Printf("[%s] Namespace operation failed: %v", c.Request.RemoteAddr, err)
            sendError(conn, "Namespace operation failed: "+err.Error())
            continue
        }

        response := models.Response{
            Type:    "namespace_status",
            Payload: result,
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
