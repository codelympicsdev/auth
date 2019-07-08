package main

import (
	"html/template"
	"log"
	"net/http"
)

var errorTemplate = template.Must(template.ParseFiles("static/layout.html", "static/error.html"))

func renderError(w http.ResponseWriter, title string, message string, returnURL string) {
	err := errorTemplate.Execute(w, map[string]interface{}{
		"title":     title,
		"message":   message,
		"returnURL": returnURL,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}

func internalServerError(w http.ResponseWriter) {
	renderError(w, "Internal Server Error", "An internal server error has occured.", "")
}
