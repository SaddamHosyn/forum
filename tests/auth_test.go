package tests

import (
	"forum-go/auth"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "SecretPassword123!"

	hashed, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashed == "" {
		t.Fatal("Expected non-empty hashed password")
	}

	if hashed == password {
		t.Fatal("Hashed password should not match plain text password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		t.Errorf("bcrypt comparison failed for valid password: %v", err)
	}
}
