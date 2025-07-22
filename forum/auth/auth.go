package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"forum-go/model"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func AddUser(db *sql.DB, username, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(username, email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error executing statement: %w", err)
	}

	return nil
}

func CheckCookie(db *sql.DB, cookie string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE cookie = ?", cookie).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking cookie: %w", err)
	}
	return count > 0, nil
}

func GetExpires(db *sql.DB, cookie string) (string, error) {
	var expires string
	err := db.QueryRow("SELECT expires FROM users WHERE cookie = ?", cookie).Scan(&expires)
	if err != nil {
		return "", fmt.Errorf("error getting expiration: %w", err)
	}
	return expires, nil
}

func Logout(db *sql.DB, username string) error {
	_, err := db.Exec("UPDATE users SET cookie = '', expires = '' WHERE username = ?", username)
	if err != nil {
		return fmt.Errorf("error logging out user: %w", err)
	}
	fmt.Printf("User %s logged out successfully\n", username)
	return nil
}

func UpdateCookie(db *sql.DB, username, newToken, expiration string) error {
	_, err := db.Exec("UPDATE users SET session_token = ?, session_expiry = ? WHERE username = ?", newToken, expiration, username)
	if err != nil {
		return fmt.Errorf("error updating cookie: %w", err)
	}
	fmt.Printf("Cookie updated for user: %s\n", username)
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPassword), nil
}

func GetUserInfo(db *sql.DB, submittedUsername string) (*model.User, error) {
	var user model.User
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	err := db.QueryRow("SELECT id, username, email, password_hash FROM users WHERE username = ?", submittedUsername).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func UserExists(db *sql.DB, username string) (string, error) {
	var userID string
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
