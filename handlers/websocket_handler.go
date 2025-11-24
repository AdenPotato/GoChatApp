package handlers

import (
	"GoChatApp/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development (adjust for production)
		return true
	},
}

type WebSocketHandler struct {
	Hub *Hub
}

func NewWebSocketHandler() *WebSocketHandler {
	hub := NewHub()
	go hub.Run() // Start the hub in a goroutine
	return &WebSocketHandler{
		Hub: hub,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get JWT token from query parameter (WebSocket doesn't support headers easily)
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	// Validate JWT token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	userID := claims.UserID
	username := claims.Username

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create new client
	client := &Client{
		Hub:      h.Hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		Rooms:    make(map[uint]bool),
	}

	// Register client with hub
	h.Hub.Register <- client

	// Start client's read and write pumps in separate goroutines
	go client.WritePump()
	go client.ReadPump()
}
