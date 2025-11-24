package routes

import (
	"GoChatApp/handlers"
	"GoChatApp/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, messageHandler *handlers.MessageHandler, roomHandler *handlers.RoomHandler, reactionHandler *handlers.ReactionHandler, dmHandler *handlers.DMHandler, blockHandler *handlers.BlockHandler, receiptHandler *handlers.ReadReceiptHandler, uploadHandler *handlers.UploadHandler, wsHandler *handlers.WebSocketHandler) {
	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Serve React build (production) or static test files (development)
	if _, err := os.Stat("./dist"); err == nil {
		// Serve React production build
		router.Static("/assets", "./dist/assets")
		router.StaticFile("/", "./dist/index.html")
		router.StaticFile("/favicon.ico", "./dist/favicon.ico")

		// SPA fallback - serve index.html for all non-API routes
		router.NoRoute(func(c *gin.Context) {
			// Don't serve index.html for API routes
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}
			c.File("./dist/index.html")
		})
	} else {
		// Fallback to static test files for development
		router.Static("/static", "./static")
		router.GET("/", func(c *gin.Context) {
			c.File("./index.html")
		})
	}

	// Serve uploaded files
	router.Static("/uploads", "./uploads")

	// API routes (public)
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Authentication routes (public)
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		// Public user routes
		api.GET("/users", userHandler.GetUsers)
		api.GET("/users/:id", userHandler.GetUserByID)

		// Public message routes (read only)
		api.GET("/messages", messageHandler.GetMessages)
		api.GET("/messages/search", messageHandler.SearchMessages)

		// Public room routes (read only)
		api.GET("/rooms", roomHandler.GetRooms)
		api.GET("/rooms/:id", roomHandler.GetRoom)

		// Public reaction routes (read only)
		api.GET("/messages/:id/reactions", reactionHandler.GetReactions)

		// Read receipts (public read)
		api.GET("/rooms/:id/receipts", receiptHandler.GetReadReceipts)

		// WebSocket endpoint (auth via token query param)
		api.GET("/ws", wsHandler.HandleWebSocket)
	}

	// Protected routes (require authentication)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Message routes (protected)
		protected.POST("/messages", messageHandler.SendMessage)

		// Room routes (protected)
		protected.POST("/rooms", roomHandler.CreateRoom)
		protected.POST("/rooms/:id/join", roomHandler.JoinRoom)
		protected.POST("/rooms/:id/leave", roomHandler.LeaveRoom)

		// Reaction routes (protected)
		protected.POST("/messages/:id/reactions", reactionHandler.ToggleReaction)

		// Direct message routes (protected)
		protected.GET("/conversations", dmHandler.GetConversations)
		protected.POST("/conversations", dmHandler.StartConversation)
		protected.GET("/conversations/:id/messages", dmHandler.GetMessages)
		protected.POST("/conversations/:id/messages", dmHandler.SendMessage)
		protected.GET("/conversations/unread", dmHandler.GetUnreadCount)

		// Block routes (protected)
		protected.GET("/blocks", blockHandler.GetBlockedUsers)
		protected.POST("/users/:id/block", blockHandler.BlockUser)
		protected.DELETE("/users/:id/block", blockHandler.UnblockUser)

		// Read receipt routes (protected)
		protected.POST("/receipts", receiptHandler.MarkAsRead)

		// File upload (protected)
		protected.POST("/upload", uploadHandler.UploadFile)
	}
}
