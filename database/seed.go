package database

import (
	"GoChatApp/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedDB seeds the database with sample data for development
func SeedDB() {
	// Check if data already exists
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already seeded, skipping...")
		return
	}

	// Create sample users
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users := []models.User{
		{
			Username:     "alice",
			Email:        "alice@example.com",
			PasswordHash: string(password),
			Avatar:       "https://i.pravatar.cc/150?img=1",
		},
		{
			Username:     "bob",
			Email:        "bob@example.com",
			PasswordHash: string(password),
			Avatar:       "https://i.pravatar.cc/150?img=2",
		},
		{
			Username:     "charlie",
			Email:        "charlie@example.com",
			PasswordHash: string(password),
			Avatar:       "https://i.pravatar.cc/150?img=3",
		},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Error seeding user %s: %v", user.Username, err)
		}
	}

	// Create sample rooms
	rooms := []models.Room{
		{
			Name:      "General",
			Type:      "public",
			CreatedBy: 1,
		},
		{
			Name:      "Random",
			Type:      "public",
			CreatedBy: 1,
		},
		{
			Name:      "Tech Talk",
			Type:      "public",
			CreatedBy: 2,
		},
	}

	for _, room := range rooms {
		if err := DB.Create(&room).Error; err != nil {
			log.Printf("Error seeding room %s: %v", room.Name, err)
		}
	}

	// Add members to rooms
	var generalRoom models.Room
	DB.Where("name = ?", "General").First(&generalRoom)
	DB.Model(&generalRoom).Association("Members").Append(&users)

	// Create sample messages
	messages := []models.Message{
		{
			UserID:  1,
			RoomID:  1,
			Content: "Welcome to the chat app!",
		},
		{
			UserID:  2,
			RoomID:  1,
			Content: "Thanks! Excited to be here!",
		},
		{
			UserID:  3,
			RoomID:  1,
			Content: "Hey everyone!",
		},
	}

	for _, message := range messages {
		if err := DB.Create(&message).Error; err != nil {
			log.Printf("Error seeding message: %v", err)
		}
	}

	log.Println("Database seeded successfully!")
}
