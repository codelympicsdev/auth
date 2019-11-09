package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// SigninResponse the response to a signin request
type SigninResponse struct {
	UserID      string `json:"user_id"`
	Requires2FA bool   `json:"requires_2fa"`
}

// SigninEmailPasswordRequest is a request that gets a user and checks their credentials by email and password
type SigninEmailPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CheckResponse is the response to a check request
type CheckResponse struct {
	Valid bool `json:"valid"`
}

// CheckOTPRequest is a request to check otp credentials for a user
type CheckOTPRequest struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
}

func signinEmailPassword(email string, password string) (string, bool, error) {
	sendData, err := json.Marshal(SigninEmailPasswordRequest{
		Email:    email,
		Password: password,
	})

	request, err := http.NewRequest("POST", apiURL+"/auth/signin/emailpassword", bytes.NewBuffer(sendData))
	if err != nil {
		return "", false, err
	}

	request.Header.Set("Authorization", rootClientID+" "+rootClientSecret)
	request.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.Do(request)
	if err != nil {
		return "", false, err
	}
	if resp == nil || resp.StatusCode != 200 {
		return "", false, errors.New("error or non 200 status")
	}

	var data = SigninResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", false, err
	}

	return data.UserID, data.Requires2FA, nil
}

func checkOTP(userID string, otp string) (bool, error) {
	sendData, err := json.Marshal(CheckOTPRequest{
		UserID: userID,
		OTP:    otp,
	})

	request, err := http.NewRequest("POST", apiURL+"/auth/check/otp", bytes.NewBuffer(sendData))
	if err != nil {
		return false, err
	}

	request.Header.Set("Authorization", rootClientID+" "+rootClientSecret)
	request.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.Do(request)
	if err != nil {
		return false, err
	}
	if resp == nil || resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Printf("%v", string(b))
		return false, errors.New("error or non 200 status")
	}

	var data = CheckResponse{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	return data.Valid, nil
}
