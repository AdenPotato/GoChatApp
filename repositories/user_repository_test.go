package repositories

import (
	"GoChatApp/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations for all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Message{},
		&models.Reaction{},
		&models.Conversation{},
		&models.DirectMessage{},
		&models.Block{},
		&models.ReadReceipt{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	if user.ID == 0 {
		t.Error("Create() did not set user ID")
	}
}

func TestUserRepository_Create_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user1 := &models.User{
		Username:     "testuser",
		Email:        "test1@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user1)

	user2 := &models.User{
		Username:     "testuser", // Duplicate
		Email:        "test2@example.com",
		PasswordHash: "hashedpassword",
	}

	err := repo.Create(user2)
	if err == nil {
		t.Error("Create() should fail for duplicate username")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a user
	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	// Find by ID
	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
		return
	}

	if found.Username != user.Username {
		t.Errorf("FindByID() Username = %v, want %v", found.Username, user.Username)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByID(999)
	if err == nil {
		t.Error("FindByID() should return error for non-existent user")
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByUsername("testuser")
	if err != nil {
		t.Errorf("FindByUsername() error = %v", err)
		return
	}

	if found.ID != user.ID {
		t.Errorf("FindByUsername() ID = %v, want %v", found.ID, user.ID)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	found, err := repo.FindByEmail("test@example.com")
	if err != nil {
		t.Errorf("FindByEmail() error = %v", err)
		return
	}

	if found.ID != user.ID {
		t.Errorf("FindByEmail() ID = %v, want %v", found.ID, user.ID)
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create multiple users
	users := []*models.User{
		{Username: "user1", Email: "user1@example.com", PasswordHash: "hash"},
		{Username: "user2", Email: "user2@example.com", PasswordHash: "hash"},
		{Username: "user3", Email: "user3@example.com", PasswordHash: "hash"},
	}

	for _, u := range users {
		repo.Create(u)
	}

	found, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
		return
	}

	if len(found) != 3 {
		t.Errorf("FindAll() returned %v users, want 3", len(found))
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	// Update user
	user.Avatar = "avatar.png"
	err := repo.Update(user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	// Verify update
	found, _ := repo.FindByID(user.ID)
	if found.Avatar != "avatar.png" {
		t.Errorf("Update() Avatar = %v, want avatar.png", found.Avatar)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	repo.Create(user)

	err := repo.Delete(user.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	// Verify deletion (soft delete)
	_, err = repo.FindByID(user.ID)
	if err == nil {
		t.Error("Delete() user should not be findable after deletion")
	}
}
