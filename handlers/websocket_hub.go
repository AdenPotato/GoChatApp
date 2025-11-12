package handlers

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uint
	Username string
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients
	Clients map[*Client]bool

	// Inbound messages from clients
	Broadcast chan []byte

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub and handles client registration, unregistration, and broadcasting
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (ID: %d). Total clients: %d", client.Username, client.UserID, len(h.Clients))

			// Broadcast user join notification
			joinMsg := map[string]interface{}{
				"type":     "user_joined",
				"user_id":  client.UserID,
				"username": client.Username,
			}
			if msgBytes, err := json.Marshal(joinMsg); err == nil {
				h.BroadcastToAll(msgBytes)
			}

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %s (ID: %d). Total clients: %d", client.Username, client.UserID, len(h.Clients))

				// Broadcast user leave notification
				leaveMsg := map[string]interface{}{
					"type":     "user_left",
					"user_id":  client.UserID,
					"username": client.Username,
				}
				if msgBytes, err := json.Marshal(leaveMsg); err == nil {
					h.BroadcastToAll(msgBytes)
				}
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.BroadcastToAll(message)
		}
	}
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			// Client's send channel is full, close and remove client
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

// GetConnectedUsers returns a list of currently connected users
func (h *Hub) GetConnectedUsers() []map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]map[string]interface{}, 0, len(h.Clients))
	for client := range h.Clients {
		users = append(users, map[string]interface{}{
			"user_id":  client.UserID,
			"username": client.Username,
		})
	}
	return users
}
