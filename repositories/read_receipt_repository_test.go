package repositories

import (
	"GoChatApp/models"
	"testing"
)

func TestReadReceiptRepository_MarkAsRead(t *testing.T) {
	db := setupTestDB(t)
	receiptRepo := NewReadReceiptRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	user := &models.User{Username: "reader", Email: "reader@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Test"}
	msgRepo.Create(msg)

	err := receiptRepo.MarkAsRead(user.ID, room.ID, msg.ID)
	if err != nil {
		t.Errorf("MarkAsRead() error = %v", err)
	}

	receipt, err := receiptRepo.GetLastRead(user.ID, room.ID)
	if err != nil {
		t.Errorf("GetLastRead() error = %v", err)
	}
	if receipt == nil {
		t.Error("MarkAsRead() should create receipt")
	}
	if receipt.LastMessageID != msg.ID {
		t.Errorf("MarkAsRead() LastMessageID = %d, want %d", receipt.LastMessageID, msg.ID)
	}
}

func TestReadReceiptRepository_MarkAsRead_Upsert(t *testing.T) {
	db := setupTestDB(t)
	receiptRepo := NewReadReceiptRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	user := &models.User{Username: "reader", Email: "reader@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg1 := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Test 1"}
	msg2 := &models.Message{UserID: user.ID, RoomID: room.ID, Content: "Test 2"}
	msgRepo.Create(msg1)
	msgRepo.Create(msg2)

	// Mark first message
	receiptRepo.MarkAsRead(user.ID, room.ID, msg1.ID)

	// Mark second message (should update, not create new)
	receiptRepo.MarkAsRead(user.ID, room.ID, msg2.ID)

	receipt, _ := receiptRepo.GetLastRead(user.ID, room.ID)
	if receipt.LastMessageID != msg2.ID {
		t.Errorf("MarkAsRead() should update to latest message, got %d want %d", receipt.LastMessageID, msg2.ID)
	}

	// Verify only one receipt exists
	receipts, _ := receiptRepo.GetRoomReadReceipts(room.ID)
	if len(receipts) != 1 {
		t.Errorf("Should have exactly 1 receipt, got %d", len(receipts))
	}
}

func TestReadReceiptRepository_GetLastRead_NotFound(t *testing.T) {
	db := setupTestDB(t)
	receiptRepo := NewReadReceiptRepository(db)

	receipt, err := receiptRepo.GetLastRead(999, 999)
	if err != nil {
		t.Errorf("GetLastRead() should not error for not found, got %v", err)
	}
	if receipt != nil {
		t.Error("GetLastRead() should return nil for not found")
	}
}

func TestReadReceiptRepository_GetRoomReadReceipts(t *testing.T) {
	db := setupTestDB(t)
	receiptRepo := NewReadReceiptRepository(db)
	userRepo := NewUserRepository(db)
	roomRepo := NewRoomRepository(db)
	msgRepo := NewMessageRepository(db)

	user1 := &models.User{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"}
	user2 := &models.User{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"}
	userRepo.Create(user1)
	userRepo.Create(user2)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	msg := &models.Message{UserID: user1.ID, RoomID: room.ID, Content: "Test"}
	msgRepo.Create(msg)

	receiptRepo.MarkAsRead(user1.ID, room.ID, msg.ID)
	receiptRepo.MarkAsRead(user2.ID, room.ID, msg.ID)

	receipts, err := receiptRepo.GetRoomReadReceipts(room.ID)
	if err != nil {
		t.Errorf("GetRoomReadReceipts() error = %v", err)
	}
	if len(receipts) != 2 {
		t.Errorf("GetRoomReadReceipts() returned %d receipts, want 2", len(receipts))
	}
}
