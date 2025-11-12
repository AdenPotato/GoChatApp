package database

import (
	"GoChatApp/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	// Run migrations
	Migrate()
}

// Migrate runs database migrations
func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Message{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// FlushDB deletes all data from all tables
func FlushDB() error {
	log.Println("Flushing database...")

	// Delete all records from each table
	if err := DB.Exec("DELETE FROM messages").Error; err != nil {
		log.Printf("Error deleting messages: %v", err)
		return err
	}

	if err := DB.Exec("DELETE FROM rooms").Error; err != nil {
		log.Printf("Error deleting rooms: %v", err)
		return err
	}

	if err := DB.Exec("DELETE FROM users").Error; err != nil {
		log.Printf("Error deleting users: %v", err)
		return err
	}

	// Reset SQLite auto-increment counters
	DB.Exec("DELETE FROM sqlite_sequence WHERE name='messages'")
	DB.Exec("DELETE FROM sqlite_sequence WHERE name='rooms'")
	DB.Exec("DELETE FROM sqlite_sequence WHERE name='users'")

	log.Println("Database flushed successfully")
	return nil
}
