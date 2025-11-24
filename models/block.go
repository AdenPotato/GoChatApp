package models

import (
	"time"

	"gorm.io/gorm"
)

// Block represents a user blocking another user
type Block struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	BlockerID   uint           `json:"blocker_id" gorm:"not null;index"`
	Blocker     User           `json:"blocker" gorm:"foreignKey:BlockerID"`
	BlockedID   uint           `json:"blocked_id" gorm:"not null;index"`
	Blocked     User           `json:"blocked" gorm:"foreignKey:BlockedID"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
