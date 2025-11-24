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
	Rooms    map[uint]bool // Rooms the client has joined
}

// RoomMessage represents a message to be sent to a specific room
type RoomMessage struct {
	RoomID  uint
	Message []byte
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients
	Clients map[*Client]bool

	// Clients organized by room
	Rooms map[uint]map[*Client]bool

	// Inbound messages from clients
	Broadcast chan []byte

	// Room-specific broadcast
	RoomBroadcast chan RoomMessage

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
		Clients:       make(map[*Client]bool),
		Rooms:         make(map[uint]map[*Client]bool),
		Broadcast:     make(chan []byte, 256),
		RoomBroadcast: make(chan RoomMessage, 256),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
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
				// Remove client from all rooms
				for roomID := range client.Rooms {
					if room, exists := h.Rooms[roomID]; exists {
						delete(room, client)
						if len(room) == 0 {
							delete(h.Rooms, roomID)
						}
					}
				}
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

		case roomMsg := <-h.RoomBroadcast:
			h.BroadcastToRoom(roomMsg.RoomID, roomMsg.Message)
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

// BroadcastToRoom sends a message to all clients in a specific room
func (h *Hub) BroadcastToRoom(roomID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, exists := h.Rooms[roomID]; exists {
		for client := range room {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
				delete(room, client)
			}
		}
	}
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Initialize room if it doesn't exist
	if _, exists := h.Rooms[roomID]; !exists {
		h.Rooms[roomID] = make(map[*Client]bool)
	}

	h.Rooms[roomID][client] = true
	client.Rooms[roomID] = true

	log.Printf("Client %s joined room %d", client.Username, roomID)

	// Notify room members
	joinMsg := map[string]interface{}{
		"type":     "room_user_joined",
		"room_id":  roomID,
		"user_id":  client.UserID,
		"username": client.Username,
	}
	if msgBytes, err := json.Marshal(joinMsg); err == nil {
		for c := range h.Rooms[roomID] {
			select {
			case c.Send <- msgBytes:
			default:
			}
		}
	}
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(client *Client, roomID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.Rooms[roomID]; exists {
		delete(room, client)
		delete(client.Rooms, roomID)

		// Clean up empty rooms
		if len(room) == 0 {
			delete(h.Rooms, roomID)
		}

		log.Printf("Client %s left room %d", client.Username, roomID)

		// Notify remaining room members
		leaveMsg := map[string]interface{}{
			"type":     "room_user_left",
			"room_id":  roomID,
			"user_id":  client.UserID,
			"username": client.Username,
		}
		if msgBytes, err := json.Marshal(leaveMsg); err == nil {
			for c := range h.Rooms[roomID] {
				select {
				case c.Send <- msgBytes:
				default:
				}
			}
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

// GetRoomUsers returns a list of users in a specific room
func (h *Hub) GetRoomUsers(roomID uint) []map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]map[string]interface{}, 0)
	if room, exists := h.Rooms[roomID]; exists {
		for client := range room {
			users = append(users, map[string]interface{}{
				"user_id":  client.UserID,
				"username": client.Username,
			})
		}
	}
	return users
}
