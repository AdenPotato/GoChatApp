package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null;default:'public'"` // public, private, direct
	CreatedBy uint           `json:"created_by"`
	Creator   User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	Members   []User         `json:"members" gorm:"many2many:room_members;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
