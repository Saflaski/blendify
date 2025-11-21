package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
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
