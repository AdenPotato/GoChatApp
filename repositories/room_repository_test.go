package repositories

import (
	"GoChatApp/models"
	"testing"
)

func TestRoomRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoomRepository(db)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}

	err := repo.Create(room)
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	if room.ID == 0 {
		t.Error("Create() did not set room ID")
	}
}

func TestRoomRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := NewRoomRepository(db)
	userRepo := NewUserRepository(db)

	// Create a user first
	user := &models.User{
		Username:     "creator",
		Email:        "creator@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	// Create a room
	room := &models.Room{
		Name:      "Test Room",
		Type:      "public",
		CreatedBy: user.ID,
	}
	roomRepo.Create(room)

	// Find by ID
	found, err := roomRepo.FindByID(room.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
		return
	}

	if found.Name != room.Name {
		t.Errorf("FindByID() Name = %v, want %v", found.Name, room.Name)
	}
}

func TestRoomRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoomRepository(db)

	// Create multiple rooms
	rooms := []*models.Room{
		{Name: "Room 1", Type: "public"},
		{Name: "Room 2", Type: "private"},
		{Name: "Room 3", Type: "public"},
	}

	for _, r := range rooms {
		repo.Create(r)
	}

	found, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
		return
	}

	if len(found) != 3 {
		t.Errorf("FindAll() returned %v rooms, want 3", len(found))
	}
}

func TestRoomRepository_AddMember(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := NewRoomRepository(db)
	userRepo := NewUserRepository(db)

	// Create user and room
	user := &models.User{
		Username:     "member",
		Email:        "member@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}
	roomRepo.Create(room)

	// Add member
	err := roomRepo.AddMember(room.ID, user.ID)
	if err != nil {
		t.Errorf("AddMember() error = %v", err)
		return
	}

	// Check membership
	isMember, err := roomRepo.IsMember(room.ID, user.ID)
	if err != nil {
		t.Errorf("IsMember() error = %v", err)
		return
	}

	if !isMember {
		t.Error("AddMember() user should be a member after adding")
	}
}

func TestRoomRepository_RemoveMember(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := NewRoomRepository(db)
	userRepo := NewUserRepository(db)

	// Create user and room
	user := &models.User{
		Username:     "member",
		Email:        "member@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}
	roomRepo.Create(room)

	// Add then remove member
	roomRepo.AddMember(room.ID, user.ID)
	err := roomRepo.RemoveMember(room.ID, user.ID)
	if err != nil {
		t.Errorf("RemoveMember() error = %v", err)
		return
	}

	// Check membership
	isMember, _ := roomRepo.IsMember(room.ID, user.ID)
	if isMember {
		t.Error("RemoveMember() user should not be a member after removal")
	}
}

func TestRoomRepository_IsMember_NotMember(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := NewRoomRepository(db)
	userRepo := NewUserRepository(db)

	user := &models.User{
		Username:     "nonmember",
		Email:        "nonmember@example.com",
		PasswordHash: "hash",
	}
	userRepo.Create(user)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}
	roomRepo.Create(room)

	isMember, err := roomRepo.IsMember(room.ID, user.ID)
	if err != nil {
		t.Errorf("IsMember() error = %v", err)
		return
	}

	if isMember {
		t.Error("IsMember() should return false for non-member")
	}
}

func TestRoomRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoomRepository(db)

	room := &models.Room{
		Name: "Test Room",
		Type: "public",
	}
	repo.Create(room)

	err := repo.Delete(room.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	_, err = repo.FindByID(room.ID)
	if err == nil {
		t.Error("Delete() room should not be findable after deletion")
	}
}
