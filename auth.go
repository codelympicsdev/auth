package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func auth(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	rawToken := s.Values["user_"+clientID]

	if rawToken == nil {
		http.Redirect(w, r, "/signin?client_id="+url.QueryEscape(clientID), 303)
		return
	}

	token, err := ParseToken(rawToken.(string))
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}
	if token.ExpirationTime.Before(time.Now()) {
		s.Values["user_"+clientID] = nil
		s.Save(r, w)
		http.Redirect(w, r, "/signin?client_id="+url.QueryEscape(clientID), 303)
		return
	}

	http.Redirect(w, r, "/approve?client_id="+clientID, 303)
}
