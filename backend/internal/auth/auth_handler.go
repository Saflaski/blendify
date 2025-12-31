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
	svc     AuthService
	config  Config
	UserKey contextKey
}

type contextKey string

const UserKey contextKey = "userid"

func NewAuthHandler(svc AuthService, cfg Config) *AuthHandler {
	return &AuthHandler{svc, cfg, UserKey}
}

func (h *AuthHandler) HandleLastFMLoginFlow(w http.ResponseWriter, r *http.Request) {
	if platform := chi.URLParam(r, "platform"); platform != "lastfm" {
		glog.Errorf("Platform %s not implemented yet", platform)
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Platform %s not implemented yet", platform)
	}

	// Shouldn't this be such that we should assume no cookie?	//TODO
	// And even if there is a cookie, we delete it? And restart login?

	url := strings.Join([]string{h.config.FrontendURL, "home"}, "/")
	//Check if cookie exists
	cookie, err := r.Cookie(h.config.FrontendCookieName)

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

		found, err := h.svc.IsSIDValid(r.Context(), cookie.Value)
		if err != nil { //Error during validity check
			glog.Errorf("Error checking cookie validity: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error checking cookie validity")
			return
		}
		if found {
			glog.Info(" FOUND VALID COOKIE SID")
			//Cookie is found and is valid. Return to home
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		} else {
			glog.Info(" DID NOT FIND VALID COOKIE SID")

			err := h.startNewLoginFlow(w, r)
			if err != nil { //Error during validity check
				glog.Errorf("Error redirecting to new login flow: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error redirect to new login flow. Try deleting all cookies")
				return
			}
		}
	}
}

func (h *AuthHandler) startNewLoginFlow(w http.ResponseWriter, r *http.Request) error {

	sessionID, state, err := h.svc.GenerateNewStateAndSID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error generating security tokens")
	}

	http.SetCookie(w, &http.Cookie{
		Name:  h.config.FrontendCookieName,
		Value: sessionID,

		Expires:  time.Now().Add(h.config.ExpiryDuration),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //TODO Change to true for Prod
		SameSite: http.SameSiteLaxMode,
	})

	loginURL := h.svc.GetInitLoginURL(state)
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)

	return nil
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

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")

	stateReturned := r.URL.Query().Get("state")
	tokenReturned := r.URL.Query().Get("token")

	//Retrieve SID
	cookieSidReturned, err := r.Cookie(h.config.FrontendCookieName)
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
	// sessionKey := xmlStruct.Session.Key
	userName := xmlStruct.Session.Name
	//Assigning the mapping for recording users for later re-auth between frontend and backend

	// h.svc.SetSessionKey(r.Context(), validationSid, sessionKey)
	_, err = h.svc.MakeNewUser(r.Context(), validationSid, userName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not register/login user")
		glog.Errorf("Error during register/login: %w", err)
		return
	}

	//Perm redirect back to the original frontend.
	// http.Redirect(w, r, "http://127.0.0.1:5173/home", http.StatusTemporaryRedirect)
	url := strings.Join([]string{h.config.FrontendURL, "home"}, "/") //This should be something that frontend handles, not backend.
	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (h *AuthHandler) HandleAPIValidation(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(h.config.FrontendCookieName)
	if err == nil {
		//Validating found cookie
		found, err := h.svc.IsSIDValid(r.Context(), cookie.Value)
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

	newCookie := h.svc.GetDeletedCookie(h.config.FrontendCookieName) //Cookie that expires immediately ie a deleted cookie

	cookie, err := r.Cookie(h.config.FrontendCookieName) //Get browser cookie

	if err != nil { //Either no cookie found or error retrieving cookie
		//Nothing to do here. Just delete the cookie
	} else { //There is a cookie with SID

		found, err := h.svc.IsSIDValid(r.Context(), cookie.Value)
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
