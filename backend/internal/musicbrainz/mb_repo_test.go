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
		testStrArray := []string{"54a3c21c-5395-44a2-b90b-b7fab8095c20", "a892a0de-3119-49e3-8237-587537dad4d9"}
		genreObj, err := repo.Genre.GetGenreByRecordings(t.Context(),
			testStrArray)
		if err != nil {
			t.Errorf(": %v", err)
		}
		for k, v := range genreObj {
			t.Log("Recording MBID:", k)
			// t.Logf("Recording MBID: %s, Genres: %+v\n", k, v)
			for _, genre := range v {
				t.Logf(" - Genre: %s (Tag Count: %d)\n", genre.Name, genre.TagCount)
			}
		}
		assert.Equal(t, len(testStrArray), len(genreObj))
		assert.NotEmpty(t, genreObj)
		assert.NoError(t, err)
	})

	t.Run("Get Recording Candidates by closest match to recording and artist name - Non Fuzzy", func(t *testing.T) {
		names := []string{
			"Faint",
			"Numb",
			"Papercut",
			"Show Me How",
			"Glamorous",
			"Bags",
			"Claire",
		}
		artists := []string{
			"Linkin Park",
			"Linkin Park",
			"Linkin Park",
			"Men I Trust",
			"Fergie feat. Ludacris",
			"Clairo",
			"Déyyess",
		}

		mbids := []string{
			"54a3c21c-5395-44a2-b90b-b7fab8095c20",
			"352dd518-23cd-4c5a-9551-ba02097b177b",
			"9aa621e1-46f2-4c91-8111-741583985612",
			"dbad831a-7a9d-416e-9ef0-11e740fef6a0",
			"74ff24ec-0a35-4093-bab9-9e8bbcae11ac",
			"9433df6b-037b-41f8-9edf-0e8c9ffaf390",
			"356c0fa3-48a9-4d4e-b3d2-616b927a0e60",
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

	t.Run("Get Recording Bulk Candidates by closest match to recording and artist name", func(t *testing.T) {
		names := []string{
			"Numb",
			"Faint",
			"Papercut",
			"Show me how",
			"Glamorous",
			"Bags",
			"Claire",
			"C.R.E.A.M. (Cash Rules Everything Around Me)",
			"In the End",
			"Breaking the Habit",
			"Somewhere I Belong",
			"Crawling",
			"Lose Yourself",
			"Stan",
			"The Real Slim Shady",
			"Juicy",
			"Big Poppa",
			"Shook Ones, Pt. II",
			"NY State of Mind",
			"California Love",
			"Changes",
			"Dear Mama",
			"Ms. Jackson",
			"Hey Ya!",
			"Electric Relaxation",
			"Can I Kick It?",
			"No Diggity",
			"Poison",
			"Regulate",
			"Gin and Juice",
			"Still D.R.E.",
			"Forgot About Dre",
			"X Gon' Give It to Ya",
			"Ruff Ryders' Anthem",
			"Hard Knock Life (Ghetto Anthem)",
			"Empire State of Mind",
			"Gold Digger",
			"Stronger",
			"Jesus Walks",
			"Heartless",
			"Come As You Are",
			"Smells Like Teen Spirit",
			"Creep",
			"Karma Police",
			"Mr. Brightside",
			"Seven Nation Army",
			"Boulevard of Broken Dreams",
			"American Idiot",
			"Toxic",
		}
		artists := []string{
			"Linkin Park",
			"Linkin Park",
			"Linkin Park",
			"Men I Trust",
			"Fergie",
			"Clairo",
			"Déyyess",
			"Wu-Tang Clan",
			"Linkin Park",
			"Linkin Park",
			"Linkin Park",
			"Linkin Park",
			"Eminem",
			"Eminem",
			"Eminem",
			"The Notorious B.I.G.",
			"The Notorious B.I.G.",
			"Mobb Deep",
			"Nas",
			"2Pac",
			"2Pac",
			"2Pac",
			"OutKast",
			"OutKast",
			"A Tribe Called Quest",
			"A Tribe Called Quest",
			"Blackstreet",
			"Bell Biv DeVoe",
			"Warren G & Nate Dogg",
			"Snoop Dogg",
			"Dr. Dre",
			"Dr. Dre",
			"DMX",
			"DMX",
			"Jay-Z",
			"Jay-Z & Alicia Keys",
			"Kanye West",
			"Kanye West",
			"Kanye West",
			"Kanye West",
			"Nirvana",
			"Nirvana",
			"Radiohead",
			"Radiohead",
			"The Killers",
			"The White Stripes",
			"Green Day",
			"Green Day",
			"Britney Spears",
		}

		results, err := repo.Recording.GetClosestRecordings(t.Context(), names, artists)
		if err != nil {
			t.Fatalf("GetClosestRecordings failed: %v", err)

		}

		if len(results) != len(names) || len(results) != len(artists) {
			t.Log("Results:", results)
			for _, name := range names {
				found := false

				for _, r := range results {
					if r.RecordingName == name {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("recording %q not found in results", name)
				}
			}
			t.Fatalf("expected %d results, got %d", len(names), len(results))
		}

		for i, r := range results {
			if r.RecordingMBID == "" {
				t.Errorf("result %d has empty MBID", i)
				t.Log("Expected names and artists:", names[i], artists[i])
			}

			if r.RecordingName == "" {
				t.Errorf("result %d has empty recording name", i)
			}

			if r.ArtistName == "" {
				t.Errorf("result %d has empty artist name", i)
			}

		}
		for _, r := range results {
			t.Logf("Recording: %s by %s (MBID: %s)\n", r.RecordingName, r.ArtistName, r.RecordingMBID)
		}
	})
}
