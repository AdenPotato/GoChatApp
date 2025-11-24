package main

import (
	"GoChatApp/database"
	"GoChatApp/handlers"
	"GoChatApp/repositories"
	"GoChatApp/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	//schizo function for database flushing
	//database.FlushDB()

	// Seed database (optional - for development)
	// Uncomment the line below to seed the database with sample data
	// database.SeedDB()

	// Get database instance
	db := database.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	roomRepo := repositories.NewRoomRepository(db)
	reactionRepo := repositories.NewReactionRepository(db)
	dmRepo := repositories.NewDMRepository(db)
	blockRepo := repositories.NewBlockRepository(db)
	receiptRepo := repositories.NewReadReceiptRepository(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	roomHandler := handlers.NewRoomHandler(roomRepo)
	reactionHandler := handlers.NewReactionHandler(reactionRepo)
	dmHandler := handlers.NewDMHandler(dmRepo)
	blockHandler := handlers.NewBlockHandler(blockRepo)
	receiptHandler := handlers.NewReadReceiptHandler(receiptRepo)
	uploadHandler := handlers.NewUploadHandler()
	wsHandler := handlers.NewWebSocketHandler()

	// Setup router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, authHandler, userHandler, messageHandler, roomHandler,
		reactionHandler, dmHandler, blockHandler, receiptHandler, uploadHandler, wsHandler)

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
