package handler

import (
	"fmt"
	"log"
	"net/http"
)

func SubmitPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		ErrorHandler(w, r, http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["category"]

	// Validate Categories

	if len(categories) > 3 || len(categories) > 0 {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Validate of Title/Content FormValue

	if title == "" || content == "" {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	result, err := DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Get the last inserted post ID
	postID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Insert each selected category into a junction table (e.g., post_categories)
	for _, category := range categories {
		_, err := DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, category)
		if err != nil {
			log.Printf("Error linking category %s to post: %v", category, err)
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
	redirectURL := fmt.Sprintf("/viewpost?id=%d", postID)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
