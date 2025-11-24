package repositories

import (
	"GoChatApp/models"
	"testing"
)

func TestMessageRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	// Create user and room
	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}
	roomRepo.Create(room)

	// Create message
	message := &models.Message{
		UserID:  user.ID,
		RoomID:  room.ID,
		Content: "Hello, World!",
	}

	err := msgRepo.Create(message)
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	if message.ID == 0 {
		t.Error("Create() did not set message ID")
	}
}

func TestMessageRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}
	roomRepo.Create(room)

	message := &models.Message{
		UserID:  user.ID,
		RoomID:  room.ID,
		Content: "Test message",
	}
	msgRepo.Create(message)

	// Find by ID
	found, err := msgRepo.FindByID(message.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
		return
	}

	if found.Content != message.Content {
		t.Errorf("FindByID() Content = %v, want %v", found.Content, message.Content)
	}

	// Should preload User
	if found.User.Username != user.Username {
		t.Errorf("FindByID() should preload User, got %v", found.User.Username)
	}
}

func TestMessageRepository_FindByRoomID(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room1 := &models.Room{Name: "Room 1", Type: "public"}
	room2 := &models.Room{Name: "Room 2", Type: "public"}
	roomRepo.Create(room1)
	roomRepo.Create(room2)

	// Create messages in different rooms
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room1.ID, Content: "Room 1 - Message 1"})
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room1.ID, Content: "Room 1 - Message 2"})
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room2.ID, Content: "Room 2 - Message 1"})

	// Find messages in room1
	messages, err := msgRepo.FindByRoomID(room1.ID, 10, 0)
	if err != nil {
		t.Errorf("FindByRoomID() error = %v", err)
		return
	}

	if len(messages) != 2 {
		t.Errorf("FindByRoomID() returned %v messages, want 2", len(messages))
	}
}

func TestMessageRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	// Create messages
	for i := 0; i < 5; i++ {
		msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room.ID, Content: "Message"})
	}

	// Test pagination
	messages, err := msgRepo.FindAll(3, 0)
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
		return
	}

	if len(messages) != 3 {
		t.Errorf("FindAll() with limit 3 returned %v messages", len(messages))
	}
}

func TestMessageRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	message := &models.Message{
		UserID:  user.ID,
		RoomID:  room.ID,
		Content: "To be deleted",
	}
	msgRepo.Create(message)

	// Soft delete
	err := msgRepo.Delete(message.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	// Message should still exist but be marked as deleted
	found, _ := msgRepo.FindByID(message.ID)
	if !found.Deleted {
		t.Error("Delete() should mark message as deleted")
	}
}

func TestMessageRepository_FindAll_ExcludesDeleted(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	// Create messages
	msg1 := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Message 1"}
	msg2 := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Message 2"}
	msgRepo.Create(msg1)
	msgRepo.Create(msg2)

	// Delete one message
	msgRepo.Delete(msg1.ID)

	// FindAll should exclude deleted messages
	messages, _ := msgRepo.FindAll(10, 0)
	if len(messages) != 1 {
		t.Errorf("FindAll() should exclude deleted messages, got %v", len(messages))
	}
}

func TestMessageRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{
		Username:     "sender",
		Email:        "sender@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	message := &models.Message{
		UserID:  user.ID,
		RoomID:  room.ID,
		Content: "Original content",
	}
	msgRepo.Create(message)

	// Update message
	message.Content = "Updated content"
	message.Edited = true
	err := msgRepo.Update(message)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	// Verify update
	found, _ := msgRepo.FindByID(message.ID)
	if found.Content != "Updated content" {
		t.Errorf("Update() Content = %v, want 'Updated content'", found.Content)
	}
	if !found.Edited {
		t.Error("Update() Edited should be true")
	}
}

func TestMessageRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	msgRepo := NewMessageRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)

	user := &models.User{Username: "sender", Email: "sender@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room1 := &models.Room{Name: "Room 1", Type: "public"}
	room2 := &models.Room{Name: "Room 2", Type: "public"}
	roomRepo.Create(room1)
	roomRepo.Create(room2)

	// Create messages with searchable content
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room1.ID, Content: "Hello world"})
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room1.ID, Content: "Hello there"})
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room2.ID, Content: "Hello again"})
	msgRepo.Create(&models.Message{UserID: user.ID, RoomID: room1.ID, Content: "Goodbye"})

	// Search all rooms
	messages, err := msgRepo.Search("Hello", 0, 10, 0)
	if err != nil {
		t.Errorf("Search() error = %v", err)
	}
	if len(messages) != 3 {
		t.Errorf("Search() returned %d messages, want 3", len(messages))
	}

	// Search specific room
	messages, _ = msgRepo.Search("Hello", room1.ID, 10, 0)
	if len(messages) != 2 {
		t.Errorf("Search() with room filter returned %d messages, want 2", len(messages))
	}

	// Search with no results
	messages, _ = msgRepo.Search("nonexistent", 0, 10, 0)
	if len(messages) != 0 {
		t.Errorf("Search() should return empty for no matches, got %d", len(messages))
	}
}
