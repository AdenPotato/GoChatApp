package handlers

import (
	"GoChatApp/models"
	"GoChatApp/repositories"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Room{}, &models.Message{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestAuthHandler_Register_Success(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	handler := NewAuthHandler(userRepo)

	router := gin.New()
	router.POST("/register", handler.Register)

	body := map[string]string{
		"username": "newuser",
		"email":    "new@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["token"] == nil {
		t.Error("Response should contain token")
	}
	if response["user"] == nil {
		t.Error("Response should contain user")
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	handler := NewAuthHandler(userRepo)

	// Create existing user
	userRepo.Create(&models.User{
		Username:     "existinguser",
		Email:        "existing@example.com",
		PasswordHash: "hash",
	})

	router := gin.New()
	router.POST("/register", handler.Register)

	body := map[string]string{
		"username": "existinguser",
		"email":    "new@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestAuthHandler_Register_InvalidInput(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	handler := NewAuthHandler(userRepo)

	router := gin.New()
	router.POST("/register", handler.Register)

	tests := []struct {
		name string
		body map[string]string
	}{
		{"missing username", map[string]string{"email": "test@example.com", "password": "password123"}},
		{"missing email", map[string]string{"username": "user", "password": "password123"}},
		{"missing password", map[string]string{"username": "user", "email": "test@example.com"}},
		{"short password", map[string]string{"username": "user", "email": "test@example.com", "password": "12345"}},
		{"invalid email", map[string]string{"username": "user", "email": "notanemail", "password": "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	handler := NewAuthHandler(userRepo)

	// First register a user
	router := gin.New()
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Register
	regBody := map[string]string{
		"username": "loginuser",
		"email":    "login@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(regBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Now login
	loginBody := map[string]string{
		"username": "loginuser",
		"password": "password123",
	}
	jsonBody, _ = json.Marshal(loginBody)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["token"] == nil {
		t.Error("Response should contain token")
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	handler := NewAuthHandler(userRepo)

	router := gin.New()
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	// Register
	regBody := map[string]string{
		"username": "loginuser",
		"email":    "login@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(regBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Login with wrong password
	loginBody := map[string]string{
		"username": "loginuser",
		"password": "wrongpassword",
	}
	jsonBody, _ = json.Marshal(loginBody)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Login_NonexistentUser(t *testing.T) {
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	handler := NewAuthHandler(userRepo)

	router := gin.New()
	router.POST("/login", handler.Login)

	loginBody := map[string]string{
		"username": "nonexistent",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(loginBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
