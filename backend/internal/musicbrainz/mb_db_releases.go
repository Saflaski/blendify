package musicbrainz

import (
	"context"
	"fmt"
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

type Input struct {
	Title  string `db:"title"`
	Artist string `db:"artist"`
}

func (s *RecordingStore) GetClosestRecordings(context context.Context, names []string, artistNames []string) ([]RecordingCandidate, error) {

	if len(names) != len(artistNames) {
		return nil, fmt.Errorf("names and artistNames slices must have the same length")
	}

	inputs := make([]Input, len(names))
	candidates := make([]RecordingCandidate, len(names))
	query := `
			WITH input(title, artist) AS (
		SELECT *
		FROM UNNEST(
			$1::text[],
			$2::text[]
		)
	)
	SELECT
		i.title,
		i.artist,
		c.recording_mbid
	FROM input i
	JOIN LATERAL (
		SELECT recording_mbid
		FROM (
			SELECT DISTINCT
				r.id,
				r.gid AS recording_mbid,
				similarity(r.name, i.title) AS title_similarity
			FROM recording r
			JOIN artist_credit ac ON ac.id = r.artist_credit
			JOIN artist_credit_name acn ON acn.artist_credit = ac.id
			JOIN artist a ON a.id = acn.artist
			WHERE
				r.name ILIKE i.title || '%' AND r.name % i.title 
				AND a.name % i.artist
				AND (r.comment IS NULL OR r.comment = '')
				AND r.name !~* '(live|remix|edit|version|demo|remaster|radio)'
		) dedup
		ORDER BY title_similarity DESC
		LIMIT 1
	) c ON TRUE;



	`
	for i, name := range names {
		inputs[i] = Input{
			Title:  name,
			Artist: artistNames[i],
		}
	}

	rows, err := s.db.QueryxContext(context, query, (names), (artistNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var title, artist, mbid string
		rows.Scan(&title, &artist, &mbid)
		candidates[i] = RecordingCandidate{
			RecordingName: title,
			ArtistName:    artist,
			RecordingMBID: mbid,
		}
		i++
	}

	return candidates, nil

	// return nil, nil
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

func norm(s []string) []string {

	for i, str := range s {
		str = strings.ToLower(str)
		str = regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(str, "")
		str = regexp.MustCompile(`[^a-z0-9\s]`).ReplaceAllString(str, "")
		s[i] = strings.Join(strings.Fields(str), " ")
	}
	return s
}
