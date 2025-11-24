package repositories

import (
	"GoChatApp/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReadReceiptRepository struct {
	db *gorm.DB
}

func NewReadReceiptRepository(db *gorm.DB) *ReadReceiptRepository {
	return &ReadReceiptRepository{db: db}
}

// MarkAsRead updates or creates a read receipt
func (r *ReadReceiptRepository) MarkAsRead(userID, roomID, messageID uint) error {
	receipt := models.ReadReceipt{
		UserID:        userID,
		RoomID:        roomID,
		LastMessageID: messageID,
		LastReadAt:    time.Now(),
	}

	// Upsert - update if exists, insert if not
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "room_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_message_id", "last_read_at"}),
	}).Create(&receipt).Error
}

// GetLastRead gets the last read message ID for a user in a room
func (r *ReadReceiptRepository) GetLastRead(userID, roomID uint) (*models.ReadReceipt, error) {
	var receipt models.ReadReceipt
	err := r.db.Where("user_id = ? AND room_id = ?", userID, roomID).First(&receipt).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &receipt, err
}

// GetRoomReadReceipts gets all read receipts for a room
func (r *ReadReceiptRepository) GetRoomReadReceipts(roomID uint) ([]models.ReadReceipt, error) {
	var receipts []models.ReadReceipt
	err := r.db.Where("room_id = ?", roomID).Preload("User").Find(&receipts).Error
	return receipts, err
}

// GetUnreadCount gets count of unread messages for a user in a room
func (r *ReadReceiptRepository) GetUnreadCount(userID, roomID uint, messageRepo *MessageRepository) (int64, error) {
	receipt, err := r.GetLastRead(userID, roomID)
	if err != nil {
		return 0, err
	}

	var count int64
	query := r.db.Model(&models.Message{}).Where("room_id = ? AND deleted = ?", roomID, false)

	if receipt != nil {
		query = query.Where("id > ?", receipt.LastMessageID)
	}

	err = query.Count(&count).Error
	return count, err
}
