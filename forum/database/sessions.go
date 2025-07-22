package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

var db *sql.DB

func CheckUserLoggedIn(r *http.Request) (bool, int) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return false, 0
	}

	var userID int
	var expiresAt time.Time

	query := `
    SELECT user_id, session_expiry 
    FROM sessions 
    WHERE session_token = ? 
    AND session_token = (
        SELECT session_token 
        FROM sessions AS s2 
        WHERE s2.user_id = sessions.user_id 
        ORDER BY session_expiry DESC 
        LIMIT 1
    )
	`

	err = DB.QueryRow(query, sessionToken.Value).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0
		}
		log.Printf("Error checking session: %v", err)
		return false, 0
	}

	if time.Now().After(expiresAt) {
		// Session has expired
		deleteSession(sessionToken.Value)
		return false, 0
	}

	return true, userID
}

func CreateSession(w http.ResponseWriter, userID int) error {
	// First, invalidate any existing session for this user
	_, err := DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		log.Printf("Error deleting old sessions: %v", err)
	}
	token, cookie, err := generateSessionToken()
	if err != nil {
		return err
	}

	query := `
    INSERT INTO sessions 
    (user_id, session_token, session_expiry) 
    VALUES (?, ?, ?)
	`

	_, err = DB.Exec(query, userID, token, cookie.Expires)
	if err != nil {
		return fmt.Errorf("failed to insert session: %v", err)
	}

	http.SetCookie(w, cookie)
	return nil
}

// -- Non-Global Functions : Only happens in this package server -- //

func deleteSession(token string) {
	_, err := DB.Exec("DELETE FROM sessions WHERE session_token = ?", token)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
	}
}

func generateSessionToken() (string, *http.Cookie, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return "", nil, err
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token.String(),
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
	}
	return token.String(), cookie, nil
}

// func invalidateUserSessions(userID int) error {
// 	_, err := db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
// 	return err
// }
