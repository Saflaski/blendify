package auth

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type ClientKey string

var LASTFM_BASE_AUTH_API = "http://www.last.fm/api/auth/"
var LASTFM_ROOT_API = "http://ws.audioscrobbler.com/2.0/"
var HOME_URL = "http://127.0.0.1:3000/"
var LASTFM_CALLBACK = HOME_URL + "oauth/lastfm/callback"

func getOAuthLoginURL(userid string, tx *Tx) string {

	q := url.Values{}
	q.Set("api_key", os.Getenv("LASTFM_API_KEY"))
	q.Set("sid", tx.SessIDVerifier)
	return LASTFM_BASE_AUTH_API
}

func getInitLoginURL(api_key string, sessIDVerifier string) string {
	q := url.Values{}
	q.Set("api_key", api_key)
	q.Set("cb", string(LASTFM_CALLBACK+"?sid="+sessIDVerifier))

	requestURL := LASTFM_BASE_AUTH_API + "?" + q.Encode()
	glog.Info(requestURL)
	return requestURL
}

func getSessionAPISignature(api_key string, token string) string {

	secret := os.Getenv("LASTFM_SECRET")
	raw_string := string("api_key" + api_key + "method" + "auth.getSession" + "token" + token + secret)

	api_sig := md5.Sum([]byte(raw_string))

	signature := hex.EncodeToString(api_sig[:])

	if len(signature) == 32 {
		return signature
	} else {
		glog.Fatal("MD5 fail")
		panic("")
	}

}

func getNewWebSessionURL(api_key string, token string) (string, url.Values) {

	q := url.Values{}
	q.Set("method", "auth.getSession")
	q.Set("api_key", api_key)
	q.Set("token", token)
	q.Set("api_sig", getSessionAPISignature(api_key, token))

	// requestURL := LASTFM_ROOT_API + "?" + q.Encode()
	return LASTFM_ROOT_API, q //Didn't realise POST would require them as different //TODO Make this cleaner
}

type Tx struct {
	SessIDVerifier string
	CreatedAt      time.Time
	IP             string
}

type memStore struct {
	mu sync.Mutex    //Mutex for keeping locks
	m  map[string]Tx //Map for mapping UUID to Tx?
}

func generateNewTx(userIP string) *Tx {
	sessIDVerifier := uuid.New().String()
	tx := Tx{
		SessIDVerifier: sessIDVerifier,
		CreatedAt:      time.Now(),
		IP:             userIP,
	}

	return &tx //Return a pointer to tx
}
