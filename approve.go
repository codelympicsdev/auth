package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)

var approveTemplate = template.Must(template.ParseFiles("static/layout.html", "static/approve.html"))

// Scope is a localized scope for the approve screen
type Scope struct {
	Icon        string
	Name        string
	Description string
}

var scopes = map[string]Scope{
	"user": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Full name, email and avatar",
	},
}

func approve(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
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

	client, err := getAPIClient(clientID)
	if err != nil {
		log.Fatalln(err.Error())
		internalServerError(w)
	}

	http.Redirect(w, r, client.RedirectURL+"?token="+url.QueryEscape(rawToken.(string)), 303)
}

func approvePage(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
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

	client, err := getAPIClient(clientID)
	if err != nil {
		log.Fatalln(err.Error())
		internalServerError(w)
		return
	}

	if client.Trusted {
		approve(w, r)
		return
	}

	var localizedScopes []Scope

	for _, id := range client.Scopes {
		if scope, ok := scopes[id]; ok {
			localizedScopes = append(localizedScopes, scope)
		} else {
			localizedScopes = append(localizedScopes, Scope{
				Icon:        "unknown",
				Name:        "Unknown Permission",
				Description: id,
			})
		}
	}

	err = approveTemplate.Execute(w, map[string]interface{}{
		"id":          client.ID,
		"name":        client.Name,
		"scopes":      localizedScopes,
		"redirectURL": client.RedirectURL,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}
