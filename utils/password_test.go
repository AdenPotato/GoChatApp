package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"simple password", "password123"},
		{"complex password", "P@ssw0rd!#$%"},
		{"empty password", ""},
		{"long password", "thisIsAVeryLongPasswordThatShouldStillWork123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword() error = %v", err)
				return
			}

			// Hash should not be empty
			if hash == "" {
				t.Error("HashPassword() returned empty hash")
			}

			// Hash should be different from original password
			if hash == tt.password {
				t.Error("HashPassword() hash equals original password")
			}

			// Same password should produce different hashes (due to salt)
			hash2, _ := HashPassword(tt.password)
			if hash == hash2 {
				t.Error("HashPassword() produced identical hashes for same password")
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testPassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		want           bool
	}{
		{"correct password", hash, password, true},
		{"wrong password", hash, "wrongPassword", false},
		{"empty password", hash, "", false},
		{"similar password", hash, "testPassword124", false},
		{"case sensitive", hash, "TESTPASSWORD123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.hashedPassword, tt.password); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPasswordWithInvalidHash(t *testing.T) {
	// Should return false for invalid hash format
	if CheckPassword("not-a-valid-hash", "password") {
		t.Error("CheckPassword() should return false for invalid hash")
	}
}
