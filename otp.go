package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)

var otpTemplate = template.Must(template.ParseFiles("static/layout.html", "static/otp.html"))

// OPTUpgradeRequest is what is used to upgrade a token with a OTP
type OPTUpgradeRequest struct {
	Token string `json:"token"`
	OTP   string `json:"otp"`
}

func doOtp(otp string, token string, clientID string) (*AuthResponse, error) {
	sendData, err := json.Marshal(OPTUpgradeRequest{OTP: otp})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", "https://api.codelympics.dev/v0/auth/upgrade/otp", bytes.NewBuffer(sendData))
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

func otp(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	otp := r.PostFormValue("otp")

	if otp == "" {
		r.Form.Add("error", "Please enter a one time password.")
		signinPage(w, r)
		return
	}

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

	rawToken := s.Values["user_otpUpgrade_"+clientID]
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
	if !token.RequiresUpgrade {
		s.Values["user_"+clientID] = rawToken
		s.Save(r, w)
		http.Redirect(w, r, "/approve?client_id="+url.QueryEscape(clientID), 303)
		return
	}

	resp, err := doOtp(otp, rawToken.(string), clientID)
	if err != nil || resp == nil || resp.Token == "" {
		log.Println(err.Error())
		r.Form.Add("error", "Failed to sign in.")
		signinPage(w, r)
		return
	}

	s.Values["user_otpUpgrade_"+clientID] = nil
	s.Values["user_"+clientID] = resp.Token
	s.Save(r, w)

	http.Redirect(w, r, "/approve?client_id="+url.QueryEscape(clientID), 303)
}

func otpPage(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")

	if clientID == "" {
		renderError(w, "Incomplete request", "The client_id field is missing.", "")
		return
	}

	e := r.FormValue("error")

	err := otpTemplate.Execute(w, map[string]interface{}{
		"query": "?client_id=" + url.QueryEscape(clientID),
		"error": e,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to render page", 500)
		return
	}
}
