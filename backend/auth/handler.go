package auth

import (
	"backend-lastfm/utility"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/golang/glog"
	_ "github.com/joho/godotenv/autoload"
)

var sessionIDTokenMap map[string]string = make(map[string]string)
var SIDCOOKIE = "sid"
var FRONTEND_ROOT_URL = "http://127.0.0.1:5173/"

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, FRONTEND_ROOT_URL+"/login", http.StatusTemporaryRedirect)
}

// When the user hits /login by virtue of not being logged in already (eg. no token found on db)
// or the user is whimsical and explicitly goes to /login, this function will initiate the token
// acquiring flow for achieving the 3 legged Login Authentication flow with LastFM
func handleLoginFlow(w http.ResponseWriter, r *http.Request) {

	//Check if cookie exists

	if resp, ok := checkCookieValidity(r); ok {
		//Redirect to /Home
		http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)

	} else {
		glog.Warning(resp)
		//Code bit to start a new login flow.
		sessionID := *generateNewTx(r.RemoteAddr)
		glog.Infof("Recorded Login \n\tFrom IP: %s\n\tAssigned SessionID: %s\n\tCreated at: %s\n",
			sessionID.IP,
			sessionID.SessIDVerifier,
			sessionID.CreatedAt)
		url := getInitLoginURL(os.Getenv("LASTFM_API_KEY"), sessionID.SessIDVerifier)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		// glog.Infof("Redirected URL: %s", url)
	}

}

// func handleCookieValidation(w http.ResponseWriter, r *http.Request) {
// 	if resp, ok := checkCookieValidity(r); ok {
// 		//Redirect to /Home

// 	}
// }

func checkCookieValidity(r *http.Request) (string, bool) {

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

	_, ok := sessionIDTokenMap[cookie.Value]
	if !ok {
		return string("SID not found in map. Given value: " + cookie.Value), false
	}

	return cookie.Value, true

}

func handleCallbackFlow(w http.ResponseWriter, r *http.Request) {

	tokenReturned := r.URL.Query().Get("token")
	callbackReturned := r.URL.Query().Get("sid")
	glog.Info("Callback returned:")
	path := strings.TrimPrefix(callbackReturned, LASTFM_CALLBACK)
	path = strings.TrimSuffix(path, "/")

	if path == "" {
		glog.Infof("No Session ID provided, ignoring")
		return
	}

	sessionID, err := url.QueryUnescape(path)
	if err != nil {
		glog.Fatal("Could not decode callback URL: ", path)
	}
	glog.Infof(sessionID)

	glog.Infof("SID-Token association[%s : %s]", sessionID, tokenReturned)

	//Fetch a web session
	webSessionURL, form := getNewWebSessionURL(
		os.Getenv("LASTFM_API_KEY"),
		tokenReturned,
	)

	glog.Info("Web Session Request URL: ", webSessionURL)
	resp, err := http.Post(
		webSessionURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		glog.Errorf("Request failed: %v", err)
		return
	}
	defer resp.Body.Close() // always close body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read body: %v", err)
		return
	}

	glog.Infof("Response body: %s", string(body))

	xmlStruct := utility.ParseXMLSessionKey(body)
	sessionKey := xmlStruct.Session.Key

	//Assigning the mapping for recording users for later re-auth between frontend and backend
	sessionIDTokenMap[sessionID] = sessionKey

	//Set cookie
	cookieToBeSet := GetClientCookie(sessionID)
	http.SetCookie(w, cookieToBeSet)

	//Perm redirect back to the original frontend.
	http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)

	glog.Info("End of authentication flow")

}

// CORS MIDDLWARE
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == "http://127.0.0.1:5173" || true {
			w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleAPIValidation(w http.ResponseWriter, r *http.Request) {

	if resp, ok := checkCookieValidity(r); ok {
		//Redirect to /Home
		// http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Cookie Valid")
		glog.Info(resp)
	} else {
		glog.Info(resp)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie Invalid")
	}
}

func ServerStart() {
	defer glog.Flush()

	glog.Info("Backend started with ClientID", os.Getenv("LASTFM_ID"))

	// http.HandleFunc("/", handleRoot)
	mux := http.NewServeMux()
	handler := cors(mux)
	mux.HandleFunc("/api/validate/", handleAPIValidation)
	mux.HandleFunc("/oauth/lastfm/login", handleLoginFlow)
	mux.HandleFunc("/oauth/lastfm/callback", handleCallbackFlow)

	http.ListenAndServe(":3000", handler) //127.0.0.1:3000

}
