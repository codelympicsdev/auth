package main

import (
	"log"

	"github.com/pquerna/otp/totp"
	"golang.org/x/oauth2"
)

const (
	authServerURL = "http://localhost:8081"
)

var (
	config = oauth2.Config{
		ClientID:     "5dbf4c151c9d440000ffa241",
		ClientSecret: "bad-secret",
		Scopes:       []string{"user.basic"},
		RedirectURL:  "http://localhost:3000/auth-callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth2/auth",
			TokenURL: authServerURL + "/oauth2/token",
		},
	}
)

func main() {
	log.Println(totp.Generate(totp.GenerateOpts{
		Issuer:      "codelympics.dev",
		AccountName: "hello@lcas.dev",
	}))
	/*
		t, err := config.Exchange(context.Background(), "GEE_Q_VGNT6STALZKHJFPW")
		log.Println(err)
		log.Println(t)
	*/

	//log.Println(config.Exchange(context.Background(), "PBOEHDMBXFA8Z4OV-LBHEQ"))

}
