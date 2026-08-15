package tests

import (
	"forum-go/pkg/utils"
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
			err := utils.ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}
