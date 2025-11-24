package handlers

import (
	"GoChatApp/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReadReceiptHandler struct {
	receiptRepo *repositories.ReadReceiptRepository
}

func NewReadReceiptHandler(receiptRepo *repositories.ReadReceiptRepository) *ReadReceiptHandler {
	return &ReadReceiptHandler{receiptRepo: receiptRepo}
}

// MarkAsRead marks messages as read up to a specific message ID
func (h *ReadReceiptHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		RoomID    uint `json:"room_id" binding:"required"`
		MessageID uint `json:"message_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.receiptRepo.MarkAsRead(userID.(uint), input.RoomID, input.MessageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

// GetReadReceipts gets read receipts for a room
func (h *ReadReceiptHandler) GetReadReceipts(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	receipts, err := h.receiptRepo.GetRoomReadReceipts(uint(roomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch read receipts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"receipts": receipts})
}
