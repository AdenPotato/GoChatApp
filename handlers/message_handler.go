package handlers

import (
	"GoChatApp/models"
	"GoChatApp/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageRepo *repositories.MessageRepository
}

func NewMessageHandler(messageRepo *repositories.MessageRepository) *MessageHandler {
	return &MessageHandler{messageRepo: messageRepo}
}

// GetMessages returns all messages with pagination
func (h *MessageHandler) GetMessages(c *gin.Context) {
	// Get pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	messages, err := h.messageRepo.FindAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage creates a new message
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var input struct {
		Content string `json:"content" binding:"required"`
		RoomID  uint   `json:"room_id" binding:"required"`
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

	// Create message
	message := models.Message{
		UserID:  userID.(uint),
		RoomID:  input.RoomID,
		Content: input.Content,
	}

	if err := h.messageRepo.Create(&message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	// Fetch the message with user and room preloaded
	createdMessage, err := h.messageRepo.FindByID(message.ID)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"message": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": createdMessage})
}

// SearchMessages searches messages by content
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	// Optional room filter
	var roomID uint
	if roomIDStr := c.Query("room_id"); roomIDStr != "" {
		if id, err := strconv.ParseUint(roomIDStr, 10, 32); err == nil {
			roomID = uint(id)
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.messageRepo.Search(query, roomID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"query":    query,
		"count":    len(messages),
	})
}
