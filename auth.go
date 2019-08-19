package main

import (
	"log"
	"net/http"
	"net/url"
)

func auth(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	clientID := r.FormValue("client_id")

	approveURL := "/approve?client_id=" + clientID

	if s.Values["user_"+clientID] == nil {
		http.Redirect(w, r, "/signin?client_id="+url.QueryEscape(clientID)+"&redirect="+url.QueryEscape(approveURL), 303)
		return
	}

	http.Redirect(w, r, approveURL, 303)
}
