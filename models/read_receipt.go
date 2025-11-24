package models

import (
	"time"
)

// ReadReceipt tracks when a user has read messages in a room
type ReadReceipt struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_room"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	RoomID        uint      `json:"room_id" gorm:"not null;uniqueIndex:idx_user_room"`
	Room          Room      `json:"room" gorm:"foreignKey:RoomID"`
	LastMessageID uint      `json:"last_message_id" gorm:"not null"` // Last message ID user has read
	LastReadAt    time.Time `json:"last_read_at"`
}
