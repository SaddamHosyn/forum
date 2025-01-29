package main

import (
	"log"
	"net/http"
)

func main() {

	RegisterServer()

	log.Println("Server is starting...\n")
	log.Println("Go on http://localhost:8000/\n")
	log.Println("To shut down the server press CTRL + C\n")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

func RegisterServer() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post/", postHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./html"))))
}
