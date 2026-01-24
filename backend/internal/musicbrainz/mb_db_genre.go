package musicbrainz

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type GenreStore struct {
	db *sqlx.DB
}

type Genre struct {
	ID       int    `db:"genre_id"`
	Name     string `db:"genre"`
	TagCount int    `db:"tag_count"`
}

type RecordingGenreRow struct {
	RecordingMBID string `db:"recording_mbid"`
	Genre
}

func (s *GenreStore) GetGenreByRecordings(context context.Context, recordingMBIDs []string) (map[string][]Genre, error) {
	rows := []RecordingGenreRow{}

	query := `
       SELECT
    r.gid            AS recording_mbid,
    t.id             AS genre_id,
    t.name           AS genre,
    rt.count         AS tag_count
	FROM musicbrainz.recording r
	JOIN musicbrainz.recording_tag rt ON rt.recording = r.id
	JOIN musicbrainz.tag t            ON t.id = rt.tag
	WHERE r.gid = ANY($1)
	ORDER BY rt.tag ASC LIMIT 5;
    `
	err := s.db.SelectContext(context, &rows, query, pq.Array(recordingMBIDs))
	if err != nil {
		return nil, err
	}

	result := make(map[string][]Genre)

	for _, row := range rows {
		result[row.RecordingMBID] = append(result[row.RecordingMBID], row.Genre)
	}

	return result, nil
}
