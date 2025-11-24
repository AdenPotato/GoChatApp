package repositories

import (
	"GoChatApp/models"

	"gorm.io/gorm"
)

type BlockRepository struct {
	db *gorm.DB
}

func NewBlockRepository(db *gorm.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

// Block blocks a user
func (r *BlockRepository) Block(blockerID, blockedID uint) error {
	block := &models.Block{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}
	return r.db.Create(block).Error
}

// Unblock removes a block
func (r *BlockRepository) Unblock(blockerID, blockedID uint) error {
	return r.db.Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		Delete(&models.Block{}).Error
}

// IsBlocked checks if blocker has blocked blocked
func (r *BlockRepository) IsBlocked(blockerID, blockedID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Block{}).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		Count(&count).Error
	return count > 0, err
}

// IsBlockedEither checks if either user has blocked the other
func (r *BlockRepository) IsBlockedEither(userID1, userID2 uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Block{}).
		Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)",
			userID1, userID2, userID2, userID1).
		Count(&count).Error
	return count > 0, err
}

// GetBlockedUsers gets all users blocked by a user
func (r *BlockRepository) GetBlockedUsers(blockerID uint) ([]models.Block, error) {
	var blocks []models.Block
	err := r.db.Where("blocker_id = ?", blockerID).Preload("Blocked").Find(&blocks).Error
	return blocks, err
}

// GetBlockedByUsers gets all users who have blocked this user
func (r *BlockRepository) GetBlockedByUsers(blockedID uint) ([]uint, error) {
	var blockerIDs []uint
	err := r.db.Model(&models.Block{}).
		Where("blocked_id = ?", blockedID).
		Pluck("blocker_id", &blockerIDs).Error
	return blockerIDs, err
}
