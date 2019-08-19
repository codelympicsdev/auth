package main

import (
	"encoding/json"
	"errors"
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
	r.HandleFunc("/approve", approvePage).Methods("GET")
	r.HandleFunc("/approve", approve).Methods("POST")

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, r)
}

// GetAPIClientResponse is the response for getting a single api client
type GetAPIClientResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	RedirectURL string   `json:"redirect_url"`
	Trusted     bool     `json:"trusted"`
	Scopes      []string `json:"scopes"`
}

func getAPIClient(clientID string) (*GetAPIClientResponse, error) {
	resp, err := http.Get("https://api.codelympics.dev/v0/apiclients/" + clientID)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.StatusCode != 200 {
		return nil, errors.New("error or non 200 status")
	}

	var data *GetAPIClientResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
