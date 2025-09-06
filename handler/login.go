package handler

import (
	"forum-go/auth"
	"forum-go/database"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if database.DB == nil {
		log.Printf("Error: database connection is nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	submittedUsername := r.FormValue("username")
	submittedPassword := r.FormValue("password")

	user, err := auth.GetUserInfo(database.DB, submittedUsername)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return 
		} else {
			log.Printf("Error retrieving user info: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}




	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(submittedPassword)); err != nil {
		log.Printf("Login failed (wrong password) for %s at %s\n", submittedUsername, time.Now().Format("2006-01-02 15:04:05"))
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create session
	if err := database.CreateSession(w, user.ID); err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	log.Printf("Login successful: User '%s' logged in at %s",
		user.Username,
		time.Now().Format("2006-01-02 15:04:05"))
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful")) // Add this line to send a success message
}
