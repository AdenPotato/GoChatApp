package repositories

import (
	"GoChatApp/models"
	"testing"
)

func TestBlockRepository_Block(t *testing.T) {
	db := setupTestDB(t)
	blockRepo := NewBlockRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "blocker", Email: "blocker@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "blocked", Email: "blocked@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	err := blockRepo.Block(user1.ID, user2.ID)
	if err != nil {
		t.Errorf("Block() error = %v", err)
	}

	isBlocked, _ := blockRepo.IsBlocked(user1.ID, user2.ID)
	if !isBlocked {
		t.Error("Block() should create block relationship")
	}
}

func TestBlockRepository_Unblock(t *testing.T) {
	db := setupTestDB(t)
	blockRepo := NewBlockRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "blocker", Email: "blocker@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "blocked", Email: "blocked@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	blockRepo.Block(user1.ID, user2.ID)
	err := blockRepo.Unblock(user1.ID, user2.ID)
	if err != nil {
		t.Errorf("Unblock() error = %v", err)
	}

	isBlocked, _ := blockRepo.IsBlocked(user1.ID, user2.ID)
	if isBlocked {
		t.Error("Unblock() should remove block relationship")
	}
}

func TestBlockRepository_IsBlocked(t *testing.T) {
	db := setupTestDB(t)
	blockRepo := NewBlockRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	// Before blocking
	isBlocked, _ := blockRepo.IsBlocked(user1.ID, user2.ID)
	if isBlocked {
		t.Error("IsBlocked() should return false before blocking")
	}

	// After blocking
	blockRepo.Block(user1.ID, user2.ID)
	isBlocked, _ = blockRepo.IsBlocked(user1.ID, user2.ID)
	if !isBlocked {
		t.Error("IsBlocked() should return true after blocking")
	}

	// Reverse direction should be false
	isBlocked, _ = blockRepo.IsBlocked(user2.ID, user1.ID)
	if isBlocked {
		t.Error("IsBlocked() should be directional")
	}
}

func TestBlockRepository_IsBlockedEither(t *testing.T) {
	db := setupTestDB(t)
	blockRepo := NewBlockRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	// user1 blocks user2
	blockRepo.Block(user1.ID, user2.ID)

	// Both directions should return true
	isBlocked, _ := blockRepo.IsBlockedEither(user1.ID, user2.ID)
	if !isBlocked {
		t.Error("IsBlockedEither() should return true when user1 blocked user2")
	}

	isBlocked, _ = blockRepo.IsBlockedEither(user2.ID, user1.ID)
	if !isBlocked {
		t.Error("IsBlockedEither() should return true when checked in reverse")
	}
}

func TestBlockRepository_GetBlockedUsers(t *testing.T) {
	db := setupTestDB(t)
	blockRepo := NewBlockRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "blocker", Email: "blocker@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "blocked1", Email: "blocked1@example.com", PasswordHash: "hash"}
	user3 := &models.User{Username: "blocked2", Email: "blocked2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)
	userRepo.Create(user3)

	blockRepo.Block(user1.ID, user2.ID)
	blockRepo.Block(user1.ID, user3.ID)

	blocks, err := blockRepo.GetBlockedUsers(user1.ID)
	if err != nil {
		t.Errorf("GetBlockedUsers() error = %v", err)
	}
	if len(blocks) != 2 {
		t.Errorf("GetBlockedUsers() returned %d blocks, want 2", len(blocks))
	}
}
