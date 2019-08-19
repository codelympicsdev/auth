package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var approveTemplate = template.Must(template.ParseFiles("static/layout.html", "static/approve.html"))

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
		renderError(w, "Client ID Missing", "The client ID is missing. You should specify a client id in the request.", "")
		return
	}

	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	token := s.Values["user_"+clientID]
	if token == nil {
		http.Redirect(w, r, "/signin?client_id="+url.QueryEscape(clientID)+"&redirect="+url.QueryEscape("/approve?client_id="+clientID), 303)
		return
	}

	client, err := getAPIClient(clientID)
	if err != nil {
		log.Fatalln(err.Error())
		internalServerError(w)
	}

	http.Redirect(w, r, client.RedirectURL+"?token="+url.QueryEscape(token.(string)), 303)
}

func approvePage(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Client ID Missing", "The client ID is missing. You should specify a client id in the request.", "")
		return
	}

	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	if s.Values["user_"+clientID] == nil {
		http.Redirect(w, r, "/signin?client_id="+url.QueryEscape(clientID)+"&redirect="+url.QueryEscape("/approve?client_id="+clientID), 303)
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
