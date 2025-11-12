package handlers

import (
	"log"
	"net/http"
	"strconv"

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
	// Get user information from query parameters (for simple auth)
	// In production, you'd validate JWT token here
	userIDStr := c.Query("user_id")
	username := c.Query("username")

	if userIDStr == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and username required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

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
		UserID:   uint(userID),
		Username: username,
	}

	// Register client with hub
	h.Hub.Register <- client

	// Start client's read and write pumps in separate goroutines
	go client.WritePump()
	go client.ReadPump()
}
