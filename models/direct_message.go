package models

import (
	"time"

	"gorm.io/gorm"
)

// Conversation represents a DM thread between two users
type Conversation struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	User1ID   uint           `json:"user1_id" gorm:"not null;index"`
	User1     User           `json:"user1" gorm:"foreignKey:User1ID"`
	User2ID   uint           `json:"user2_id" gorm:"not null;index"`
	User2     User           `json:"user2" gorm:"foreignKey:User2ID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// DirectMessage represents a message in a DM conversation
type DirectMessage struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ConversationID uint           `json:"conversation_id" gorm:"not null;index"`
	Conversation   Conversation   `json:"-" gorm:"foreignKey:ConversationID"`
	SenderID       uint           `json:"sender_id" gorm:"not null"`
	Sender         User           `json:"sender" gorm:"foreignKey:SenderID"`
	Content        string         `json:"content" gorm:"not null"`
	Read           bool           `json:"read" gorm:"default:false"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
