package main

import (
	"html/template"
	"log"
	"net/http"
)

var signinTemplate = template.Must(template.ParseFiles("static/layout.html", "static/signin.html"))

func signin(w http.ResponseWriter, r *http.Request) {

}

func signinPage(w http.ResponseWriter, r *http.Request) {
	err := signinTemplate.Execute(w, map[string]interface{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}
