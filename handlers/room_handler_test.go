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
)

func TestRoomHandler_GetRooms(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	handler := NewRoomHandler(roomRepo)

	// Create some rooms
	roomRepo.Create(&models.Room{Name: "Room 1", Type: "public"})
	roomRepo.Create(&models.Room{Name: "Room 2", Type: "private"})

	router := gin.New()
	router.GET("/rooms", handler.GetRooms)

	req := httptest.NewRequest("GET", "/rooms", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string][]models.Room
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response["rooms"]) != 2 {
		t.Errorf("Expected 2 rooms, got %d", len(response["rooms"]))
	}
}

func TestRoomHandler_GetRoom(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	handler := NewRoomHandler(roomRepo)

	room := &models.Room{Name: "Test Room", Type: "public"}
	roomRepo.Create(room)

	router := gin.New()
	router.GET("/rooms/:id", handler.GetRoom)

	req := httptest.NewRequest("GET", "/rooms/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRoomHandler_GetRoom_NotFound(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	handler := NewRoomHandler(roomRepo)

	router := gin.New()
	router.GET("/rooms/:id", handler.GetRoom)

	req := httptest.NewRequest("GET", "/rooms/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestRoomHandler_CreateRoom_Success(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	handler := NewRoomHandler(roomRepo)

	router := gin.New()
	router.POST("/rooms", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		handler.CreateRoom(c)
	})

	body := map[string]string{
		"name": "New Room",
		"type": "public",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestRoomHandler_CreateRoom_Unauthorized(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	handler := NewRoomHandler(roomRepo)

	router := gin.New()
	router.POST("/rooms", handler.CreateRoom) // No user_id in context

	body := map[string]string{
		"name": "New Room",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRoomHandler_JoinRoom_Success(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	userRepo := repositories.NewUserRepository(db)
	handler := NewRoomHandler(roomRepo)

	// Create user and room
	user := &models.User{Username: "joiner", Email: "join@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Join Room", Type: "public"}
	roomRepo.Create(room)

	router := gin.New()
	router.POST("/rooms/:id/join", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		handler.JoinRoom(c)
	})

	req := httptest.NewRequest("POST", "/rooms/1/join", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify membership
	isMember, _ := roomRepo.IsMember(room.ID, user.ID)
	if !isMember {
		t.Error("User should be a member after joining")
	}
}

func TestRoomHandler_JoinRoom_AlreadyMember(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	userRepo := repositories.NewUserRepository(db)
	handler := NewRoomHandler(roomRepo)

	user := &models.User{Username: "member", Email: "member@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Join Room", Type: "public"}
	roomRepo.Create(room)
	roomRepo.AddMember(room.ID, user.ID)

	router := gin.New()
	router.POST("/rooms/:id/join", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		handler.JoinRoom(c)
	})

	req := httptest.NewRequest("POST", "/rooms/1/join", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

func TestRoomHandler_LeaveRoom_Success(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	userRepo := repositories.NewUserRepository(db)
	handler := NewRoomHandler(roomRepo)

	user := &models.User{Username: "leaver", Email: "leave@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Leave Room", Type: "public"}
	roomRepo.Create(room)
	roomRepo.AddMember(room.ID, user.ID)

	router := gin.New()
	router.POST("/rooms/:id/leave", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		handler.LeaveRoom(c)
	})

	req := httptest.NewRequest("POST", "/rooms/1/leave", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify not a member
	isMember, _ := roomRepo.IsMember(room.ID, user.ID)
	if isMember {
		t.Error("User should not be a member after leaving")
	}
}

func TestRoomHandler_LeaveRoom_NotMember(t *testing.T) {
	db := setupTestDB(t)
	roomRepo := repositories.NewRoomRepository(db)
	userRepo := repositories.NewUserRepository(db)
	handler := NewRoomHandler(roomRepo)

	user := &models.User{Username: "nonmember", Email: "non@example.com", PasswordHash: "hash"}
	userRepo.Create(user)

	room := &models.Room{Name: "Leave Room", Type: "public"}
	roomRepo.Create(room)

	router := gin.New()
	router.POST("/rooms/:id/leave", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		handler.LeaveRoom(c)
	})

	req := httptest.NewRequest("POST", "/rooms/1/leave", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
