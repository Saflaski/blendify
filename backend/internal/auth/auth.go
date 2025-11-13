package auth

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type ClientKey string
var SIDCOOKIE = "sid"
var LASTFM_BASE_AUTH_API = "http://www.last.fm/api/auth/"
var LASTFM_ROOT_API = "http://ws.audioscrobbler.com/2.0/"

var HOME_URL = "http://localhost:3000"
var LASTFM_CALLBACK = HOME_URL + "/v1/auth/callback/lastfm"

func GetClientCookie(val string) *http.Cookie {

	secure := false //TODO SET TRUE FOR PROD
	cookieReturn := http.Cookie{
		Name:  "sid",
		Value: val,
		Path:  "/",
		// MaxAge:   0,
		Expires:  time.Now().Add(24 * time.Hour), //Cookie expires in 24 hours
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return &cookieReturn
}

func GetDeletedCookie() *http.Cookie {

	secure := false //TODO Set true for PROD
	cookieReturn := http.Cookie{
		Name:  "sid",
		Value: "",
		Path:  "/",
		// MaxAge:   0,
		Expires:  time.Unix(0, 0), //Cookie expires in 24 hours
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return &cookieReturn
}

// Returns
// string: Cookie value if success, error message if not
// bool: success value of operation
func CheckCookieValidity(r *http.Request) (string, bool) {

	cookie, err := r.Cookie(SIDCOOKIE)

	//Check if we even get a cookie first
	if err != nil {
		if err == http.ErrNoCookie {
			return "No cookie recieved, starting fresh login flow", false
		} else {
			return "Error during cookie retrieval", false
		}
	}

	if cookie.Value == "" {
		return "Cookie value \"\"", false
	}

	// _, ok := sessionIDTokenMap[cookie.Value]
	_, ok := GetSidKey(cookie.Value)
	if !ok {
		return string("SID not found in map. Given value: " + cookie.Value), false
	}

	return cookie.Value, true

}

func GetInitLoginURL(api_key string, state string) string {
	q := url.Values{}
	q.Set("api_key", api_key)
	q.Set("cb", string(LASTFM_CALLBACK+"?state="+state))

	requestURL := LASTFM_BASE_AUTH_API + "?" + q.Encode()
	// glog.Info(requestURL)
	return requestURL
}

func GetSessionAPISignature(api_key string, token string) string {

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

func GetNewWebSessionURL(api_key string, token string) (string, url.Values) {

	q := url.Values{}
	q.Set("method", "auth.getSession")
	q.Set("api_key", api_key)
	q.Set("token", token)
	q.Set("api_sig", GetSessionAPISignature(api_key, token))

	// requestURL := LASTFM_ROOT_API + "?" + q.Encode()
	return LASTFM_ROOT_API, q //Didn't realise POST would require them as different //TODO Make this cleaner
}

type Tx struct {
	SessIDVerifier string
	CreatedAt      time.Time
	IP             string
}

func GenerateNewTx(userIP string) *Tx {
	sessIDVerifier := uuid.New().String()
	tx := Tx{
		SessIDVerifier: sessIDVerifier,
		CreatedAt:      time.Now(),
		IP:             userIP,
	}
	// glog.Info(sessIDVerifier)
	// glog.Info(userIP)

	return &tx //Return a pointer to tx
}
