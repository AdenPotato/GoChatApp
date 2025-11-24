package models

import (
	"time"

	"gorm.io/gorm"
)

type Reaction struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MessageID uint           `json:"message_id" gorm:"not null;index"`
	Message   Message        `json:"-" gorm:"foreignKey:MessageID"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Emoji     string         `json:"emoji" gorm:"not null"` // e.g., "👍", "❤️", "😂"
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UniqueIndex ensures one user can only react with same emoji once per message
// Add in migration: CREATE UNIQUE INDEX idx_unique_reaction ON reactions(message_id, user_id, emoji)
