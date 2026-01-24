package musicbrainz

import (
	"context"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type RecordingStore struct {
	db *sqlx.DB
}

func (s *RecordingStore) GetReleasesByArtist(artistID string) (any, error) {
	return nil, nil
}

type RecordingCandidate struct {
	RecordingID     int     `db:"recording_id"`
	RecordingMBID   string  `db:"recording_mbid"`
	RecordingName   string  `db:"recording_name"`
	ArtistName      string  `db:"artist_name"`
	ArtistMBID      string  `db:"artist_mbid"`
	TitleSimilarity float64 `db:"title_similarity"`
}

func (s *RecordingStore) GetClosestRecording(context context.Context, name string, artistName string) (RecordingCandidate, error) {
	if artistName == "" {
		artistName = "%"
	}
	candidate := RecordingCandidate{}
	query := `
	SELECT DISTINCT
        r.id AS recording_id,
        r.gid AS recording_mbid,
		r.name AS recording_name,
		a.name AS artist_name,
		a.gid AS artist_mbid,
        similarity(r.name, $1) AS title_similarity
		FROM recording r
		JOIN artist_credit ac ON ac.id = r.artist_credit
		JOIN artist_credit_name acn ON acn.artist_credit = ac.id
		JOIN artist a ON a.id = acn.artist
		WHERE
			r.name % $1
			AND a.name % $2
			AND (r.comment IS NULL OR r.comment = '')
			AND r.name !~* '(live|remix|edit|version|demo|remaster|radio)'
		ORDER BY title_similarity DESC
		LIMIT 1 
	)
	SELECT *
	FROM candidates;

	`

	err := s.db.SelectContext(context, &candidate, query, name, artistName)
	if err != nil {
		return RecordingCandidate{}, err
	}

	return candidate, nil
}

func (s *RecordingStore) DoesRecordExistByMBID(context context.Context, mbid string) (bool, error) {

	query := `
	SELECT EXISTS (
		SELECT 1
		FROM recording
		WHERE gid = $1
		);`
	var exists bool
	err := s.db.GetContext(context, &exists, query, mbid)
	if err != nil {
		return false, err
	}
	return exists, nil

}

func (s *RecordingStore) GetRegexForRecording(name string) (any, error) {
	//WIP - to be implemented
	panic("not implemented")
	// return nil, nil
}

func norm(s string) string {
	s = strings.ToLower(s)

	// Remove content within parentheses
	s = regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(s, "")

	// Remove content within brackets
	s = regexp.MustCompile(`[^a-z0-9\s]`).ReplaceAllString(s, "")

	// Replace multiple spaces with a single space
	return strings.Join(strings.Fields(s), " ")
}
