package handler

import (
	"database/sql"
	"fmt"
	"forum-go/database"
	"forum-go/model"
	"forum-go/template"
	"time"

	"log"
	"net/http"
	"strconv"
)

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/viewpost" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// Extract post ID from URL
	postIDs := r.URL.Query().Get("id")

	postID, err := strconv.Atoi(postIDs)
	if err != nil {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Fetch the post
	post, err := database.FetchPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorHandler(w, r, http.StatusNotFound)
		} else {
			log.Printf("Error fetching post: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError)
		}
		return
	}

	// Fetch comments for the post
	comments, err := database.FetchCommentsByPostID(postID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		// Decide how to handle this error (continue without comments or return an error)
	}

	for i := range comments {
		comments[i].TimeAgo = calculateTimeAgo(comments[i].CreatedAt)
	}

	isLoggedIn, userID := database.CheckUserLoggedIn(r)
	var user *model.User
	if isLoggedIn {
		user, err = database.FetchUserById(userID)
		if err != nil {
			log.Printf("Error fetching user data: %v", err)
			// Continue without user data
		}
	}

	data := struct {
		*model.Post
		IsLoggedIn bool
		User       *model.User
		Comments   []model.Comment
	}{
		Post:       post,
		IsLoggedIn: isLoggedIn,
		User:       user,
		Comments:   comments,
	}

	err = template.Templates.ExecuteTemplate(w, "viewPost.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func calculateTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 365*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(duration.Hours()/24/30))
	default:
		return fmt.Sprintf("%d years ago", int(duration.Hours()/24/365))
	}
}
