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
	"user.basic": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Read access to full name, email and avatar",
	},
	"user.read": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Read access to full name, email and avatar and 2FA status",
	},
	"user": Scope{
		Icon:        "user",
		Name:        "User",
		Description: "Read and write access to full name, email, avatar and 2FA status",
	},
	"auth": Scope{
		Icon:        "auth",
		Name:        "Auth",
		Description: "Change your password and enable/disable 2FA",
	},
	"challenge.attempt.read": Scope{
		Icon:        "challenge",
		Name:        "Challenge Attempt",
		Description: "Read access to your past challenge attempts",
	},
	"challenge.attempt.write": Scope{
		Icon:        "challenge",
		Name:        "Challenge Attempt",
		Description: "Submit new challenge attempts",
	},
	"challenge.attempt": Scope{
		Icon:        "challenge",
		Name:        "Challenge Attempt",
		Description: "Submit new challenge attempts access your past challenge attempts",
	},
	"challenge": Scope{
		Icon:        "challenge",
		Name:        "Challenge",
		Description: "Submit new challenge attempts access your past challenge attempts",
	},
	"admin.user": Scope{
		Icon:        "admin",
		Name:        "User Admin",
		Description: "Access all user data on the platform",
	},
	"admin.attempts": Scope{
		Icon:        "admin",
		Name:        "Challenge Attempts Admin",
		Description: "Access all challenge attempts on the platform",
	},
	"admin.challenges": Scope{
		Icon:        "admin",
		Name:        "Challenge Admin",
		Description: "Access all challenges on the platform",
	},
	"admin": Scope{
		Icon:        "admin",
		Name:        "Admin",
		Description: "Do everything. Be everyone.",
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
