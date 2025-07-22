package handler

import (
	"fmt"
	"forum-go/database"
	"forum-go/model"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/newpost" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	isLoggedIn, userID := database.CheckUserLoggedIn(r)
	if !isLoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := database.FetchUserById(userID)
	if err != nil {
		log.Printf("Error fetching user data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:

		categories, err := database.FetchCategories()
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Categories []model.Category
			IsLoggedIn bool
			User       *model.User
		}{
			Categories: categories,
			IsLoggedIn: isLoggedIn,
			User:       user,
		}

		// Display the new post form
		tmpl, err := template.ParseFiles(
			"./assets/html/newPost.html",
			"./assets/html/header.html",
			"./assets/html/footer.html")
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		// Process the form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["category"]

		if title == "" || content == "" || len(categories) == 0 {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		categoriesStr := strings.Join(categories, ", ")

		// Create a new post
		post := &model.Post{
			Title:      title,
			Content:    content,
			UserID:     userID,
			Categories: categoriesStr,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Save the post to the database
		postID, err := savePost(post)
		if err != nil {
			log.Printf("Error saving post: %v", err)
			http.Error(w, "Error saving post", http.StatusInternalServerError)
			return
		}

		// Redirect to the new post
		http.Redirect(w, r, fmt.Sprintf("/viewpost?id=%d", postID), http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func savePost(post *model.Post) (int64, error) {
	query := `
        INSERT INTO posts (title, content, user_id, categories, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	result, err := database.DB.Exec(query, post.Title, post.Content, post.UserID, post.Categories, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("error saving post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return id, nil
}
