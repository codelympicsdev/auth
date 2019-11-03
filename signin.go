package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var signinTemplate = template.Must(template.ParseFiles("static/layout.html", "static/signin.html"))

// SigninRequest is the request for a jwt token with email and password
type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

// AuthResponse is the response with a token
type AuthResponse struct {
	Token string `json:"token"`
}

func doSignin(email string, password string, clientID string) (*AuthResponse, error) {
	sendData, err := json.Marshal(SigninRequest{Email: email, Password: password, ClientID: clientID})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", apiURL+"/auth/signin", bytes.NewBuffer(sendData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", rootClientID+" "+rootClientSecret)
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.StatusCode != 200 {
		if resp != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			s := buf.String()
			return nil, errors.New("error or non 200 status: " + s)
		}

		return nil, errors.New("error or non 200 status")
	}
	var data *AuthResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func signin(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	if email == "" || password == "" {
		r.Form.Add("error", "Please enter email and password.")
		signinPage(w, r)
		return
	}

	resp, err := doSignin(email, password, clientID)
	if err != nil {
		log.Println(err.Error())
		r.Form.Add("error", "Failed to sign in.")
		signinPage(w, r)
		return
	}

	s, err := store.Get(r, "session")
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	token, err := ParseToken(resp.Token)
	if err != nil {
		log.Println(err.Error())
		internalServerError(w)
		return
	}

	if token.RequiresUpgrade {
		s.Values["user_otpUpgrade_"+clientID] = resp.Token
		s.Save(r, w)

		http.Redirect(w, r, "/otp?client_id="+url.QueryEscape(clientID), 303)
	} else {
		s.Values["user_"+clientID] = resp.Token
		err = s.Save(r, w)
		if err != nil {
			log.Println(err.Error())
			internalServerError(w)
			return
		}

		http.Redirect(w, r, "/approve?client_id="+url.QueryEscape(clientID), 303)
	}
}

func signinPage(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	e := r.FormValue("error")

	err := signinTemplate.Execute(w, map[string]interface{}{
		"query": "?client_id=" + url.QueryEscape(clientID),
		"error": e,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}
