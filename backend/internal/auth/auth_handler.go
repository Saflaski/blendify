package auth

import (
	"backend-lastfm/internal/utility"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/glog"
	_ "github.com/joho/godotenv/autoload"
)

type AuthHandler struct {
	frontendUrl         string
	sessionIdCookieName string
	svc                 AuthService
}

func NewAuthHandler(frontendUrl, sessionIdCookieName string, svc AuthService) *AuthHandler {
	return &AuthHandler{frontendUrl, sessionIdCookieName, svc}
}

func (h *AuthHandler) HandleLastFMLoginFlow(w http.ResponseWriter, r *http.Request) {
	if platform := chi.URLParam(r, "platform"); platform != "lastfm" {
		glog.Errorf("Platform %s not implemented yet", platform)
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Platform %s not implemented yet", platform)
	}

	url := strings.Join([]string{h.frontendUrl, "home"}, "/")
	//Check if cookie exists
	cookie, err := r.Cookie(h.sessionIdCookieName)

	if err != nil { //Either no cookie found or error retrieving cookie
		if err == http.ErrNoCookie {
			//Start login flow
			err := h.startNewLoginFlow(w, r)
			if err != nil {
				glog.Errorf("Error starting new login flow: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error starting new login flow")
				return
			}
		} else {
			glog.Errorf("Error retrieving cookie: %v", err)
			w.WriteHeader(http.StatusBadRequest) //Bad Request
			fmt.Fprintf(w, "Error retrieving cookie.")
			return
		}
	} else { //There is a cookie with SID

		found, err := h.svc.CheckCookieValidity(r.Context(), cookie.Value)
		if err != nil { //Error during validity check
			glog.Errorf("Error checking cookie validity: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error checking cookie validity")
			return
		}
		if found {
			//Cookie is found and is valid. Return to home
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
	}
}

func (h *AuthHandler) startNewLoginFlow(w http.ResponseWriter, r *http.Request) error {

	sessionID, state, err := h.svc.GenerateNewStateAndSID(r.Context()) //Core logic so...Service?
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating security tokens")
	}

	http.SetCookie(w, &http.Cookie{
		Name:  h.sessionIdCookieName,
		Value: sessionID,

		Expires:  time.Now().Add(time.Minute * 100),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	loginURL := h.svc.GetInitLoginURL(state)
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)

	return nil
}

// // When the user hits /login by virtue of not being logged in already (eg. no token found on db)
// // or the user is whimsical and explicitly goes to /login, this function will initiate the token
// // acquiring flow for achieving the 3 legged Login Authentication flow with LastFM
// func (h *AuthHandler) HandleLastFMLoginFlow(w http.ResponseWriter, r *http.Request) {

// 	if platform := chi.URLParam(r, "platform") ; platform != "lastfm"{
// 		glog.Errorf("Platform %s not implemented yet", platform)
// 		return
// 	}

// 	//TODO Delete this after full HSR implementation
// 	//h.svc.IsSessionValid(X)
// 	//if yes then do that
// 	//if no then redirect

// 	//Check if cookie exists

// 	//if _, ok := h.svc.IsSessionValid(Cookie, h.sessionIdCookieName); ok {
// 	if _, ok := h.svc.CheckCookieValidity(r, h.sessionIdCookieName); ok {
// 		//Since there is a valid cookie, user is redirected to /home/
// 		url := strings.Join([]string{h.frontendUrl, "home"}, "/")
// 		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

// 	} else {

// 		//Code bit to start a new login flow.
// 		sessionID := *h.svc.GenerateNewTx(r.RemoteAddr)	//Core logic so...Service?

// 		glog.Infof("Recorded Login Attempt\n\tFrom IP: %s\n\tAssigned SessionID: %s\n\tCreated at: %s\n",
// 			sessionID.IP,
// 			sessionID.SessIDVerifier,					//Still need this?
// 			sessionID.CreatedAt)

// 		http.SetCookie(w, &http.Cookie{					//Keep in H
// 			Name:  h.sessionIdCookieName,
// 			Value: sessionID.SessIDVerifier,

// 			Expires:  time.Now().Add(time.Minute * 100),
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   false,
// 			SameSite: http.SameSiteLaxMode,
// 		})

// 		randomStateByte := make([]byte, 16)
// 		_, err := rand.Read(randomStateByte)
// 		randomState := hex.EncodeToString(randomStateByte)	//Move to S
// 		if err != nil {
// 			glog.Warning("Cannot create state token")
// 			panic("")
// 		}

// 		//Saving state : SID map internally
// 		SetStateSid(randomState, sessionID.SessIDVerifier)	//Move to S

// 		//Sending login url with callback and state token
// 		url := h.svc.GetInitLoginURL(os.Getenv("LASTFM_API_KEY"), randomState)
// 		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

// 	}

// }

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

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")

	stateReturned := r.URL.Query().Get("state")
	tokenReturned := r.URL.Query().Get("token")

	//Retrieve SID
	cookieSidReturned, err := r.Cookie(h.sessionIdCookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Cookie not found or invalid. Retry request.")
		return
	}

	//Perform SID and State verification check
	validationSid, err := h.svc.ConsumeStateAndReturnSID(r.Context(), stateReturned)
	if err != nil {
		glog.Errorf("Error starting new login flow, failed callback: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with callback verification from LastFM. Retry login flow.")
	}

	if validationSid != cookieSidReturned.Value {
		glog.Warning("Invalid validation sid match try")
		w.WriteHeader(http.StatusRequestTimeout)
		fmt.Fprintf(w, "Error involving SessionID verification. Do you have cookies enabled? Retry login flow.")
	}

	//If execution has reached this state, then we have verified that the callback is genuine

	//Fetch a web session
	webSessionURL, form := h.svc.GetNewWebSessionURL(tokenReturned)

	resp, err := http.Post(
		webSessionURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		glog.Errorf("Request failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error requesting web session from LastFM. Clear cookies and retry login flow. Is LastFM up?")
		return
	}

	defer resp.Body.Close() // always close body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Request failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading response from LastFM. Clear cookies and retry login flow. Is LastFM up?")
		return
	}

	//Parsing XML response

	xmlStruct := utility.ParseXMLSessionKey(body)
	sessionKey := xmlStruct.Session.Key

	//Assigning the mapping for recording users for later re-auth between frontend and backend

	h.svc.SetSessionKey(r.Context(), validationSid, sessionKey)

	//Perm redirect back to the original frontend.
	// http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)
	url := strings.Join([]string{h.frontendUrl, "home"}, "/")
	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (h *AuthHandler) HandleAPIValidation(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(h.sessionIdCookieName)
	if err == nil {
		//Validating found cookie
		found, err := h.svc.CheckCookieValidity(r.Context(), cookie.Value)
		if err != nil {
			glog.Errorf("Error checking cookie validity: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error checking cookie validity")
			return
		} else {
			if found {
				//Valid Cookie found - Redirect to /home
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "Cookie Valid")
				return
			} else { //Cookie expired or DB mismatch
				glog.Infof("Cookie invalid or expired for request from %s. UA: %s", r.RemoteAddr, r.UserAgent())
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Cookie Expired or Invalid")
			}
		}
	} else {
		//No cookie found or invalid cookie
		glog.Infof("No cookie found for request from %s. UA: %s", r.RemoteAddr, r.UserAgent())
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "No Cookie Found")
		return
	}

}

func (h *AuthHandler) HandleLastFMLogOut(w http.ResponseWriter, r *http.Request) {

	newCookie := h.svc.GetDeletedCookie(h.sessionIdCookieName) //Cookie that expires immediately ie a deleted cookie

	cookie, err := r.Cookie(h.sessionIdCookieName) //Get browser cookie

	if err != nil { //Either no cookie found or error retrieving cookie
		//Nothing to do here. Just delete the cookie
	} else { //There is a cookie with SID

		found, err := h.svc.CheckCookieValidity(r.Context(), cookie.Value)
		if err != nil { //Error during validity check
			glog.Errorf("Error checking cookie validity: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error checking cookie validity")
			return
		}
		if found {
			//Cookie is found and is valid. Return to home
			http.SetCookie(w, newCookie) //Set deleted cookie
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Logout successful.")

			//Proceed to delete from repo
			err := h.svc.DelSidKey(r.Context(), cookie.Value)
			if err != nil {
				glog.Errorf("Error deleting SID key from repository: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error during logout process. Retry logout.")
				return
			}
		} else {
			//Cookie invalid or expired. Just delete cookie
			http.SetCookie(w, newCookie) //Set deleted cookie
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Logout successful.")

			//Cant delete from repo since no valid SID found.
		}
	}

}
