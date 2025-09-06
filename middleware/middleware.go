package middleware

import (
	"context"
	"database/sql"
	"forum-go/database"
	"net/http"
	"time"
)

var DB *sql.DB

func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		var userID int
		var expiresAt time.Time

		err = database.DB.QueryRow("SELECT user_id, session_expiry FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now(),
			})
			next.ServeHTTP(w, r)
			return
		}

		// Add user information to the request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
