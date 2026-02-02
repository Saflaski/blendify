package musicbrainz

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	Genre interface {
		GetGenreByRecordings(context.Context, []string) (map[string][]Genre, error)
		GetGenreByArtistMBIDs(context.Context, []string) (map[string][]Genre, error)
	}
	Recording interface {
		GetReleasesByArtist(string) (any, error)
		GetClosestRecording(context.Context, string, string) (RecordingCandidate, error)
		GetClosestRecordings(context.Context, []string, []string) ([]RecordingCandidate, error)
		GetRegexForRecording(string) (any, error)
		DoesRecordExistByMBID(context.Context, string) (bool, error)
	}
	Artist interface {
		GetClosestArtistsByName(context.Context, []string) ([]ArtistCandidate, error)
	}
}

func NewPostgresMusicBrainzRepo(db *sqlx.DB) Storage {
	return Storage{
		Genre:     &GenreStore{db: db},
		Recording: &RecordingStore{db: db},
		Artist:    &ArtistStore{db: db},
	}
}

func NewOnlineMusicBrainzRepo() Storage {
	return Storage{}
}
