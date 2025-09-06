package handler

import (
	"forum-go/database"
	"log"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if database.DB == nil {
		log.Println("Database connection is nil in LogoutHandler")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No active session", http.StatusBadRequest)
			return
		}
		log.Printf("Error retrieving cookie: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sessionToken := cookie.Value

	_, err = database.DB.Exec("DELETE FROM sessions WHERE session_token = ?", sessionToken)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
