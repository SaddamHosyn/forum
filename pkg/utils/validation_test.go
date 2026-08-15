package utils

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid Password", "Password123!", false},
		{"Too Short", "Pass1!", true},
		{"Missing Uppercase", "password123!", true},
		{"Missing Lowercase", "PASSWORD123!", true},
		{"Missing Digit", "Password!", true},
		{"Missing Special Char", "Password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"Valid Alpha", "john", true},
		{"Valid Alphanumeric", "john123", true},
		{"Valid Underscore and Dash", "john_doe-99", true},
		{"Invalid Special Char", "john@doe", false},
		{"Invalid Space", "john doe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUsername(tt.username); got != tt.want {
				t.Errorf("isValidUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}
