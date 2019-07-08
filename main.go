package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var rootClientID = os.Getenv("ROOT_CLIENT_ID")
var rootClientSecret = os.Getenv("ROOT_CLIENT_SECRET")
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth", auth).Methods("GET")
	r.HandleFunc("/signin", signinPage).Methods("GET")
	r.HandleFunc("/signin", signin).Methods("POST")
	r.HandleFunc("/otp", otpPage).Methods("GET")
	r.HandleFunc("/otp", otp).Methods("POST")

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, r)
}
