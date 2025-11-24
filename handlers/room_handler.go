package handlers

import (
	"GoChatApp/models"
	"GoChatApp/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomRepo *repositories.RoomRepository
}

func NewRoomHandler(roomRepo *repositories.RoomRepository) *RoomHandler {
	return &RoomHandler{roomRepo: roomRepo}
}

// GetRooms returns all rooms
func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.roomRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// GetRoom returns a specific room by ID
func (h *RoomHandler) GetRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	room, err := h.roomRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// CreateRoom creates a new room
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type"` // public, private, direct
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Default to public if type not specified
	roomType := input.Type
	if roomType == "" {
		roomType = "public"
	}

	room := models.Room{
		Name:      input.Name,
		Type:      roomType,
		CreatedBy: userID.(uint),
	}

	if err := h.roomRepo.Create(&room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	// Auto-join creator to the room
	if err := h.roomRepo.AddMember(room.ID, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add creator to room"})
		return
	}

	// Fetch room with relations
	createdRoom, err := h.roomRepo.FindByID(room.ID)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"room": room})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"room": createdRoom})
}

// JoinRoom adds the authenticated user to a room
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	idStr := c.Param("id")
	roomID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if room exists
	room, err := h.roomRepo.FindByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Check if already a member
	isMember, err := h.roomRepo.IsMember(uint(roomID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check membership"})
		return
	}
	if isMember {
		c.JSON(http.StatusConflict, gin.H{"error": "Already a member of this room"})
		return
	}

	// Add member
	if err := h.roomRepo.AddMember(uint(roomID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined room successfully", "room": room})
}

// LeaveRoom removes the authenticated user from a room
func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	idStr := c.Param("id")
	roomID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if room exists
	_, err = h.roomRepo.FindByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Check if member
	isMember, err := h.roomRepo.IsMember(uint(roomID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check membership"})
		return
	}
	if !isMember {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not a member of this room"})
		return
	}

	// Remove member
	if err := h.roomRepo.RemoveMember(uint(roomID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left room successfully"})
}
