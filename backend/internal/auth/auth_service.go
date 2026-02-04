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
	"github.com/redis/go-redis/v9"
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

type AuthService struct {
	repo         AuthRepository
	lastFMAPIKey string
	lfmapi       musicapi.LastFMAPIExternal
	config       Config
}

type Config struct {
	ExpiryDuration     time.Duration
	FrontendCookieName string
	FrontendURL        string
	BackendURL         string
}

type SessionID string

func NewAuthService(repo AuthRepository, cfg Config) *AuthService {
	return &AuthService{
		repo:         repo,
		lastFMAPIKey: os.Getenv("LASTFM_API_KEY"),
		config:       cfg,
	}
}

type authService interface {
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
	GetUserByLFMUsername(context context.Context, lfmName string) (string, error)
}

func (s *AuthService) MakeNewUser(context context.Context, validationSid string, userName string) (uuid.UUID, error) {

	res, err := s.CheckIfExistingUserFromLFM(context, userName)
	if err != nil {
		return res, err //res will be empty value
	}

	if res == uuid.Nil { //We did not find an existing username
		// glog.Info("Generating full new user")
		newuuid := uuid.New()
		err = s.repo.MakeNewUser(context, validationSid, userName, newuuid)

		//Download blend data

		return newuuid, err
	} else {
		//Assign new SID to user
		// glog.Info("Found existing user, just assigning sid to existing user then")
		err := s.repo.AddNewSidToExistingUser(context, res, validationSid)
		if err != nil {
			return res, fmt.Errorf(" could not add new sid to existing user: %s with err: %w", res.String(), err)
		}
		return res, nil

	}

}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	user, ok := ctx.Value(UserKey).(string)
	if !ok {
		return "", fmt.Errorf(" did not find userid in context")
	}
	return user, nil
}

// This will return err if and only if there is an error with checking the userName
// If it does not find an existing user, it will return an empty UUID value
func (s *AuthService) CheckIfExistingUserFromLFM(context context.Context, userName string) (uuid.UUID, error) {

	res, err := s.repo.GetUserIdByLFM(context, userName)
	if err != redis.Nil && err != nil {
		glog.Error("actual full error")
		return uuid.Nil, err
	} else if err == redis.Nil {
		glog.Error("Did not find user")
		return uuid.Nil, nil
	} else {
		uuidAsBytes, err := uuid.Parse(res)
		if err != nil {
			glog.Errorf("Found actual user but could not parse, %w", err)
			return uuid.Nil, err
		}
		glog.Error("Found actual user and returning")
		return uuidAsBytes, err
	}
}

func (s *AuthService) GetUserByAnySessionID(context context.Context, sid string) (string, error) {
	return s.repo.GetUserByAnySessionID(context, sid)
}

func (s *AuthService) GetUserByLFMUsername(context context.Context, lfmName string) (string, error) {
	return s.repo.GetUserIdByLFM(context, lfmName)
	// if err != nil {
	// 	glog.Errorf("error during getting userid by lfmName sec index, %w", err)
	// }
	// if err == redis.Nil {
	// 	return "", nil
	// }
	// return res, err
}

func (s *AuthService) AddUserIdToLFMIndex(context context.Context, userid, lfmName string) {
	err := s.repo.AddUserIdToLFMIndex(context, userid, lfmName)
	if err != nil {
		glog.Errorf("Could not add user to LFM-User index in repo, %w", err)
	}
}

func (s *AuthService) GetUserByValidSessionID(context context.Context, sid string) (string, error) {
	return s.repo.GetUserByValidSessionID(context, sid, s.config.ExpiryDuration)
}

func (s *AuthService) IsSIDValid(context context.Context, sid string) (bool, error) {
	u, err := s.repo.GetUserByValidSessionID(context, sid, s.config.ExpiryDuration)
	return u != "", err
}

func (s *AuthService) DeleteSessionID(context context.Context, sid string) error {

	return s.repo.DeleteSingularSessionID(context, sid)

}

func (s *AuthService) DeleteUser(context context.Context, userid string) error {
	return s.repo.DeleteUser(context, userid)
}

func (s *AuthService) NewSid() SessionID {
	return SessionID(uuid.New().String())
}

// func (s *AuthService) SetSessionKey(ctx context.Context, sessionID, userKey string) error {

// 	err := s.repo.SetSidWebSesssionKey(ctx, sessionID, userKey, time.Hour*24*10) //Set for 10 days
// 	if err != nil {
// 		return fmt.Errorf("Cannot set session key in repository: %v", err)
// 	}
// 	return nil
// }

func (s *AuthService) GetDeletedCookie(cookieName string) *http.Cookie {

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

func (s *AuthService) CheckCookieValidity(ctx context.Context, sidValue string) (bool, error) {

	found, err := s.repo.GetSidKey(ctx, sidValue)
	if err != nil { //Err checking
		return false, err
	} else if found != "" { //Found valid key
		return true, nil
	} else { //Did not find any key associated
		return false, nil
	}
}

func (s *AuthService) DelSidKey(ctx context.Context, sessionID string) error {
	err := s.repo.DelSidKey(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("Cannot delete SID key from repository: %v", err)
	}
	return nil
}

func (s *AuthService) ConsumeStateAndReturnSID(ctx context.Context, state string) (string, error) {
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

func (s *AuthService) GenerateNewStateAndSID(ctx context.Context) (string, string, error) {

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

func (s *AuthService) GetInitLoginURL(state string) string {
	q := url.Values{}
	q.Set("api_key", s.lastFMAPIKey)
	q.Set("cb", string(s.config.BackendURL+"/auth/callback/lastfm"+"?state="+state))

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

func (s *AuthService) GetNewWebSessionURL(token string) (string, url.Values) {

	q := url.Values{}
	q.Set("method", "auth.getSession")
	q.Set("api_key", s.lastFMAPIKey)
	q.Set("token", token)
	q.Set("api_sig", getSessionAPISignature(s.lastFMAPIKey, token))

	// requestURL := LASTFM_ROOT_API + "?" + q.Encode()
	return LASTFM_ROOT_API, q //Didn't realise POST would require them as different //TODO Make this cleaner
}

func (s *AuthService) GenerateNewTx(userIP string) *Tx {
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
