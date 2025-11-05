package repositories

import (
	"GoChatApp/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message
func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// FindByID finds a message by ID
func (r *MessageRepository) FindByID(id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("User").Preload("Room").First(&message, id).Error
	return &message, err
}

// FindByRoomID finds all messages in a room with pagination
func (r *MessageRepository) FindByRoomID(roomID uint, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.
		Where("room_id = ? AND deleted = ?", roomID, false).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// FindAll returns all messages with pagination
func (r *MessageRepository) FindAll(limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.
		Where("deleted = ?", false).
		Preload("User").
		Preload("Room").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// Update updates a message
func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

// Delete soft deletes a message
func (r *MessageRepository) Delete(id uint) error {
	return r.db.Model(&models.Message{}).Where("id = ?", id).Update("deleted", true).Error
}
