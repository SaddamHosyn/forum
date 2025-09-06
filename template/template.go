package template

import (
	"html/template"
	"log"
)

var Templates *template.Template

func InitTemplates() {
	var err error
	Templates, err = template.ParseFiles(
		"./assets/html/index.html",
		"./assets/html/header.html",
		"./assets/html/footer.html",
		"./assets/html/error.html",
		"./assets/html/newPost.html",
		"./assets/html/viewPost.html",
		"./assets/html/profile.html",
	)
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
}
