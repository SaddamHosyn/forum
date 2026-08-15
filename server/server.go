package server

import (
	"database/sql"
	"fmt"
	"forum-go/handler"
	"forum-go/middleware"
	"forum-go/render"
	"log"
	"net/http"
	"os"
)

func Startserver(db *sql.DB) {

	render.InitTemplates()
	RegisterServer(db)

}

func RegisterServer(db *sql.DB) {

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/favicon.ico", handler.FaviconHandler)
	http.HandleFunc("/viewpost", handler.ViewPostHandler)
	http.HandleFunc("/newpost", handler.NewPostHandler)

	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/profile", handler.ProfileHandler)

	http.HandleFunc("/submit-post", middleware.SessionMiddleware(handler.SubmitPostHandler))
	http.HandleFunc("/submitComment", middleware.SessionMiddleware(handler.SubmitCommentHandler))

	http.HandleFunc("/vote", handler.VoteHandler)
	http.HandleFunc("/vote-comment", handler.VoteCommentHandler)

	http.HandleFunc("/logout", handler.LogoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8999"
	}

	fmt.Printf("Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.EnableCORS(http.DefaultServeMux)))
}
