package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// GetAPIClientResponse is the response for getting a single api client
type GetAPIClientResponse struct {
	ID          string   `json:"id"`
	RedirectURL string   `json:"redirect_url"`
	Trusted     bool     `json:"trusted"`
	Scopes      []string `json:"scopes"`
}

func auth(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	resp, err := http.Get("https://api.codelympics.dev/v0/apiclients/" + clientID)
	if err != nil || resp.StatusCode != 200 {
		renderError(w, "Authentication Error", "The provided api client does not exist.", "")
		return
	}

	var data *GetAPIClientResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	stringifiedData, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	s.Values["apiclient"] = string(stringifiedData)

	err = s.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	if s.Values["user"] == nil {
		http.Redirect(w, r, "/signin", 303)
		return
	}

	http.Redirect(w, r, "/approve", 303)
}
