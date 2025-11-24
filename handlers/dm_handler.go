package handlers

import (
	"GoChatApp/models"
	"GoChatApp/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DMHandler struct {
	dmRepo *repositories.DMRepository
}

func NewDMHandler(dmRepo *repositories.DMRepository) *DMHandler {
	return &DMHandler{dmRepo: dmRepo}
}

// GetConversations gets all DM conversations for authenticated user
func (h *DMHandler) GetConversations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	convs, err := h.dmRepo.GetUserConversations(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": convs})
}

// StartConversation starts or gets existing conversation with another user
func (h *DMHandler) StartConversation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.UserID == userID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot start conversation with yourself"})
		return
	}

	conv, err := h.dmRepo.FindOrCreateConversation(userID.(uint), input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversation": conv})
}

// GetMessages gets messages in a conversation
func (h *DMHandler) GetMessages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	convIDStr := c.Param("id")
	convID, err := strconv.ParseUint(convIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Verify user is participant
	isParticipant, err := h.dmRepo.IsParticipant(uint(convID), userID.(uint))
	if err != nil || !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a participant in this conversation"})
		return
	}

	// Parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.dmRepo.GetMessages(uint(convID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Mark messages as read
	h.dmRepo.MarkAsRead(uint(convID), userID.(uint))

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage sends a DM
func (h *DMHandler) SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	convIDStr := c.Param("id")
	convID, err := strconv.ParseUint(convIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// Verify user is participant
	isParticipant, err := h.dmRepo.IsParticipant(uint(convID), userID.(uint))
	if err != nil || !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a participant in this conversation"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &models.DirectMessage{
		ConversationID: uint(convID),
		SenderID:       userID.(uint),
		Content:        input.Content,
	}

	if err := h.dmRepo.CreateMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

// GetUnreadCount gets total unread DM count
func (h *DMHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	count, err := h.dmRepo.GetUnreadCount(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
