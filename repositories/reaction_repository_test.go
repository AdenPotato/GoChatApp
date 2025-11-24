package repositories

import (
	"GoChatApp/models"
	"testing"
)

func TestReactionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	// Setup
	user := &models.User{Username: "reactor", Email: "reactor@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Test message"}
	msgRepo.Create(msg)

	// Test
	reaction := &models.Reaction{
		MessageID: msg.ID,
		UserID:    user.ID,
		Emoji:     "👍",
	}

	err := reactionRepo.Create(reaction)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if reaction.ID == 0 {
		t.Error("Create() did not set reaction ID")
	}
}

func TestReactionRepository_Toggle(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	// Setup
	user := &models.User{Username: "toggler", Email: "toggler@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Test message"}
	msgRepo.Create(msg)

	// First toggle - should add
	added, err := reactionRepo.Toggle(msg.ID, user.ID, "❤️")
	if err != nil {
		t.Errorf("Toggle() error = %v", err)
	}
	if !added {
		t.Error("Toggle() first call should add reaction")
	}

	// Second toggle - should remove
	added, err = reactionRepo.Toggle(msg.ID, user.ID, "❤️")
	if err != nil {
		t.Errorf("Toggle() error = %v", err)
	}
	if added {
		t.Error("Toggle() second call should remove reaction")
	}
}

func TestReactionRepository_GetReactionCounts(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	// Setup
	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg := &models.Message{UserID: user1.ID, RoomID: room.ID, Content: "Test message"}
	msgRepo.Create(msg)

	// Add reactions
	reactionRepo.Create(&models.Reaction{MessageID: msg.ID, UserID: user1.ID, Emoji: "👍"})
	reactionRepo.Create(&models.Reaction{MessageID: msg.ID, UserID: user2.ID, Emoji: "👍"})
	reactionRepo.Create(&models.Reaction{MessageID: msg.ID, UserID: user1.ID, Emoji: "❤️"})

	counts, err := reactionRepo.GetReactionCounts(msg.ID)
	if err != nil {
		t.Errorf("GetReactionCounts() error = %v", err)
	}

	if counts["👍"] != 2 {
		t.Errorf("GetReactionCounts() 👍 = %d, want 2", counts["👍"])
	}
	if counts["❤️"] != 1 {
		t.Errorf("GetReactionCounts() ❤️ = %d, want 1", counts["❤️"])
	}
}

func TestReactionRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	reactionRepo := NewReactionRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	user := &models.User{Username: "user", Email: "user@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Test message"}
	msgRepo.Create(msg)

	// Before adding
	exists, _ := reactionRepo.Exists(msg.ID, user.ID, "👍")
	if exists {
		t.Error("Exists() should return false before adding reaction")
	}

	// After adding
	reactionRepo.Create(&models.Reaction{MessageID: msg.ID, UserID: user.ID, Emoji: "👍"})
	exists, _ = reactionRepo.Exists(msg.ID, user.ID, "👍")
	if !exists {
		t.Error("Exists() should return true after adding reaction")
	}
}
