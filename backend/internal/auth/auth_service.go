package auth

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type ClientKey string

var LASTFM_BASE_AUTH_API = "http://www.last.fm/api/auth/"
var LASTFM_ROOT_API = "http://ws.audioscrobbler.com/2.0/"

var HOME_URL = "http://localhost:3000"
var LASTFM_CALLBACK = HOME_URL + "/v1/auth/callback/lastfm"

// func GetClientCookie(val string) *http.Cookie {

// 	secure := false //TODO SET TRUE FOR PROD
// 	cookieReturn := http.Cookie{
// 		Name:  "sid",
// 		Value: val,
// 		Path:  "/",
// 		// MaxAge:   0,
// 		Expires:  time.Now().Add(24 * time.Hour), //Cookie expires in 24 hours
// 		Secure:   secure,
// 		HttpOnly: true,
// 		SameSite: http.SameSiteLaxMode,
// 	}
// 	return &cookieReturn
// }

type Tx struct {
	SessIDVerifier string
	CreatedAt      time.Time
	IP             string
}

type authService struct {
	repo         AuthRepository
	lastFMAPIKey string
	lfmapi       musicapi.LastFMAPIExternal
	config       Config
}

type Config struct {
	ExpiryDuration     time.Duration
	FrontendCookieName string
	FrontendURL        string
}

type SessionID string

func NewAuthService(repo AuthRepository, cfg Config) AuthService {
	return &authService{
		repo:         repo,
		lastFMAPIKey: os.Getenv("LASTFM_API_KEY"),
		config:       cfg,
	}
}

type AuthService interface {
	GetDeletedCookie(cookieName string) *http.Cookie
	CheckCookieValidity(ctx context.Context, cookieValue string) (bool, error)
	GenerateNewTx(userIP string) *Tx
	GetInitLoginURL(state string) string
	GetNewWebSessionURL(token string) (string, url.Values)

	//Sets new temp state (with expiry) and maps to new permanent sessionID.
	GenerateNewStateAndSID(ctx context.Context) (string, string, error)

	ConsumeStateAndReturnSID(ctx context.Context, state string) (string, error)

	//SetSessionKey(ctx context.Context, sessionID, userKey string) error
	DelSidKey(ctx context.Context, sessionID string) error
	MakeNewUser(context context.Context, validationSid string, userName string) (uuid.UUID, error)
	GetUserByValidSessionID(context context.Context, sid string) (string, error)
	IsSIDValid(context context.Context, sid string) (bool, error)
}

func (s *authService) MakeNewUser(context context.Context, validationSid string, userName string) (uuid.UUID, error) {

	newuuid := uuid.New()
	err := s.repo.MakeNewUser(context, validationSid, userName, newuuid)

	//Download blend data

	return newuuid, err

}

func (s *authService) GetUserByAnySessionID(context context.Context, sid string) (string, error) {
	return s.repo.GetUserByAnySessionID(context, sid)
}

func (s *authService) GetUserByValidSessionID(context context.Context, sid string) (string, error) {
	return s.repo.GetUserByValidSessionID(context, sid, s.config.ExpiryDuration)
}

func (s *authService) IsSIDValid(context context.Context, sid string) (bool, error) {
	u, err := s.repo.GetUserByValidSessionID(context, sid, s.config.ExpiryDuration)
	return u != "", err
}

func (s *authService) DeleteSessionID(context context.Context, sid string) error {

	return s.repo.DeleteSingularSessionID(context, sid)

}

func (s *authService) DeleteUser(context context.Context, userid string) error {
	return s.repo.DeleteUser(context, userid)
}

func (s *authService) NewSid() SessionID {
	return SessionID(uuid.New().String())
}

// func (s *authService) SetSessionKey(ctx context.Context, sessionID, userKey string) error {

// 	err := s.repo.SetSidWebSesssionKey(ctx, sessionID, userKey, time.Hour*24*10) //Set for 10 days
// 	if err != nil {
// 		return fmt.Errorf("Cannot set session key in repository: %v", err)
// 	}
// 	return nil
// }

func (s *authService) GetDeletedCookie(cookieName string) *http.Cookie {

	secure := false //TODO Set true for PROD
	cookieReturn := http.Cookie{
		Name:  cookieName,
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

func (s *authService) CheckCookieValidity(ctx context.Context, sidValue string) (bool, error) {

	found, err := s.repo.GetSidKey(ctx, sidValue)
	if err != nil { //Err checking
		return false, err
	} else if found != "" { //Found valid key
		return true, nil
	} else { //Did not find any key associated
		return false, nil
	}
}

func (s *authService) DelSidKey(ctx context.Context, sessionID string) error {
	err := s.repo.DelSidKey(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("Cannot delete SID key from repository: %v", err)
	}
	return nil
}

func (s *authService) ConsumeStateAndReturnSID(ctx context.Context, state string) (string, error) {
	sid, err := s.repo.ConsumeStateSID(ctx, state)
	if err != nil {
		glog.Warning("Cannot consume state token from repository: %v", err)
		return "", err
	}
	return sid, nil

}

func makeNewState() (string, error) {
	//Make new state
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	newState := hex.EncodeToString(randomBytes)
	if err != nil {
		glog.Warning("Cannot create state token")
		return "", err
	}
	return newState, nil
}

func makeNewSID() string {
	sessIDVerifier := uuid.New().String()
	return sessIDVerifier
}

func (s *authService) GenerateNewStateAndSID(ctx context.Context) (string, string, error) {

	//Make new state and SID
	newState, err := makeNewState()
	if err != nil {
		glog.Warning("Cannot create state token")
		return "", "", err
	}

	sessID := makeNewSID()

	//Save to Repository
	err = s.repo.SetNewStateSid(ctx, newState, sessID, s.config.ExpiryDuration)
	if err != nil {
		glog.Warning("Cannot save state-SID to repository: %v", err)
		return "", "", err
	}

	return sessID, newState, nil
}

func (s *authService) GetInitLoginURL(state string) string {
	q := url.Values{}
	q.Set("api_key", s.lastFMAPIKey)
	q.Set("cb", string(LASTFM_CALLBACK+"?state="+state))

	requestURL := LASTFM_BASE_AUTH_API + "?" + q.Encode()
	// glog.Info(requestURL)
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

func (s *authService) GetNewWebSessionURL(token string) (string, url.Values) {

	q := url.Values{}
	q.Set("method", "auth.getSession")
	q.Set("api_key", s.lastFMAPIKey)
	q.Set("token", token)
	q.Set("api_sig", getSessionAPISignature(s.lastFMAPIKey, token))

	// requestURL := LASTFM_ROOT_API + "?" + q.Encode()
	return LASTFM_ROOT_API, q //Didn't realise POST would require them as different //TODO Make this cleaner
}

func (s *authService) GenerateNewTx(userIP string) *Tx {
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
