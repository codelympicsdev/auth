package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/gbrlsnchs/jwt/v3"

	"github.com/gorilla/mux"

	"gopkg.in/boj/redistore.v1"
)

var rootClientID = os.Getenv("ROOT_CLIENT_ID")
var rootClientSecret = os.Getenv("ROOT_CLIENT_SECRET")
var store *redistore.RediStore

func main() {
	var redisURI = os.Getenv("REDIS_URI")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	var err error
	store, err = redistore.NewRediStore(10, "tcp", redisURI, redisPassword, []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		panic(err)
	}
	defer store.Close()

	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("https://codelympics.dev", 303))
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

// Token is the structure for the JWT token
type Token struct {
	jwt.Payload

	RequiresUpgrade bool `json:"requires_upgrade"`

	ID        string `json:"id,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`

	Scopes []string `json:"scopes"`
}

// ParseToken takes a jwt and extracts the content
func ParseToken(token string) (*Token, error) {
	var t *Token

	_, err := jwt.Verify([]byte(token), jwt.None(), &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
