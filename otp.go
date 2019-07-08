package main

import (
	"html/template"
	"log"
	"net/http"
)

var otpTemplate = template.Must(template.ParseFiles("static/layout.html", "static/otp.html"))

func otp(w http.ResponseWriter, r *http.Request) {

}

func otpPage(w http.ResponseWriter, r *http.Request) {
	err := otpTemplate.Execute(w, map[string]interface{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}
