package render

import (
	"html/template"
	"log"
)

var Templates *template.Template

func InitTemplates() {
	var err error
	Templates, err = template.ParseFiles(
		"./templates/index.html",
		"./templates/header.html",
		"./templates/footer.html",
		"./templates/error.html",
		"./templates/newPost.html",
		"./templates/viewPost.html",
		"./templates/profile.html",
	)
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
}
