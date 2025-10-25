package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

var (
	oAuthConf = oauth2.Config{
		ClientID:     os.Getenv("GIT_CLIENT_ID"),
		ClientSecret: os.Getenv("GIT_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:3000/oauth2/callback",
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:       "https://github.com/login/oauth/authorize",
			TokenURL:      "https://github.com/login/oauth/access_token",
			DeviceAuthURL: "https://github.com/login/device/code",
		},
	}
	randomState, err = _generateState()
)

func _generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(b) // no padding
	return s, nil
}

func _handleLoginFlow(w http.ResponseWriter, r *http.Request) {
	glog.Info("Login Flow started")

	url := oAuthConf.AuthCodeURL(string(randomState))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func _handleCallback(w http.ResponseWriter, r *http.Request) {
	glog.Info("Callback")

	if r.FormValue("state") != randomState {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := oAuthConf.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	fmt.Println(token)
}

func _main() {
	defer glog.Flush()
	fmt.Println(oAuthConf.ClientID)
	http.HandleFunc("/oauth2/login", handleLoginFlow)
	http.HandleFunc("/oauth2/callback", _handleCallback)
	http.ListenAndServe(":3000", nil)

}
