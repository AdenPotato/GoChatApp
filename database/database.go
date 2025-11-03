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
