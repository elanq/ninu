package ninu

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	tokenKey               = "oauth:access:token"
	defaultTokenExpiryTime = 24 * time.Hour
)

var (
	credential  *oauth2.Config
	mu          sync.Once
	cache       Cache
	oauthClient *http.Client

	errEmptySavedToken = errors.New("empty saved token")
)

func InitCredential() {
	mu.Do(func() {
		googleCred := os.Getenv("GOOGLE_CREDENTIAL")
		if googleCred == "" {
			panic(errors.New("Empty google credential"))
		}

		cred, err := google.ConfigFromJSON([]byte(googleCred), drive.DriveScope)
		if err != nil {
			panic(err)
		}
		credential = cred

		if token, _ := SavedToken(); token != nil {
			oauthClient = credential.Client(context.Background(), token)
			log.Println("Authorized via saved key")
		}
	})
}

func Credential() *oauth2.Config {
	return credential
}

func Client() *http.Client {
	return oauthClient
}

func SavedToken() (*oauth2.Token, error) {
	val, err := Redis.Get(tokenKey)
	if err != nil {
		return nil, err
	}

	if !json.Valid(val) {
		return nil, errEmptySavedToken
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(val, token); err != nil {
		return nil, err
	}

	return token, nil
}

func saveToken(token *oauth2.Token) error {
	payload, err := json.Marshal(token)
	if err != nil {
		return err
	}

	duration := defaultTokenExpiryTime
	if !token.Expiry.IsZero() {
		duration = time.Until(token.Expiry)
	}
	return Redis.Set(tokenKey, payload, duration)
}

func Authorize(authCode string) error {
	token, err := SavedToken()
	if err != nil && err != errEmptySavedToken {
		return err
	}

	if token != nil {
		oauthClient = Credential().Client(context.Background(), token)
		return nil
	}

	token, err = Credential().Exchange(context.Background(), authCode)
	if err != nil {
		return err
	}

	oauthClient = Credential().Client(context.Background(), token)
	if err := saveToken(token); err != nil {
		return err
	}
	return err
}

func AuthURL() string {
	return Credential().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
