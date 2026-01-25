package musicbrainz_test

import (
	"backend-lastfm/internal/musicbrainz"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgresMusicBrainzRepo(t *testing.T) {
	godotenv.Load("../../.env")
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("godotenv.Load failed")
	}

	db_dsn, exists := os.LookupEnv("MUSICBRAINZ_DB_DSN")
	if !exists {
		t.Fatal("MUSICBRAINZ_DB_DSN not set in environment")
		return
	}
	db, err := sqlx.Connect("pgx", db_dsn)
	assert.NoError(t, err)
	defer db.Close()

	repo := musicbrainz.NewPostgresMusicBrainzRepo(db)

	assert.NotNil(t, repo.Genre)

	//Test GetGenresByRecording
	t.Run("Get Genres by some recording", func(t *testing.T) {
		genreObj, err := repo.Genre.GetGenreByRecordings(t.Context(), []string{"5afa33bb-83f6-42c6-b789-8bfa34c48b50"})
		for k, v := range genreObj {
			t.Log("Recording MBID:", k)
			// t.Logf("Recording MBID: %s, Genres: %+v\n", k, v)
			for _, genre := range v {
				t.Logf(" - Genre: %s (Tag Count: %d)\n", genre.Name, genre.TagCount)
			}
		}
		assert.NotEmpty(t, genreObj)
		assert.NoError(t, err)
	})

	t.Run("Get Recording Candidates by closest match to recording and artist name", func(t *testing.T) {
		names := []string{
			"Faint",
			"Numb",
			"Papercut",
			"Show me how",
			"Glamorous",
			"Bags",
		}
		artists := []string{
			"Linkin Park",
			"Linkin Park",
			"Linkin Park",
			"Men I Trust",
			"Fergie feat. Ludacris",
			"Clairo",
		}

		mbids := []string{
			"54a3c21c-5395-44a2-b90b-b7fab8095c20",
			"352dd518-23cd-4c5a-9551-ba02097b177b",
			"9aa621e1-46f2-4c91-8111-741583985612",
			"dbad831a-7a9d-416e-9ef0-11e740fef6a0",
			"a1cb2d52-7702-4a96-a5ed-41281995dbdb",
			"9433df6b-037b-41f8-9edf-0e8c9ffaf390",
		}

		results, err := repo.Recording.GetClosestRecordings(t.Context(), names, artists)
		if err != nil {
			t.Fatalf("GetClosestRecordings failed: %v", err)

		}

		if len(results) != len(names) {
			t.Log("Results:", results)
			t.Fatalf("expected %d results, got %d", len(names), len(results))
		}

		for i, r := range results {
			if r.RecordingMBID == "" {
				t.Errorf("result %d has empty MBID", i)
			}

			if r.RecordingName == "" {
				t.Errorf("result %d has empty recording name", i)
			}

			if r.ArtistName == "" {
				t.Errorf("result %d has empty artist name", i)
			}

			assert.Equal(t, mbids[i], r.RecordingMBID, "mbid mismatch for result %d", i)
		}
		for _, r := range results {
			t.Logf("Recording: %s by %s (MBID: %s)\n", r.RecordingName, r.ArtistName, r.RecordingMBID)
		}
	})
}
