package main

import (
	"backend-lastfm/internal/auth"
	"backend-lastfm/internal/utility"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	_ "github.com/joho/godotenv/autoload"
)


var FRONTEND_ROOT_URL = "http://127.0.0.1:5173"

// When the user hits /login by virtue of not being logged in already (eg. no token found on db)
// or the user is whimsical and explicitly goes to /login, this function will initiate the token
// acquiring flow for achieving the 3 legged Login Authentication flow with LastFM
func handleLoginFlow(w http.ResponseWriter, r *http.Request) {

	//Check if cookie exists
	glog.Info("Pass1")
	if resp, ok := auth.CheckCookieValidity(r); ok {
		//Redirect to /Home
		http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)

	} else {
		glog.Info(resp)
		//Code bit to start a new login flow.
		sessionID := *auth.GenerateNewTx(r.RemoteAddr)

		glog.Infof("Recorded Login Attempt\n\tFrom IP: %s\n\tAssigned SessionID: %s\n\tCreated at: %s\n",
			sessionID.IP,
			sessionID.SessIDVerifier,
			sessionID.CreatedAt)

		//Set sid cookie to client
		http.SetCookie(w, &http.Cookie{
			Name:  "sid",
			Value: sessionID.SessIDVerifier,

			Expires:  time.Now().Add(time.Second * 100),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		randomStateByte := make([]byte, 16)
		_, err := rand.Read(randomStateByte)
		randomState := hex.EncodeToString(randomStateByte)
		if err != nil {
			glog.Warning("Cannot create state token")
			panic("")
		}

		//Saving state : SID map internally
		auth.SetStateSid(randomState, sessionID.SessIDVerifier)

		//Sending login url with callback and state token
		url := auth.GetInitLoginURL(os.Getenv("LASTFM_API_KEY"), randomState)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		// glog.Infof("Redirected URL: %s", url)
	}

}

func handleCallbackFlow(w http.ResponseWriter, r *http.Request) {

	//Sample URL Path given is
	//http://127.0.0.1:3000/oauth/lastfm/callback?sid={SID}&token={TOKEN}

	//Validate
	// _, timeout := context.WithTimeout(r.Context(), 10*time.Second)
	// defer timeout()

	//Security fix to set no cache and no referrer
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")

	//Retrieve State, Token
	stateReturned := r.URL.Query().Get("state")
	tokenReturned := r.URL.Query().Get("token")

	//Retrieve SID
	cookieSidReturned, err := r.Cookie("sid")
	if err != nil {
		glog.Info(cookieSidReturned.Value) //DEV

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie not found or invalid. Retry request.")

	}

	//Perform SID and State verification check

	validationSid, ok := auth.GetStateSid(stateReturned)
	if !ok {
		glog.Warning("State does not exist on state sid map")
	}

	if validationSid != cookieSidReturned.Value {
		glog.Warning("Invalid validation sid match try")
	}

	//If execution has reached this state, then we have verified that the callback is genuine

	//Fetch a web session
	webSessionURL, form := auth.GetNewWebSessionURL(
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
	auth.SetSidKey(validationSid, sessionKey)
	auth.DelStateSid(stateReturned)

	//Perm redirect back to the original frontend.
	// http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)
	http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusSeeOther)

	glog.Info("End of authentication flow")

}

// // CORS MIDDLWARE
// func cors(next http.Handler) http.Handler {

// 	allowed := map[string]bool{
// 		"http://127.0.0.1:5173": true,
// 		"http://localhost:5173": true,
// 	}
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
// 		if origin != "" && allowed[origin] {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Set("Vary", "Origin")
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
// 			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		}

// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func handleAPIValidation(w http.ResponseWriter, r *http.Request) {

	if _, ok := auth.CheckCookieValidity(r); ok {
		//Redirect to /Home
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Cookie Valid")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie Invalid")

	}
}

func handleLogOut(w http.ResponseWriter, r *http.Request) {
	if resp, ok := auth.CheckCookieValidity(r); ok {
		newCookie := auth.GetDeletedCookie() //Cookie value set to auto-expire yesterday
		http.SetCookie(w, newCookie)
		// delete(sessionIDTokenMap, resp)
		auth.DelSidKey(resp)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Log out successful")

	} else {

		glog.Warningf("%s unsuccessfully tried to request log-out. UA: %s", r.RemoteAddr, r.UserAgent())
		glog.Warning(resp)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Log out unsuccessful")
	}
}

// func ServerStart() { //--------------- Scheduled to be deleted
// 	defer glog.Flush()

// 	glog.Info("Backend started with ClientID", os.Getenv("LASTFM_ID"))
	



// 	// http.HandleFunc("/", handleRoot)
// 	mux := http.NewServeMux()
// 	handler := cors(mux)

// 	mux.HandleFunc("/api/logout/", handleLogOut)
// 	mux.HandleFunc("/api/validate/", handleAPIValidation)
// 	mux.HandleFunc("/oauth/lastfm/login", handleLoginFlow)
// 	mux.HandleFunc("/oauth/lastfm/callback", handleCallbackFlow)

// 	http.ListenAndServe(":3000", handler) //127.0.0.1:3000

// }

