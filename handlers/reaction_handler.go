package handlers

import (
	"GoChatApp/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	reactionRepo *repositories.ReactionRepository
}

func NewReactionHandler(reactionRepo *repositories.ReactionRepository) *ReactionHandler {
	return &ReactionHandler{reactionRepo: reactionRepo}
}

// ToggleReaction adds or removes a reaction (toggle behavior)
func (h *ReactionHandler) ToggleReaction(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		Emoji string `json:"emoji" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	added, err := h.reactionRepo.Toggle(uint(messageID), userID.(uint), input.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle reaction"})
		return
	}

	action := "removed"
	if added {
		action = "added"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction " + action,
		"added":   added,
		"emoji":   input.Emoji,
	})
}

// GetReactions gets all reactions for a message
func (h *ReactionHandler) GetReactions(c *gin.Context) {
	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	reactions, err := h.reactionRepo.FindByMessageID(uint(messageID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"})
		return
	}

	counts, _ := h.reactionRepo.GetReactionCounts(uint(messageID))

	c.JSON(http.StatusOK, gin.H{
		"reactions": reactions,
		"counts":    counts,
	})
}
