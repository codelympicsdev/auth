package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis"
	oredis "gopkg.in/go-oauth2/redis.v3"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var oauth2server *server.Server

func init() {
	gob.Register(url.Values{})

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.SetValidateURIHandler(validateURI)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(NewAccessGenerator())
	manager.MapClientStorage(NewClientStore())
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr:     redisURI,
		Password: redisPassword,
	}, "clyauthtok"))
	config := server.NewConfig()
	oauth2server = server.NewServer(config, manager)
	oauth2server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.ClientCredentials, oauth2.Refreshing, oauth2.Implicit)
	oauth2server.SetUserAuthorizationHandler(userAuthorizationHandler)

	oauth2server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	oauth2server.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
}

func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	store, err := rs.Get(r, "codelympics-auth")
	if err != nil {
		return "", err
	}

	uid, ok := store.Values["signed_in_user_id"]
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Values["return_form"] = r.Form
		err = store.Save(r, w)
		if err != nil {
			return "", err
		}

		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}

	userID := uid.(string)
	delete(store.Values, "signed_in_user_id")
	err = store.Save(r, w)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	store, err := rs.Get(r, "codelympics-auth")
	if err != nil {
		internalServerError(w)
		return
	}

	var form url.Values
	if v, ok := store.Values["return_form"]; ok {
		form = v.(url.Values)
	}
	r.Form = form

	delete(store.Values, "return_form")
	err = store.Save(r, w)
	if err != nil {
		internalServerError(w)
		return
	}

	err = oauth2server.HandleAuthorizeRequest(w, r)
	if err != nil {
		renderError(w, "Bad Request", "Please check the parameters supplied to the authorization url.", "")
		return
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	err := oauth2server.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateURI(urls string, redirectURI string) error {
	allowedURLs := []string{}
	err := json.Unmarshal([]byte(urls), &allowedURLs)
	if err != nil {
		return err
	}

	for _, allowedURL := range allowedURLs {
		allowed, err := url.Parse(allowedURL)
		if err != nil {
			return err
		}

		redirect, err := url.Parse(redirectURI)
		if err != nil {
			return err
		}
		if strings.HasSuffix(redirect.Hostname(), allowed.Hostname()) && strings.HasSuffix(redirect.Path, allowed.Path) && strings.HasSuffix(redirect.Scheme, allowed.Scheme) {
			return nil
		}
	}

	return errors.ErrInvalidRedirectURI
}
