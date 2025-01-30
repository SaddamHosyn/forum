package main

import (
	"log"
	"net/http"
	"text/template"
)

type Text struct {
	ErrorNum int
	ErrorMes string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("./html/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error 500: Rendering Error on about.html")
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("./html/post.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error 500: Rendering Error on about.html")
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		t, err := template.ParseFiles("./html/error.html")
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		em := "HTTP status 404: Page Not Found"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusInternalServerError {
		t, err := template.ParseFiles("./html/error.html")
		if err != nil {
			log.Printf("HTTP status 500: Internal Server Error -missing error.html file%v", err)
		}
		em := "HTTP status 500: Internal Server Error"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusBadRequest {
		t, err := template.ParseFiles("./html/error.html")
		if err != nil {
			log.Printf("HTTP status 500: Internal Server Error -missing error.html file%v", err)
		}
		em := "HTTP status 400: Bad Request\nPlease select artist from the Home Page"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
}
