package auth

import (
	"backend-lastfm/utility"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/golang/glog"
	_ "github.com/joho/godotenv/autoload"
)

var stateTokenMap map[string]string = make(map[string]string)

// When the user hits /login by virtue of not being logged in already (eg. no token found on db)
// or the user is whimsical and explicitly goes to /login, this function will initiate the token
// acquiring flow for achieving the 3 legged Login Authentication flow with LastFM
func handleLoginFlow(w http.ResponseWriter, r *http.Request) {

	state := *generateNewTx(net.IP(r.RemoteAddr))

	glog.Infof("Recorded Login \n\tFrom IP: %s\n\tAssigned State: %s\n\tCreated at: %s\n",
		state.IP, state.StateVerifier, state.CreatedAt)
	url := getInitLoginURL(os.Getenv("LASTFM_API_KEY"), state.StateVerifier)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	glog.Infof("Redirected URL: %s", url)

}

func handleCallbackFlow(w http.ResponseWriter, r *http.Request) {

	tokenReturned := r.URL.Query().Get("token")
	callbackReturned := r.URL.Query().Get("state")
	glog.Info("Callback returned:")
	path := strings.TrimPrefix(callbackReturned, LASTFM_CALLBACK)
	path = strings.TrimSuffix(path, "/")

	if path == "" {
		glog.Infof("No stateVerifier provided, ignoring")
		return
	}

	stateVerifier, err := url.QueryUnescape(path)
	if err != nil {
		glog.Fatal("Could not decode callback URL: ", path)
	}
	glog.Infof(stateVerifier)

	glog.Infof("State-Token association[%s : %s]", stateVerifier, tokenReturned)

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
	defer resp.Body.Close() // always close body!

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read body: %v", err)
		return
	}

	glog.Infof("Response body: %s", string(body))

	xmlStruct := utility.ParseXMLSessionKey(body)
	sessionKey := xmlStruct.Session.Key

	//Assigning the mapping for recording users for later re-auth between frontend and backend
	stateTokenMap[stateVerifier] = sessionKey

	glog.Info("End of authentication flow")

}

func ServerStart() {
	defer glog.Flush()

	glog.Info("Backend started with ClientID", os.Getenv("LASTFM_ID"))

	http.HandleFunc("/oauth/lastfm/login", handleLoginFlow)
	http.HandleFunc("/oauth/lastfm/callback", handleCallbackFlow)

	http.ListenAndServe(":3000", nil) //127.0.0.1:3000

}
