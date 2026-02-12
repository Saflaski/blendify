package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"backend-lastfm/internal/musicbrainz"
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type StubBlendService struct {
}

func TestBlend(t *testing.T) {
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
	BlendifysqlxDB := sqlx.MustConnect("pgx", os.Getenv("BLENDIFY_DB_DSN"))
	BlendifysqlxDB.SetMaxOpenConns(25)
	BlendifysqlxDB.SetMaxIdleConns(25)
	BlendifysqlxDB.SetConnMaxLifetime(5 * time.Minute)

	redisStore := NewBlendStore(redis.NewClient(&redis.Options{
		Addr:     DB_ADDR,
		Password: DB_PASS,
		DB:       2,
		Protocol: DB_PROTOCOL,
	}),
		BlendifysqlxDB,
	)

	lfm_adapter := musicapi.NewLastFMExternalAdapter(
		LASTFM_API_KEY,
		"https://ws.audioscrobbler.com/2.0/",
		true,
		200,
	)
	sqlxDB := sqlx.MustConnect("pgx", os.Getenv("MUSICBRAINZ_DB_DSN"))
	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(25)
	sqlxDB.SetConnMaxLifetime(5 * time.Minute)

	mbRepo := musicbrainz.NewPostgresMusicBrainzRepo(sqlxDB)
	blendService := NewBlendService(*redisStore, *lfm_adapter, *musicbrainz.NewMBService(mbRepo))
	_ = blendService
	_ = redisStore
	//Mock Data

	t.Run("Get User - Top 50 Listened to Artists -- 3 months", func(t *testing.T) {
		userName := "test2002"
		// response, err := blendService.getTopArists(userName, BlendTimeDurationThreeMonth)
		TopArtistResponse, err := lfm_adapter.GetUserTopArtists(
			t.Context(),
			userName,
			string(BlendTimeDurationThreeMonth),
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

	// t.Run("Get blend between two users: saflas and test2002", func(t *testing.T) {
	// 	userA := "saflas"
	// 	userB := "test2002"

	// 	blendNumber, err := blendService.getArtistBlend(t.Context(), userA, userB, BlendTimeDurationOneMonth)
	// 	if err != nil {
	// 		t.Errorf("Error during getting artist blend: %q", err)
	// 	}

	// 	if blendNumber > 100 || blendNumber <= 0 { //I know the blend between these two isn't 0
	// 		t.Errorf("Number is not within acceptable range: %d", blendNumber)
	// 	}
	// })

	t.Run("Create and Delete all blends by user", func(t *testing.T) {
		// userName := "internaltesting"
		// response, err := blendService.getTopArists(userName, BlendTimeDurationThreeMonth)
		ctx := t.Context()
		ctx = context.WithValue(ctx, "user", "saflas")
		ctx2 := t.Context()
		ctx = context.WithValue(ctx, "user", "other")
		user := userid("123-123-123-123-123")
		user2 := userid("456-456-456-456-456")

		blendService.repo.redisClient.HSet(ctx2, "user:", "LFM Username", string(user2))

		link, err := blendService.GenerateNewLinkAndAssignToUser(ctx, user)
		if err != nil {
			t.Fatalf("Generating new link and assigning to userid: %s", err)
		}

		blendid, err := blendService.AddOrMakeBlendFromLink(ctx2, user2, link)
		if err != nil || blendid == "" {
			t.Fatalf("Making Blend from invite : %s", err)
		}
		t.Log("blendid")
		t.Log(blendid)
		// _, err = blendService.GetDuoBlendData(ctx, blendid)
		// if err != nil {
		// 	t.Fatalf("Duo Blend: %s", err)
		// }

		t.Log("Making sure user blends were deleted")

		blendids, err := blendService.repo.GetBlendsByUser(ctx2, user2)
		if err != nil {
			t.Fatal(err)
		}

		if len(blendids) == 0 {
			t.Fatalf("Could not make blends %s", err)
		}

		err = blendService.DeleteUserBlends(ctx2, string(user2))
		if err != nil {
			t.Fatalf(" Couldnt delete: %s", err)
		}

		blendids, err = blendService.repo.GetBlendsByUser(ctx2, user2)
		if err != nil {
			t.Fatal(err)
		}

		if len(blendids) != 0 {
			t.Fatalf("Could not delete blends %s", err)
		}

		blendService.repo.redisClient.HDel(ctx2, "user:", "LFM Username", string(user2))
		blendService.repo.redisClient.HDel(ctx2, "user:", "LFM Username", string(user))

	})

	t.Run("Create and Delete user and try deleting 0 blends", func(t *testing.T) {
		// userName := "internaltesting"
		// response, err := blendService.getTopArists(userName, BlendTimeDurationThreeMonth)
		ctx := t.Context()
		ctx = context.WithValue(ctx, "user", "saflas")
		ctx2 := t.Context()
		ctx = context.WithValue(ctx, "user", "other")
		user := userid("123-123-123-123-123")
		user2 := userid("456-456-456-456-456")

		blendService.repo.redisClient.HSet(ctx2, "user:", "LFM Username", string(user2))

		t.Log("Making sure user blends were deleted")

		blendids, err := blendService.repo.GetBlendsByUser(ctx2, user2)
		if err != nil {
			t.Fatal(err)
		}

		if len(blendids) != 0 {
			t.Fatalf("Blends already exist for user? %s", err)
		}

		err = blendService.DeleteUserBlends(ctx2, string(user2))
		if err != nil {
			t.Fatalf(" Couldnt delete: %s", err)
		}

		blendids, err = blendService.repo.GetBlendsByUser(ctx2, user2)
		if err != nil {
			t.Fatal(err)
		}

		if len(blendids) != 0 {
			t.Fatalf("Could not delete blends %s", err)
		}

		blendService.repo.redisClient.HDel(ctx2, "user:", "LFM Username", string(user2))
		blendService.repo.redisClient.HDel(ctx2, "user:", "LFM Username", string(user))

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
	BlendifysqlxDB := sqlx.MustConnect("pgx", os.Getenv("BLENDIFY_DB_DSN"))
	BlendifysqlxDB.SetMaxOpenConns(25)
	BlendifysqlxDB.SetMaxIdleConns(25)
	BlendifysqlxDB.SetConnMaxLifetime(5 * time.Minute)

	blendStore := NewBlendStore(redis.NewClient(&redis.Options{
		Addr:     DB_ADDR,
		Password: DB_PASS,
		DB:       0,
		Protocol: DB_PROTOCOL,
	}),
		BlendifysqlxDB,
	)

	lfm_adapter := musicapi.NewLastFMExternalAdapter(
		LASTFM_API_KEY,
		"https://ws.audioscrobbler.com/2.0/",
		true,
		200,
	)

	blendService := NewBlendService(*blendStore, *lfm_adapter, *musicbrainz.NewMBService(musicbrainz.NewPostgresMusicBrainzRepo(nil)))
	_ = blendService
	_ = blendStore

	t.Run("Hydrate and cache user", func(t *testing.T) {
		blendService.GetNewDataForUser(context.Background(), userid("dc2e4fcf-0d07-4871-b287-9b3488599c3d"))
	})

	t.Run("Try to get from cache or download", func(t *testing.T) {
		// err = blendService.PopulateUsersByBlend(t.Context(), blendId)
		// if err != nil {
		// 	return DuoBlend{}, fmt.Errorf(" Could not populate user data: %w", err)
		// }
		userA := userid("3c7a687b-e8df-4f13-ad94-6bee68d67aa1")
		userB := userid("dc2e4fcf-0d07-4871-b287-9b3488599c3d")
		blendService.GetNewDataForUser(t.Context(), userA)
		blendService.GetNewDataForUser(t.Context(), userB)
		blend, err := blendService.GenerateBlendOfTwo(t.Context(),
			"0000", //mock jobid
			userB,
			userA,
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

func TestMBService(t *testing.T) {
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

	BlendifysqlxDB := sqlx.MustConnect("pgx", os.Getenv("BLENDIFY_DB_DSN"))
	BlendifysqlxDB.SetMaxOpenConns(25)
	BlendifysqlxDB.SetMaxIdleConns(25)
	BlendifysqlxDB.SetConnMaxLifetime(5 * time.Minute)

	blendStore := NewBlendStore(redis.NewClient(&redis.Options{
		Addr:     DB_ADDR,
		Password: DB_PASS,
		DB:       2,
		Protocol: DB_PROTOCOL,
	}),
		BlendifysqlxDB,
	)

	lfm_adapter := musicapi.NewLastFMExternalAdapter(
		LASTFM_API_KEY,
		"https://ws.audioscrobbler.com/2.0/",
		true,
		200,
	)
	sqlxDB := sqlx.MustConnect("pgx", os.Getenv("MUSICBRAINZ_DB_DSN"))
	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(25)
	sqlxDB.SetConnMaxLifetime(5 * time.Minute)

	mbRepo := musicbrainz.NewPostgresMusicBrainzRepo(sqlxDB)
	blendService := NewBlendService(*blendStore, *lfm_adapter, *musicbrainz.NewMBService(mbRepo))
	_ = blendService
	_ = blendStore
	//Mock Data

	t.Run("Test Populate MapCatStats with MBID Closest search", func(t *testing.T) {
		trackToPlays := map[string]CatalogueStats{}

		// TIMEDURATION DOESNT WORK
		// ------------------------
		topTracks, err := lfm_adapter.GetUserTopTracks(
			t.Context(),
			"saflas",
			string(BlendTimeDurationYear),
			6,
			50,
		)

		if len(topTracks.TopTracks.Track) == 0 {
			t.Error("got empty track list from lastfm adapter")
		}

		if err != nil {
			t.Error("error during get user top tracks: %w", err)
		}
		// for _, v := range topTracks.TopTracks.Track {
		// 	playcount, err := strconv.Atoi(v.Playcount)
		// 	if err != nil {
		// 		return trackToPlays, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		// 	}
		// 	trackToPlays[v.Name] = playcount
		// }

		for _, v := range topTracks.TopTracks.Track {

			playcount, err := strconv.Atoi(v.Playcount)
			if err != nil {
				t.Errorf("got unparseable string during string -> int conversation: %v", err)
			}

			// // imageURL := getCatalogueImageURL(v.LFMImages) //Selects a good pic out of the ones given
			// genreObjects, err := blendService.MBService.GetGenresByRecordingMBID(t.Context(), v.MBID)
			// if err != nil {
			// 	t.Errorf("could not get genres during gettoptracks: %v", err)

			// }
			// genres := make([]string, len(genreObjects))
			// for i, g := range genreObjects {
			// 	genres[i] = g.Name //Capitalize first letter of each genre
			// }
			catStat := CatalogueStats{
				Artist:      v.Artist,
				Count:       playcount,
				PlatformURL: v.URL,
				Image:       "",
				PlatformID:  v.MBID,
				// Genres:      genres,
			}
			trackToPlays[v.Name] = catStat

		}
		// ------------------------
		t.Log("Length of MapCatStats: ", len(trackToPlays))
		trackToPlays, err = blendService.PopulateTrackMBIDs(t.Context(), trackToPlays)
		if err != nil {

			t.Errorf(" could not populate mbids for map cat stats: %v", err)
		}
		t.Log("Populated CatalogueStats with MBIDs")
		t.Log("Populating with Genres")
		trackToPlays, err = blendService.PopulateGenresForMapCatStats(t.Context(), trackToPlays, "recording")
		assert.NoError(t, err)

		t.Log("Populated CatalogueStats with Genres")

		t.Log("Printing populated CatalogueStats with MBIDs")
		i := 0
		sum := 0
		genreSum := 0
		for trackName, catStat := range trackToPlays {
			i += 1
			if catStat.PlatformID == "" {
				t.Logf("Could not populate MBID for track: %s", trackName)
			} else {
				sum += 1
			}

			if len(catStat.Genres) == 0 {
				t.Logf("Could not populate Genres for track: %s, MBID: %s", trackName, catStat.PlatformID)
			} else {
				genreSum += 1
			}

			t.Logf("MBID: %s, Genres: %s\n", catStat.PlatformID, catStat.Genres)
		}
		t.Logf("Successfully populated %d out of %d MBIDs", sum, i)
		t.Logf("Successfully populated %d out of %d Genres", genreSum, i)
	})

	t.Run("Calculating common genres with placeholders", func(t *testing.T) {
		userAGenres := []string{"Rock", "Pop", "Jazz", "Classical", "Hip Hop", "Electronic"}
		userBGenres := []string{"Rock", "Pop", "Country", "Classical", "Reggae", "Electronic"}

		allUserGenres := make([][]string, 2)

		allUserGenres[0] = userAGenres
		allUserGenres[1] = userBGenres
		commonGenres, err := blendService.CalculateIntersectionOfStringSlices(allUserGenres)
		assert.NoError(t, err)
		t.Log("Common Genres:")
		for _, genre := range commonGenres {
			t.Log(genre)
		}
	})

}
