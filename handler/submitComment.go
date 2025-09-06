package handler

import (
	"forum-go/database"
	"net/http"
)

func SubmitCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		ErrorHandler(w, r, http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if postID == "" || content == "" {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/viewpost?id="+postID, http.StatusSeeOther)
}
