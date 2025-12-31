package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func TestAuthService(t *testing.T) {

	godotenv.Load("../../.env")
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("godotenv.Load failed")
	}
	godotenv.Load("../../.env.test")
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Fatal("godotenv.Load failed [TEST]")
	}

	DB_ADDR := os.Getenv("DB_ADDR")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NUM, _ := strconv.Atoi(os.Getenv("DB_NUM"))
	DB_PROTOCOL, _ := strconv.Atoi(os.Getenv("DB_PROTOCOL"))
	LASTFM_API_KEY := os.Getenv("LASTFM_API_KEY")
	// LAST_FM_URL := "https://ws.audioscrobbler.com/2.0/"
	if len(DB_ADDR) == 0 || len(LASTFM_API_KEY) == 0 {
		t.Errorf("key Environment Value is empty")
	}

	// WS_KEY := os.Getenv("WEB_SESSION_KEY")

	redisStore := NewRedisStateStore(redis.NewClient(&redis.Options{
		Addr:     DB_ADDR,
		Password: DB_PASS,
		DB:       0,
		Protocol: DB_PROTOCOL,
	}), time.Duration(648000)*time.Second)
	_ = DB_NUM

	// lfmApi := musicapi.NewLastFMExternalAdapter(LASTFM_API_KEY, LAST_FM_URL, true)
	authService := AuthService{
		repo:         redisStore,
		lastFMAPIKey: LASTFM_API_KEY,
		// config: auth.Config{
		// 	ExpiryDuration: time.Duration(app.config.sessionExpiry) * time.Second,
		// 	// ExpiryDuration:     time.Duration(app.config.sessionExpiry) * time.Second,
		// 	FrontendCookieName: "sid",
		// 	FrontendURL:        os.Getenv("FRONTEND_URL"),
		// 	BackendURL:         os.Getenv("BACKEND_URL"),
		// },
	}

	ctx := context.Background()
	t.Run("Create New User", func(t *testing.T) {
		sid := uuid.New().String()
		_, err := authService.MakeNewUser(ctx, sid, "saflas")
		if err != nil {
			t.Errorf("could not create new user:%s", err)
		}
		// t.Logf("userid: %s", userid)

	})

	t.Run("Simulate account creation then device logout", func(t *testing.T) {
		//Create New User
		newSid := uuid.New().String()
		newUserId, err := authService.MakeNewUser(ctx, newSid, "saflas2")
		if err != nil {
			t.Errorf("could not create new user:%s", err)
		}
		// t.Log("Success new user\n")
		newUserIdStr := newUserId.String()

		//Find User with given sid on /logout
		returnedUserId, err := authService.GetUserByAnySessionID(ctx, newSid) //Get userid
		if err != nil {
			t.Errorf("could not read user: %s", err)
		}

		if returnedUserId != newUserIdStr {
			t.Errorf("Returned userid did not match. Want %s, got %s", newUserIdStr, returnedUserId)
		}
		// t.Log("Success correct user returned\n")

		//Delete sid with returned user
		authService.DeleteSessionID(ctx, newSid)

		returnedUserId, err = authService.GetUserByAnySessionID(ctx, newSid) //Get userid
		if err != nil {
			t.Errorf("could not read user: %s", err)
		}

		if returnedUserId != "" {
			t.Errorf("Could not delete user. Expected %s, got %s", "", returnedUserId)
		}

	})

	t.Run("Simulate account creation then deletion", func(t *testing.T) {
		newSid := uuid.New().String()
		newUserId, err := authService.MakeNewUser(ctx, newSid, "saflas2")
		if err != nil {
			t.Errorf("could not create new user:%s", err)
		}
		// t.Log("Success new user\n")
		newUserIdStr := newUserId.String()

		//User issues logout with and we get sid from cookie
		returnedUserId, err := authService.GetUserByAnySessionID(ctx, newSid) //Get userid
		if err != nil {
			t.Errorf("could not read user: %s", err)
		}

		if returnedUserId != newUserIdStr {
			t.Errorf("Returned userid did not match. Want %s, got %s", newUserIdStr, returnedUserId)
		}
		// t.Log("Success correct user returned\n")

		err = authService.DeleteUser(ctx, newUserIdStr)
		if err != nil {
			t.Errorf("Error in deleting user: %s", err)
		}

		// t.Log("Success deleted user successfully\n")

	})

	t.Run("get lfm by userid", func(t *testing.T) {
		str, err := authService.repo.GetLFMByUserId(t.Context(), "dc2e4fcf-0d07-4871-b287-9b3488599c3d")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(str)
	})

}
