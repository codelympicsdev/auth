package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/utils/uuid"
)

// GenerateTokenResponse is the response for generating a token
type GenerateTokenResponse struct {
	Token string `json:"token"`
}

// GenerateTokenRequest is a request to generate a token
type GenerateTokenRequest struct {
	ClientID string   `json:"client_id"`
	UserID   string   `json:"user_id"`
	Scopes   []string `json:"scopes"`
}

func generateToken(userID string, clientID string, scopes []string) (string, error) {
	sendData, err := json.Marshal(GenerateTokenRequest{
		ClientID: clientID,
		UserID:   userID,
		Scopes:   scopes,
	})

	request, err := http.NewRequest("POST", apiURL+"/auth/generatetoken", bytes.NewBuffer(sendData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", rootClientID+" "+rootClientSecret)
	request.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.Do(request)
	if err != nil {
		return "", err
	}
	if resp == nil || resp.StatusCode != 200 {
		return "", errors.New("error or non 200 status")
	}

	var data = GenerateTokenResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Token, nil
}

// NewAccessGenerator creates a new access token generator
func NewAccessGenerator() *AccessGenerator {
	return &AccessGenerator{}
}

// AccessGenerator generates access tokens using the codelympics API
type AccessGenerator struct{}

// Token based on the UUID generated token
func (a *AccessGenerator) Token(data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	access, err := generateToken(data.UserID, data.Client.GetID(), strings.Split(data.TokenInfo.GetScope(), ","))
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).Bytes())
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}
