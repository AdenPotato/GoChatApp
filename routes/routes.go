package routes

import (
	"GoChatApp/handlers"
	"GoChatApp/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, messageHandler *handlers.MessageHandler, wsHandler *handlers.WebSocketHandler) {
	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Serve static files
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./index.html")
	})

	// API routes
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Authentication routes
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		// User routes
		api.GET("/users", userHandler.GetUsers)
		api.GET("/users/:id", userHandler.GetUserByID)

		// Message routes
		api.GET("/messages", messageHandler.GetMessages)
		api.POST("/messages", messageHandler.SendMessage)

		// WebSocket endpoint for real-time chat
		api.GET("/ws", wsHandler.HandleWebSocket)
	}
}
