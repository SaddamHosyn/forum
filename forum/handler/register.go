package handler

import (
	"database/sql"
	"fmt"
	"forum-go/auth"
	"forum-go/database"
	"forum-go/pkg/utils"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Validate input (you should add more thorough validation, below the complete validation logic)
	if username == "" || email == "" || password == "" {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// **Validate input using `ValidateInputs()`**
	if err := utils.ValidateInputs(database.DB, username, email, password); err != nil {
		log.Printf("Validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUserID, err := auth.UserExists(database.DB, username)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking user existence: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if existingUserID != "" {
		ErrorHandler(w, r, http.StatusConflict)
		return
	}

	// Add the user to the database
	err = auth.AddUser(database.DB,username, email, password)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")
	log.Printf("User [ %v ] registered successfully.", username)
}
