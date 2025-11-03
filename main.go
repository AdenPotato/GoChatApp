package main

import (
	"GoChatApp/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	// Seed database (optional - for development)
	// Uncomment the line below to seed the database with sample data
	// database.SeedDB()

	router := gin.Default()

	// CORS middleware for React frontend
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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

		// Chat routes
		api.GET("/messages", getMessages)
		api.POST("/messages", sendMessage)

		// User routes
		api.GET("/users", getUsers)
		api.POST("/login", login)
		api.POST("/register", register)

		// WebSocket endpoint for real-time chat
		api.GET("/ws", handleWebSocket)
	}

	router.Run(":8080")
}

// Handler functions
func getMessages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"messages": []string{}})
}

func sendMessage(c *gin.Context) {
	var message struct {
		Content string `json:"content"`
		User    string `json:"user"`
	}

	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent"})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": []string{}})
}

func login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": "sample_token"})
}

func register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func handleWebSocket(c *gin.Context) {
	// WebSocket upgrade logic will go here
	c.JSON(http.StatusNotImplemented, gin.H{"message": "WebSocket not implemented yet"})
}
