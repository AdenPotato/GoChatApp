package handlers

import (
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
	var message struct {
		Content string `json:"content" binding:"required"`
		User    string `json:"user"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement real message creation with authenticated user
	c.JSON(http.StatusCreated, gin.H{"message": "Message sent"})
}
