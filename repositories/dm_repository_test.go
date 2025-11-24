package repositories

import (
	"GoChatApp/models"
	"testing"
)

func TestDMRepository_FindOrCreateConversation(t *testing.T) {
	db := setupTestDB(t)
	dmRepo := NewDMRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	// Create conversation
	conv, err := dmRepo.FindOrCreateConversation(user1.ID, user2.ID)
	if err != nil {
		t.Errorf("FindOrCreateConversation() error = %v", err)
	}
	if conv.ID == 0 {
		t.Error("FindOrCreateConversation() did not create conversation")
	}

	// Find existing conversation (order reversed)
	conv2, err := dmRepo.FindOrCreateConversation(user2.ID, user1.ID)
	if err != nil {
		t.Errorf("FindOrCreateConversation() error = %v", err)
	}
	if conv2.ID != conv.ID {
		t.Error("FindOrCreateConversation() should return same conversation regardless of order")
	}
}

func TestDMRepository_GetUserConversations(t *testing.T) {
	db := setupTestDB(t)
	dmRepo := NewDMRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	user3 := &models.User{Username: "user3", Email: "user3@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)
	userRepo.Create(user3)

	// Create conversations
	dmRepo.FindOrCreateConversation(user1.ID, user2.ID)
	dmRepo.FindOrCreateConversation(user1.ID, user3.ID)

	convs, err := dmRepo.GetUserConversations(user1.ID)
	if err != nil {
		t.Errorf("GetUserConversations() error = %v", err)
	}
	if len(convs) != 2 {
		t.Errorf("GetUserConversations() returned %d conversations, want 2", len(convs))
	}
}

func TestDMRepository_CreateAndGetMessages(t *testing.T) {
	db := setupTestDB(t)
	dmRepo := NewDMRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	conv, _ := dmRepo.FindOrCreateConversation(user1.ID, user2.ID)

	// Create messages
	msg1 := &models.DirectMessage{ConversationID: conv.ID, SenderID: user1.ID, Content: "Hello"}
	msg2 := &models.DirectMessage{ConversationID: conv.ID, SenderID: user2.ID, Content: "Hi there"}
	dmRepo.CreateMessage(msg1)
	dmRepo.CreateMessage(msg2)

	messages, err := dmRepo.GetMessages(conv.ID, 10, 0)
	if err != nil {
		t.Errorf("GetMessages() error = %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("GetMessages() returned %d messages, want 2", len(messages))
	}
}

func TestDMRepository_MarkAsRead(t *testing.T) {
	db := setupTestDB(t)
	dmRepo := NewDMRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	conv, _ := dmRepo.FindOrCreateConversation(user1.ID, user2.ID)

	// Create unread message from user1
	msg := &models.DirectMessage{ConversationID: conv.ID, SenderID: user1.ID, Content: "Hello", Read: false}
	dmRepo.CreateMessage(msg)

	// Mark as read by user2
	err := dmRepo.MarkAsRead(conv.ID, user2.ID)
	if err != nil {
		t.Errorf("MarkAsRead() error = %v", err)
	}

	// Verify message is read
	messages, _ := dmRepo.GetMessages(conv.ID, 10, 0)
	if !messages[0].Read {
		t.Error("MarkAsRead() should mark messages as read")
	}
}

func TestDMRepository_IsParticipant(t *testing.T) {
	db := setupTestDB(t)
	dmRepo := NewDMRepository(db)
	userRepo := NewUserRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	user3 := &models.User{Username: "user3", Email: "user3@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)
	userRepo.Create(user3)

	conv, _ := dmRepo.FindOrCreateConversation(user1.ID, user2.ID)

	// user1 is participant
	isParticipant, _ := dmRepo.IsParticipant(conv.ID, user1.ID)
	if !isParticipant {
		t.Error("IsParticipant() should return true for participant")
	}

	// user3 is not participant
	isParticipant, _ = dmRepo.IsParticipant(conv.ID, user3.ID)
	if isParticipant {
		t.Error("IsParticipant() should return false for non-participant")
	}
}
