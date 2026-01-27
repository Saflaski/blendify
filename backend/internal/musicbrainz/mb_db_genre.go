package musicbrainz

import (
	"context"
	"fmt"

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
	ID            int    `db:"genre_id"`
	Name          string `db:"genre"`
	TagCount      int    `db:"tag_count"`
}

func (s *GenreStore) GetGenreByRecordings(context context.Context, recordingMBIDs []string) (map[string][]Genre, error) {
	rows := []RecordingGenreRow{}

	query := `
	WITH all_potential_tags AS (
    SELECT
        r.gid::text AS recording_mbid,
        t.id AS genre_id,
        t.name AS genre,
        rt.count AS tag_count,
        'recording' AS source
    FROM musicbrainz.recording r
    JOIN musicbrainz.recording_tag rt ON rt.recording = r.id
    JOIN musicbrainz.tag t ON t.id = rt.tag
    WHERE r.gid = ANY($1::uuid[])
    UNION ALL
    SELECT
        r.gid::text AS recording_mbid,
        t.id AS genre_id,
        t.name AS genre,
        rgt.count AS tag_count,
        'release_group' AS source
    FROM musicbrainz.recording r
    JOIN musicbrainz.track tr ON tr.recording = r.id
    JOIN musicbrainz.medium m ON m.id = tr.medium
    JOIN musicbrainz.release rel ON rel.id = m.release
    JOIN musicbrainz.release_group rg ON rg.id = rel.release_group
    JOIN musicbrainz.release_group_tag rgt ON rgt.release_group = rg.id
    JOIN musicbrainz.tag t ON t.id = rgt.tag
    WHERE r.gid = ANY($1::uuid[])
      AND NOT EXISTS (
          SELECT 1 FROM musicbrainz.recording_tag rt2 
          WHERE rt2.recording = r.id
      )
),
ranked_tags AS (
    SELECT 
        *,
        ROW_NUMBER() OVER (
            PARTITION BY recording_mbid 
            ORDER BY tag_count DESC, genre ASC 
        ) as rn
    FROM all_potential_tags
	WHERE genre NOT IN (
    	'5+ wochen',
    	'offizielle charts'
    	'ph_3_stars',
    	'ph_temp_checken'
    	'offizielle charts',
    	'ph_3_stars',
    	'offizielle charts',
    	'ph_temp_checken'
)
)
SELECT DISTINCT
    recording_mbid, 
    genre_id, 
    genre, 
    tag_count
FROM ranked_tags
WHERE rn <= 10
ORDER BY recording_mbid, tag_count DESC;
    `
	if err := s.db.SelectContext(context, &rows, query, pq.Array(recordingMBIDs)); err != nil {
		return nil, fmt.Errorf(" error during selectcontext in GetGenreByRecordings: %v", err)
	}

	result := make(map[string][]Genre)

	for _, row := range rows {
		result[row.RecordingMBID] = append(result[row.RecordingMBID], Genre{
			ID:       row.ID,
			Name:     row.Name,
			TagCount: row.TagCount,
		})
		fmt.Printf("Debug: RecordingMBID: %s, Genre: %s, TagCount: %d\n", row.RecordingMBID, row.Name, row.TagCount)
	}

	return result, nil
}
