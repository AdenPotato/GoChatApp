package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	// Will add WebSocket hub here later
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// TODO: Implement WebSocket upgrade logic in Phase 4
	c.JSON(http.StatusNotImplemented, gin.H{"message": "WebSocket not implemented yet"})
}
