package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)

var signinTemplate = template.Must(template.ParseFiles("static/layout.html", "static/signin.html"))
var otpTemplate = template.Must(template.ParseFiles("static/layout.html", "static/otp.html"))
var approveTemplate = template.Must(template.ParseFiles("static/layout.html", "static/approve.html"))

func signinHandler(w http.ResponseWriter, r *http.Request) {
	store, err := rs.Get(r, "codelympics-auth")
	if err != nil {
		internalServerError(w)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" || password == "" {
			r.Form.Add("error", "Please enter email and password.")
			signinPage(w, r)
			return
		}

		userID, requires2FA, err := signinEmailPassword(email, password)
		if err != nil {
			r.Form.Add("error", "Signin failed.")
			signinPage(w, r)
			return
		}
		if requires2FA {
			store.Values["requires_2fa_user_id"] = userID
			store.Values["requires_2fa_timeout"] = time.Now().Add(10 * time.Minute).Unix()
			err := store.Save(r, w)
			if err != nil {
				internalServerError(w)
				return
			}
			w.Header().Set("Location", "/otp")
			w.WriteHeader(http.StatusFound)
		} else {
			store.Values["signed_in_user_id"] = userID
			err := store.Save(r, w)
			if err != nil {
				internalServerError(w)
				return
			}
			w.Header().Set("Location", "/approve")
			w.WriteHeader(http.StatusFound)
		}
		return
	}

	signinPage(w, r)
}

func signinPage(w http.ResponseWriter, r *http.Request) {
	e := r.FormValue("error")

	err := signinTemplate.Execute(w, map[string]interface{}{
		"error": e,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}

func otpHandler(w http.ResponseWriter, r *http.Request) {
	store, err := rs.Get(r, "codelympics-auth")
	if err != nil {
		internalServerError(w)
		return
	}

	v, ok := store.Values["requires_2fa_user_id"]
	if !ok {
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID := v.(string)

	if timeout, ok := store.Values["requires_2fa_timeout"]; !ok || time.Now().After(time.Unix(timeout.(int64), 0)) {
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusFound)
		return
	}

	if r.Method == "POST" {
		otp := r.FormValue("otp")

		if otp == "" {
			r.Form.Add("error", "Please enter a one time password.")
			otpPage(w, r)
			return
		}

		valid, err := checkOTP(userID, otp)
		if err != nil {
			r.Form.Add("error", "Validation failed.")
			otpPage(w, r)
			return
		}
		if !valid {
			r.Form.Add("error", "One time password is not valid.")
			otpPage(w, r)
			return
		}

		delete(store.Values, "requires_2fa_user_id")
		delete(store.Values, "requires_2fa_timeout")
		store.Values["signed_in_user_id"] = userID
		err = store.Save(r, w)
		if err != nil {
			internalServerError(w)
			return
		}
		w.Header().Set("Location", "/approve")
		w.WriteHeader(http.StatusFound)

		return
	}

	otpPage(w, r)
}

func otpPage(w http.ResponseWriter, r *http.Request) {
	e := r.FormValue("error")

	err := otpTemplate.Execute(w, map[string]interface{}{
		"error": e,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}

func approveHandler(w http.ResponseWriter, r *http.Request) {
	store, err := rs.Get(r, "codelympics-auth")
	if err != nil {
		internalServerError(w)
		return
	}

	if _, ok := store.Values["signed_in_user_id"]; !ok {
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusFound)
		return
	}

	if r.Method == "POST" {
		w.Header().Set("Location", "/oauth2/auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	var form url.Values
	if v, ok := store.Values["return_form"]; ok {
		form = v.(url.Values)
	}

	clientID := form.Get("client_id")
	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	client, err := getAPIClient(clientID)
	if err != nil {
		internalServerError(w)
		return
	}

	if client.Trusted {
		w.Header().Set("Location", "/oauth2/auth")
		w.WriteHeader(http.StatusFound)
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
		"redirectURL": form.Get("redirect_uri"),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}
