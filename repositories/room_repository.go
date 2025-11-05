package repositories

import (
	"GoChatApp/models"

	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create creates a new room
func (r *RoomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

// FindByID finds a room by ID
func (r *RoomRepository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	err := r.db.Preload("Creator").Preload("Members").First(&room, id).Error
	return &room, err
}

// FindAll returns all rooms
func (r *RoomRepository) FindAll() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Preload("Creator").Preload("Members").Find(&rooms).Error
	return rooms, err
}

// Update updates a room
func (r *RoomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

// Delete deletes a room
func (r *RoomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

// AddMember adds a user to a room
func (r *RoomRepository) AddMember(roomID, userID uint) error {
	return r.db.Exec("INSERT INTO room_members (room_id, user_id) VALUES (?, ?)", roomID, userID).Error
}

// RemoveMember removes a user from a room
func (r *RoomRepository) RemoveMember(roomID, userID uint) error {
	return r.db.Exec("DELETE FROM room_members WHERE room_id = ? AND user_id = ?", roomID, userID).Error
}

// IsMember checks if a user is a member of a room
func (r *RoomRepository) IsMember(roomID, userID uint) (bool, error) {
	var count int64
	err := r.db.Table("room_members").Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count).Error
	return count > 0, err
}
