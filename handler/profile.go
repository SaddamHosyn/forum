package handler

import (
	"fmt"
	"forum-go/database"
	"forum-go/model"
	"forum-go/render"
	"log"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	loggedIn, userID := database.CheckUserLoggedIn(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch user information
	user, err := database.FetchUserById(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Fetch user's posts
	posts, err := database.FetchPostsByUserID(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// Fetch liked posts
	likedPosts, err := database.FetchLikedPostsByUserID(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	dislikedPosts, err := database.FetchDislikedPostsByUserID(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	//fmt.Printf("User: %s, LikedPosts Count: %d\n", user.Username, len(likedPosts))

	// Prepare data for the template
	data := struct {
		Title         string
		User          *model.User
		Posts         []*model.Post
		LikedPosts    []*model.Post
		DislikedPosts []*model.Post
		IsLoggedIn    bool
	}{
		Title:         fmt.Sprintf("%s's Profile", user.Username),
		User:          user,
		Posts:         posts,
		LikedPosts:    likedPosts,
		DislikedPosts: dislikedPosts,
		IsLoggedIn:    true,
	}

	err = render.Templates.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
