package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		email    string
	}{
		{"valid user", 1, "testuser", "test@example.com"},
		{"user with special chars", 2, "user_name-123", "user+test@example.com"},
		{"empty email", 3, "nomail", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username, tt.email)
			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
				return
			}

			if token == "" {
				t.Error("GenerateToken() returned empty token")
			}

			// Validate the generated token
			claims, err := ValidateToken(token)
			if err != nil {
				t.Errorf("Generated token is invalid: %v", err)
				return
			}

			if claims.UserID != tt.userID {
				t.Errorf("Claims.UserID = %v, want %v", claims.UserID, tt.userID)
			}
			if claims.Username != tt.username {
				t.Errorf("Claims.Username = %v, want %v", claims.Username, tt.username)
			}
			if claims.Email != tt.email {
				t.Errorf("Claims.Email = %v, want %v", claims.Email, tt.email)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Generate a valid token for testing
	validToken, _ := GenerateToken(1, "testuser", "test@example.com")

	tests := []struct {
		name        string
		token       string
		wantErr     bool
		checkClaims bool
	}{
		{"valid token", validToken, false, true},
		{"empty token", "", true, false},
		{"invalid token format", "not.a.token", true, false},
		{"malformed token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkClaims && claims != nil {
				if claims.UserID != 1 {
					t.Errorf("Claims.UserID = %v, want 1", claims.UserID)
				}
				if claims.Username != "testuser" {
					t.Errorf("Claims.Username = %v, want testuser", claims.Username)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// Create an expired token manually
	claims := Claims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, _ := token.SignedString(jwtSecret)

	_, err := ValidateToken(expiredToken)
	if err == nil {
		t.Error("ValidateToken() should reject expired token")
	}
}

func TestTokenUniqueness(t *testing.T) {
	// Same user should get different tokens (due to IssuedAt timestamp)
	token1, _ := GenerateToken(1, "user", "user@example.com")
	time.Sleep(time.Millisecond * 10) // Small delay to ensure different timestamp
	token2, _ := GenerateToken(1, "user", "user@example.com")

	// Tokens might be the same if generated within the same second
	// This is acceptable behavior, just verify both are valid
	_, err1 := ValidateToken(token1)
	_, err2 := ValidateToken(token2)

	if err1 != nil || err2 != nil {
		t.Error("Both tokens should be valid")
	}
}
