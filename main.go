package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"gopkg.in/boj/redistore.v1"
)

var apiURL = os.Getenv("API_URL")
var rootClientID = os.Getenv("ROOT_CLIENT_ID")
var rootClientSecret = os.Getenv("ROOT_CLIENT_SECRET")
var redisURI = os.Getenv("REDIS_URI")
var redisPassword = os.Getenv("REDIS_PASSWORD")
var rs *redistore.RediStore

func main() {
	if apiURL == "" {
		apiURL = "https://api.codelympics.dev/v0"
	}

	var err error
	rs, err = redistore.NewRediStore(10, "tcp", redisURI, redisPassword, []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		panic(err)
	}
	defer rs.Close()

	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("https://codelympics.dev", 303))

	r.HandleFunc("/signin", signinHandler)
	r.HandleFunc("/otp", otpHandler)
	r.HandleFunc("/approve", approveHandler)

	r.HandleFunc("/oauth2/auth", authHandler)
	r.HandleFunc("/oauth2/token", tokenHandler)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalln(err)
	}
}
