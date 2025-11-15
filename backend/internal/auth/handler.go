package auth

import (
	"backend-lastfm/internal/utility"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/glog"
	_ "github.com/joho/godotenv/autoload"
)



type AuthHandler struct {
    frontendUrl string
	sessionIdCookieName string

}

func NewAuthHandler(frontendUrl, sessionIdCookieName string) *AuthHandler {
	return &AuthHandler{frontendUrl, sessionIdCookieName}
}


// When the user hits /login by virtue of not being logged in already (eg. no token found on db)
// or the user is whimsical and explicitly goes to /login, this function will initiate the token
// acquiring flow for achieving the 3 legged Login Authentication flow with LastFM
func (h *AuthHandler) HandleLastFMLoginFlow(w http.ResponseWriter, r *http.Request) {

	if platform := chi.URLParam(r, "platform") ; platform != "lastfm"{
		glog.Errorf("Platform %s not implemented yet", platform)
		return
	}

	//Check if cookie exists
	glog.Info("Pass1")
	if resp, ok := CheckCookieValidity(r, h.sessionIdCookieName); ok {
		//Redirect to /Home
		url := strings.Join([]string{h.frontendUrl, "home"}, "/")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	} else {
		glog.Info(resp)
		//Code bit to start a new login flow.
		sessionID := *GenerateNewTx(r.RemoteAddr)

		glog.Infof("Recorded Login Attempt\n\tFrom IP: %s\n\tAssigned SessionID: %s\n\tCreated at: %s\n",
			sessionID.IP,
			sessionID.SessIDVerifier,
			sessionID.CreatedAt)
		
		//Set sid cookie to client


		


		http.SetCookie(w, &http.Cookie{
			Name:  h.sessionIdCookieName,
			Value: sessionID.SessIDVerifier,

			Expires:  time.Now().Add(time.Minute * 100),
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
		SetStateSid(randomState, sessionID.SessIDVerifier)

		//Sending login url with callback and state token
		url := GetInitLoginURL(os.Getenv("LASTFM_API_KEY"), randomState)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		// glog.Infof("Redirected URL: %s", url)
	}

}

func (h *AuthHandler) HandleLastFMCallbackFlow(w http.ResponseWriter, r *http.Request) {

	if platform := chi.URLParam(r, "platform"); platform != "lastfm" {
		glog.Errorf("Platform %s not implemented yet", platform)
		return
	}

	//Sample URL Path given is
	//http://127.0.0.1:3000/v1/auth/callback?sid={SID}&token={TOKEN}

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
	cookieSidReturned, err := r.Cookie(h.sessionIdCookieName)
	if err != nil {
		// glog.Info(cookieSidReturned.Value) //DEV

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie not found or invalid. Retry request.")
		return
	}

	//Perform SID and State verification check

	validationSid, ok := GetStateSid(stateReturned)
	if !ok {
		glog.Warning("State does not exist on state sid map")
	}

	if validationSid != cookieSidReturned.Value {
		glog.Warning("Invalid validation sid match try")
	}

	//If execution has reached this state, then we have verified that the callback is genuine

	//Fetch a web session
	webSessionURL, form := GetNewWebSessionURL(
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
	SetSidKey(validationSid, sessionKey)
	DelStateSid(stateReturned)

	//Perm redirect back to the original frontend.
	// http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)
	url := strings.Join([]string{h.frontendUrl, "home"}, "/")
	http.Redirect(w, r, url, http.StatusSeeOther)

	glog.Info("End of authentication flow")

}


func (h *AuthHandler) HandleAPIValidation(w http.ResponseWriter, r *http.Request) {

	if err, ok := CheckCookieValidity(r, h.sessionIdCookieName); ok {
		//Redirect to /Home
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Cookie Valid")
	} else {
		glog.Infof(err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie Invalid")
	}
}

func (h *AuthHandler) HandleLastFMLogOut(w http.ResponseWriter, r *http.Request) {
	if resp, ok := CheckCookieValidity(r, h.sessionIdCookieName); ok {
		newCookie := GetDeletedCookie(h.sessionIdCookieName) //Cookie value set to auto-expire yesterday
		http.SetCookie(w, newCookie)
		// delete(sessionIDTokenMap, resp)
		DelSidKey(resp)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Log out successful")

	} else {

		glog.Warningf("%s unsuccessfully tried to request log-out. UA: %s", r.RemoteAddr, r.UserAgent())
		glog.Warning(resp)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Log out unsuccessful")
	}
}




