package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/oauth2.v3"
)

var httpclient = http.Client{}

// APIClient is an oauth2 api client
type APIClient struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Trusted      bool     `json:"trusted"`
	Scopes       []string `json:"scopes"`
	RedirectURLs []string `json:"redirect_urls"`
	Secret       string   `json:"secret"`
}

// GetID of the client
func (c *APIClient) GetID() string {
	return c.ID
}

// GetSecret of the client
func (c *APIClient) GetSecret() string {
	return c.Secret
}

// GetDomain of the client
func (c *APIClient) GetDomain() string {
	data, _ := json.Marshal(c.RedirectURLs)

	if data == nil {
		return ""
	}

	return string(data)
}

// GetUserID of the client
func (c *APIClient) GetUserID() string {
	return ""
}

func getAPIClient(clientID string) (*APIClient, error) {
	request, err := http.NewRequest("GET", apiURL+"/apiclients/"+clientID, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", rootClientID+" "+rootClientSecret)
	request.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.StatusCode != 200 {
		return nil, errors.New("error or non 200 status")
	}

	var client = APIClient{}

	err = json.NewDecoder(resp.Body).Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// NewClientStore creates client store that handles the client information
func NewClientStore() *ClientStore {
	return &ClientStore{}
}

// ClientStore handles the client information
type ClientStore struct{}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	client, err := getAPIClient(id)
	if err != nil {
		return nil, fmt.Errorf("not found: %w", err)
	}

	return client, nil
}
