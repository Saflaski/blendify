package musicbrainz

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type ArtistStore struct {
	db *sqlx.DB
}

type ArtistCandidate struct {
	ArtistName string `db:"artist_name"`
	ArtistMBID string `db:"artist_mbid"`
}

func (a *ArtistStore) GetClosestArtistsByName(context context.Context, artistNames []string) ([]ArtistCandidate, error) {

	candidates := make([]ArtistCandidate, len(artistNames))
	query := `
select gid as mbid, name from artist
where artist."name" = any($1::text[])
order by artist.begin_date_year asc
limit 1
;`
	rows, err := a.db.QueryxContext(context, query, (artistNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var artist, mbid string
		rows.Scan(&mbid, &artist)
		candidates[i] = ArtistCandidate{
			ArtistName: artist,
			ArtistMBID: mbid,
		}
		i++
	}

	return candidates, nil
}
