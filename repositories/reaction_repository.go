package repositories

import (
	"GoChatApp/models"

	"gorm.io/gorm"
)

type ReactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

// Create adds a reaction
func (r *ReactionRepository) Create(reaction *models.Reaction) error {
	return r.db.Create(reaction).Error
}

// Delete removes a reaction
func (r *ReactionRepository) Delete(messageID, userID uint, emoji string) error {
	return r.db.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).
		Delete(&models.Reaction{}).Error
}

// FindByMessageID gets all reactions for a message
func (r *ReactionRepository) FindByMessageID(messageID uint) ([]models.Reaction, error) {
	var reactions []models.Reaction
	err := r.db.Where("message_id = ?", messageID).Preload("User").Find(&reactions).Error
	return reactions, err
}

// GetReactionCounts gets reaction counts grouped by emoji for a message
func (r *ReactionRepository) GetReactionCounts(messageID uint) (map[string]int64, error) {
	var results []struct {
		Emoji string
		Count int64
	}

	err := r.db.Model(&models.Reaction{}).
		Select("emoji, COUNT(*) as count").
		Where("message_id = ?", messageID).
		Group("emoji").
		Find(&results).Error

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Emoji] = r.Count
	}

	return counts, err
}

// Exists checks if a user already reacted with specific emoji
func (r *ReactionRepository) Exists(messageID, userID uint, emoji string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Reaction{}).
		Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).
		Count(&count).Error
	return count > 0, err
}

// Toggle adds reaction if not exists, removes if exists
func (r *ReactionRepository) Toggle(messageID, userID uint, emoji string) (added bool, err error) {
	exists, err := r.Exists(messageID, userID, emoji)
	if err != nil {
		return false, err
	}

	if exists {
		err = r.Delete(messageID, userID, emoji)
		return false, err
	}

	reaction := &models.Reaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	}
	err = r.Create(reaction)
	return true, err
}
