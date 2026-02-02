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
	Source   string `db:"source"`
	Rank     int    `db:"rn"`
}

type RecordingGenreRow struct {
	RecordingMBID string `db:"recording_mbid"`
	ID            int    `db:"genre_id"`
	Name          string `db:"genre"`
	TagCount      int    `db:"tag_count"`
	Source        string `db:"source"`
	Rank          int    `db:"rn"`
}

type ArtistGenreRow struct {
	ArtistMBID string `db:"artist_mbid"`
	GenreName  string `db:"genre"`
	TagCount   int    `db:"tag_count"`
	Rank       int    `db:"rn"`
}

func (s *GenreStore) GetGenreByArtistMBIDs(context context.Context, artistMBIDs []string) (map[string][]Genre, error) {

	rows := []ArtistGenreRow{}
	query := `
	with all_potential_tags as (select 
	a.gid::text AS artist_mbid,
	t."name" as genre,
	t.ref_count as tag_count
from artist a
join artist_tag at
on at.artist = a.id 
join tag t 
on at.tag = t.id 
where a."gid" = any($1::uuid[])
),
ranked_tags as (
	select *, row_number() over (
		partition by artist_mbid
		order by tag_count desc, genre asc
	) as rn
	from all_potential_tags
)
select artist_mbid, genre, tag_count, rn from ranked_tags
where rn <= 10
order by artist_mbid, tag_count DESC`

	if err := s.db.SelectContext(context, &rows, query, pq.Array(artistMBIDs)); err != nil {
		return nil, fmt.Errorf(" error during selectcontext in GetGenreByRecordings: %v", err)
	}

	result := make(map[string][]Genre)

	for _, row := range rows {
		result[row.ArtistMBID] = append(result[row.ArtistMBID], Genre{
			Name:     row.GenreName,
			TagCount: row.TagCount,
			Rank:     row.Rank,
		})
		// fmt.Printf("Debug: Artist: %s, Genre: %s, TagCount: %d\n", row.ArtistMBID, row.Name, row.TagCount)
	}

	return result, nil

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
			Rank:     row.Rank,
			Source:   row.Source,
		})
		// fmt.Printf("Debug: RecordingMBID: %s, Genre: %s, TagCount: %d\n", row.RecordingMBID, row.Name, row.TagCount)
	}

	return result, nil
}
