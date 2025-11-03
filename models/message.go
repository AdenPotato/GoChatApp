package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	RoomID    uint           `json:"room_id" gorm:"not null"`
	Room      Room           `json:"room" gorm:"foreignKey:RoomID"`
	Content   string         `json:"content" gorm:"not null"`
	Edited    bool           `json:"edited" gorm:"default:false"`
	Deleted   bool           `json:"deleted" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
