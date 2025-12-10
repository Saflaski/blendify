package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type StubBlendService struct {
}

func TestGetBlend(t *testing.T) {
	godotenv.Load("../../.env")
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("godotenv.Load failed")
	}

	DB_ADDR := os.Getenv("DB_ADDR")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NUM, _ := strconv.Atoi(os.Getenv("DB_NUM"))
	DB_PROTOCOL, _ := strconv.Atoi(os.Getenv("DB_PROTOCOL"))
	LASTFM_API_KEY := os.Getenv("LASTFM_API_KEY")

	if len(DB_ADDR) == 0 || len(LASTFM_API_KEY) == 0 {
		t.Errorf("key Environment Value is empty")
	}

	redisStore := NewRedisStateStore(redis.NewClient(&redis.Options{
		Addr:     DB_ADDR,
		Password: DB_PASS,
		DB:       DB_NUM,
		Protocol: DB_PROTOCOL,
	}))

	lfm_adapter := musicapi.NewLastFMExternalAdapter(
		LASTFM_API_KEY,
		"https://ws.audioscrobbler.com/2.0/",
		true,
	)
	blendService := NewBlendService(*redisStore, *lfm_adapter)
	_ = blendService
	_ = redisStore
	//Mock Data

	t.Run("Get User - Top 50 Listened to Artists -- 3 months", func(t *testing.T) {
		userName := "test2002"
		// response, err := blendService.getTopArists(userName, BlendTimeDurationThreeMonth)
		TopArtistResponse, err := lfm_adapter.GetUserTopArtists(
			userName,
			durationMap[BlendTimeDurationThreeMonth],
			1,
			50,
		)

		if err != nil {
			t.Errorf("Expected no error, got %q", err)
		}

		mostListenedToArtist := TopArtistResponse.TopArtists.Artist[0].Name
		if mostListenedToArtist == "" {
			t.Errorf("Got empty string after processing LastFM response: %q", mostListenedToArtist)

		}
		if len(TopArtistResponse.TopArtists.Artist) == 0 {
			t.Errorf("empty error from processing topartists")
		}

	})

	t.Run("Get blend between two users: saflas and test2002", func(t *testing.T) {
		userA := "saflas"
		userB := "test2002"

		blendNumber, err := blendService.getArtistBlend(userA, userB, BlendTimeDurationOneMonth)
		if err != nil {
			t.Errorf("Error during getting artist blend: %q", err)
		}

		if blendNumber > 100 || blendNumber <= 0 { //I know the blend between these two isn't 0
			t.Errorf("Number is not within acceptable range: %d", blendNumber)
		}
	})

}

func TestDownloadAndCache(t *testing.T) {
	godotenv.Load("../../.env")
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("godotenv.Load failed")
	}

	DB_ADDR := os.Getenv("DB_ADDR")
	DB_PASS := os.Getenv("DB_PASS")
	// DB_NUM, _ := strconv.Atoi(os.Getenv("DB_NUM"))
	DB_PROTOCOL, _ := strconv.Atoi(os.Getenv("DB_PROTOCOL"))
	LASTFM_API_KEY := os.Getenv("LASTFM_API_KEY")

	if len(DB_ADDR) == 0 || len(LASTFM_API_KEY) == 0 {
		t.Errorf("key Environment Value is empty")
	}

	redisStore := NewRedisStateStore(redis.NewClient(&redis.Options{
		Addr:     DB_ADDR,
		Password: DB_PASS,
		DB:       0,
		Protocol: DB_PROTOCOL,
	}))

	lfm_adapter := musicapi.NewLastFMExternalAdapter(
		LASTFM_API_KEY,
		"https://ws.audioscrobbler.com/2.0/",
		true,
	)

	blendService := NewBlendService(*redisStore, *lfm_adapter)
	_ = blendService
	_ = redisStore

	t.Run("Hydrate and cache user", func(t *testing.T) {
		blendService.GetNewDataForUser(context.Background(), userid("dc2e4fcf-0d07-4871-b287-9b3488599c3d"))
	})

	t.Run("Try to get from cache or download", func(t *testing.T) {
		blend, err := blendService.GenerateBlendOfTwo(t.Context(),
			userid("dc2e4fcf-0d07-4871-b287-9b3488599c3d"),
			userid("3c7a687b-e8df-4f13-ad94-6bee68d67aa1"),
		)
		if err != nil {
			t.Errorf("%s", err)

		}
		t.Log(blend.AlbumBlend.OneYear)
		t.Log(blend.AlbumBlend.ThreeMonth)
		t.Log(blend.AlbumBlend.OneMonth)
		t.Log(blend.ArtistBlend.OneYear)
		t.Log(blend.ArtistBlend.ThreeMonth)
		t.Log(blend.ArtistBlend.OneMonth)
		t.Log(blend.TrackBlend.OneYear)
		t.Log(blend.TrackBlend.ThreeMonth)
		t.Log(blend.TrackBlend.OneMonth)

	})

	t.Run("Get LFM by user", func(t *testing.T) {
		user, err := blendService.repo.GetLFMByUserId(t.Context(), "dc2e4fcf-0d07-4871-b287-9b3488599c3d")
		t.Log(user)
		if err != nil {
			t.Error(err)
		}
	})
}
