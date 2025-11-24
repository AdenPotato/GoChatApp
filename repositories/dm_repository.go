package repositories

import (
	"GoChatApp/models"

	"gorm.io/gorm"
)

type DMRepository struct {
	db *gorm.DB
}

func NewDMRepository(db *gorm.DB) *DMRepository {
	return &DMRepository{db: db}
}

// FindOrCreateConversation gets existing conversation or creates new one
func (r *DMRepository) FindOrCreateConversation(user1ID, user2ID uint) (*models.Conversation, error) {
	var conv models.Conversation

	// Always store with smaller ID first for consistency
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	err := r.db.Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).
		Preload("User1").Preload("User2").
		First(&conv).Error

	if err == gorm.ErrRecordNotFound {
		conv = models.Conversation{
			User1ID: user1ID,
			User2ID: user2ID,
		}
		if err := r.db.Create(&conv).Error; err != nil {
			return nil, err
		}
		r.db.Preload("User1").Preload("User2").First(&conv, conv.ID)
	} else if err != nil {
		return nil, err
	}

	return &conv, nil
}

// GetUserConversations gets all conversations for a user
func (r *DMRepository) GetUserConversations(userID uint) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := r.db.Where("user1_id = ? OR user2_id = ?", userID, userID).
		Preload("User1").Preload("User2").
		Find(&convs).Error
	return convs, err
}

// CreateMessage creates a new DM
func (r *DMRepository) CreateMessage(msg *models.DirectMessage) error {
	return r.db.Create(msg).Error
}

// GetMessages gets messages in a conversation with pagination
func (r *DMRepository) GetMessages(conversationID uint, limit, offset int) ([]models.DirectMessage, error) {
	var messages []models.DirectMessage
	err := r.db.Where("conversation_id = ?", conversationID).
		Preload("Sender").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// MarkAsRead marks all messages in conversation as read for a user
func (r *DMRepository) MarkAsRead(conversationID, userID uint) error {
	return r.db.Model(&models.DirectMessage{}).
		Where("conversation_id = ? AND sender_id != ? AND read = ?", conversationID, userID, false).
		Update("read", true).Error
}

// GetUnreadCount gets unread message count for a user
func (r *DMRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64

	// Get all conversations where user is participant
	var convIDs []uint
	r.db.Model(&models.Conversation{}).
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Pluck("id", &convIDs)

	if len(convIDs) == 0 {
		return 0, nil
	}

	err := r.db.Model(&models.DirectMessage{}).
		Where("conversation_id IN ? AND sender_id != ? AND read = ?", convIDs, userID, false).
		Count(&count).Error

	return count, err
}

// IsParticipant checks if user is part of conversation
func (r *DMRepository) IsParticipant(conversationID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Conversation{}).
		Where("id = ? AND (user1_id = ? OR user2_id = ?)", conversationID, userID, userID).
		Count(&count).Error
	return count > 0, err
}
