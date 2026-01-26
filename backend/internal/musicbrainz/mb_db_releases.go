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
			FROM UNNEST($1::text[], $2::text[])
		)
		SELECT
			i.title,
			i.artist,
			c.recording_mbid
		FROM input i
		JOIN LATERAL (
			SELECT r.gid AS recording_mbid
			FROM recording r
			JOIN artist_credit ac ON ac.id = r.artist_credit
			JOIN artist_credit_name acn ON acn.artist_credit = ac.id
			JOIN artist a ON a.id = acn.artist
			LEFT JOIN LATERAL (
				SELECT SUM(rt.count) AS tag_count
				FROM recording_tag rt
				WHERE rt.recording = r.id
			) t ON true
			LEFT JOIN release rel ON rel.id = (
				SELECT t2.medium
				FROM track t2
				WHERE t2.recording = r.id
				LIMIT 1
			)
			LEFT JOIN release_group rg ON rg.id = rel.release_group
			LEFT JOIN release_group_primary_type rgpt ON rg.type = rgpt.id
			WHERE r.name ILIKE i.title || '%'
			AND r.name % i.title
			AND a.name % i.artist
			AND (r.comment IS NULL OR r.comment = '')
			AND r.name !~* '(live|remix|edit|version|demo|remaster|radio)'
			ORDER BY
				COALESCE(t.tag_count, 0) DESC,
				similarity(r.name, i.title) DESC,
				CASE rgpt.name
					WHEN 'Album' THEN 1
					WHEN 'Single' THEN 2
					WHEN 'EP' THEN 3
					ELSE 4
				END,
				r.length IS NOT NULL DESC,
				r.id
			LIMIT 1
		) c ON TRUE;
	`

	queryTwo := `
	WITH input_pairs AS (
  SELECT
    row_number() OVER () AS pair_id,
    t.name,
    t.artist
  FROM unnest(
    $1::text[],
    $2::text[]
  ) AS t(name, artist)
),
exact_matches AS (
  SELECT
    ip.pair_id,
    r.gid AS mbid,
    r.name AS rname,
    ac.name AS artist,
    COUNT(rt.tag) AS occurrence_count
  FROM input_pairs ip
  JOIN recording r
    ON r.name = ip.name
	AND r.comment = ''
  JOIN artist_credit ac
    ON r.artist_credit = ac.id AND ac.name = ip.artist
  LEFT JOIN recording_tag rt
    ON r.id = rt.recording
  GROUP BY ip.pair_id, r.gid, r.name, ac.name
), 
fuzzy_matches AS (
  SELECT
    ip.pair_id,
    r.gid AS mbid,
    r.name AS rname,
    ac.name AS artist,
    COUNT(rt.tag) AS occurrence_count
  FROM input_pairs ip
  LEFT JOIN exact_matches em
    ON ip.pair_id = em.pair_id
  JOIN recording r
    ON r.name % ip.name
	AND r.comment = ''
  JOIN artist_credit ac
    ON r.artist_credit = ac.id AND ac.name % ip.artist
  LEFT JOIN recording_tag rt
    ON r.id = rt.recording
  WHERE em.pair_id IS NULL
  GROUP BY ip.pair_id, r.gid, r.name, ac.name
),
final as ( select * from (
  SELECT *, ROW_NUMBER() OVER (PARTITION BY pair_id ORDER BY occurrence_count DESC) AS rn
  FROM (
    SELECT * FROM exact_matches
    UNION ALL
    SELECT * FROM fuzzy_matches
  ) combined
) t
WHERE rn = 1
ORDER BY pair_id
)
select mbid, rname, artist from final;

	`
	_ = query
	_ = queryTwo
	for i, name := range names {
		inputs[i] = Input{
			Title:  name,
			Artist: artistNames[i],
		}
	}

	rows, err := s.db.QueryxContext(context, queryTwo, (names), (artistNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var title, artist, mbid string
		rows.Scan(&mbid, &title, &artist)
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
