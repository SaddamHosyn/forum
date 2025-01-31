package server

import (
	"forum-go/handler"
	"forum-go/util"
	"log"
	"net/http"
)

func StartServer() {

	RegisterServer()

	log.Println("Server is starting...\n")
	log.Println("Go on http://localhost:8000/\n")
	log.Println("To shut down the server press CTRL + C\n")

	err := http.ListenAndServe(":8000", nil)
	util.CheckErr("Server failed to start: %v", err)
}

func RegisterServer() {

	fs := http.FileServer(http.Dir("public"))

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/post/", handler.PostHandler)
	http.HandleFunc("/login", handler.LoginHandler)

	http.HandleFunc("/action", handler.LoginHandler)
	http.HandleFunc("/scifi", handler.LoginHandler)
	http.HandleFunc("/romance", handler.LoginHandler)
	http.HandleFunc("/horror", handler.LoginHandler)
	http.HandleFunc("/romance", handler.LoginHandler)
	http.HandleFunc("/romance", handler.LoginHandler)



	http.Handle("/public/", http.StripPrefix("/public/", fs))

}
