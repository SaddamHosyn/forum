package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateInputs(DB *sql.DB, username, email, password string) error {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if username == "" || email == "" || password == "" {
		return ValidationError{Field: "general", Message: "all fields are required"}
	}

	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "invalid email format"}
	}

	if len(username) < 5 || len(username) > 15 {
		return ValidationError{Field: "username", Message: "username must be between 5 and 30 characters long"}
	}

	if err := ValidatePassword(password); err != nil {
		return ValidationError{Field: "password", Message: err.Error()}
	}

	if !isValidUsername(username) {
		return ValidationError{Field: "username", Message: "username can only contain letters, numbers, underscores, and dashes"}
	}

	usernameAvailable, err := UsernameNotTaken(DB, username)
	if err != nil {
		return fmt.Errorf("error checking username availability: %w", err)
	}
	if !usernameAvailable {
		return ValidationError{Field: "username", Message: "username already taken"}
	}

	emailAvailable, err := EmailNotTaken(DB, email)
	if err != nil {
		return fmt.Errorf("error checking email availability: %w", err)
	}
	if !emailAvailable {
		return ValidationError{Field: "email", Message: "email already registered"}
	}

	return nil
}

func isValidUsername(username string) bool {
	for _, char := range username {
		if !(isValidCharacter(char) || char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func isValidCharacter(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func EmailNotTaken(DB *sql.DB, email string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking email: %w", err)
	}
	return count == 0, nil
}

func UsernameNotTaken(DB *sql.DB, username string) (bool, error) {

	if DB == nil {
		return false, fmt.Errorf("database connection is nil123")
	}
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking username: %w", err)
	}
	return count == 0, nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	lowercase := regexp.MustCompile(`[a-z]`)
	uppercase := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`[0-9]`)
	specialChar := regexp.MustCompile(`[@$!%*?&]`)

	if !lowercase.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !uppercase.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !digit.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	if !specialChar.MatchString(password) {
		return errors.New("password must contain at least one special character (@, $, !, %, *, ?, &)")
	}

	return nil
}
