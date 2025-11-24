package handlers

import (
	"GoChatApp/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlockHandler struct {
	blockRepo *repositories.BlockRepository
}

func NewBlockHandler(blockRepo *repositories.BlockRepository) *BlockHandler {
	return &BlockHandler{blockRepo: blockRepo}
}

// BlockUser blocks another user
func (h *BlockHandler) BlockUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blockedIDStr := c.Param("id")
	blockedID, err := strconv.ParseUint(blockedIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if uint(blockedID) == userID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot block yourself"})
		return
	}

	// Check if already blocked
	isBlocked, _ := h.blockRepo.IsBlocked(userID.(uint), uint(blockedID))
	if isBlocked {
		c.JSON(http.StatusConflict, gin.H{"error": "User already blocked"})
		return
	}

	if err := h.blockRepo.Block(userID.(uint), uint(blockedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User blocked"})
}

// UnblockUser unblocks a user
func (h *BlockHandler) UnblockUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blockedIDStr := c.Param("id")
	blockedID, err := strconv.ParseUint(blockedIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.blockRepo.Unblock(userID.(uint), uint(blockedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unblocked"})
}

// GetBlockedUsers gets list of blocked users
func (h *BlockHandler) GetBlockedUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blocks, err := h.blockRepo.GetBlockedUsers(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocked users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocked_users": blocks})
}
