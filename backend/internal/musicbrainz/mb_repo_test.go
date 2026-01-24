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
}
