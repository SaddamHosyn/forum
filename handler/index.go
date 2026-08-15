package handler

import (
	"forum-go/database"
	"forum-go/model"
	"forum-go/render"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	category := r.URL.Query().Get("category")

	var posts []model.Post
	var err error

	if category != "" && category != "All Movies" { // "All Movies" should return all posts
		// Fetch posts filtered by category
		posts, err = database.FetchPostsByCategory(category)
		if err != nil {
			log.Printf("Error fetching posts by category: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// Fetch all posts
		posts, err = database.FetchPosts()
		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	categories, err := database.FetchCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	success := r.URL.Query().Get("success") == "true"

	isLoggedIn, userID := database.CheckUserLoggedIn(r)
	var user *model.User
	if isLoggedIn {
		user, err = database.FetchUserById(userID)
		if err != nil {
			log.Printf("Error fetching user data: %v", err)
			// Continue without user data
		}
	}

	data := model.HomePageData{
		Posts:      posts,
		Categories: categories,
		User:       user,
		Success:    success,
		IsLoggedIn: isLoggedIn,
	}

	err = render.Templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
