package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
)

type blendId string
type blendLinkValue string

type userid string

type blendCategory string
type blendTimeDuration string

var durationMap = map[blendTimeDuration]musicapi.Period{
	BlendTimeDurationOneMonth:   musicapi.ONE_MONTH,
	BlendTimeDurationThreeMonth: musicapi.THREE_MONTHS,
	BlendTimeDurationYear:       musicapi.YEAR,
}

const (
	BlendCategoryArtist blendCategory = "artist"
	BlendCategoryTrack  blendCategory = "track"
	BlendCategoryAlbum  blendCategory = "album"
)

const (
	BlendTimeDurationOneMonth   blendTimeDuration = "1month"
	BlendTimeDurationThreeMonth blendTimeDuration = "3month"
	// BlendTimeDurationSixMonth   blendTimeDuration = "6month"
	BlendTimeDurationYear blendTimeDuration = "12month"
	// BlendTimeDurationAllTime    blendTimeDuration = "alltime"
)

var durationRange = []blendTimeDuration{
	BlendTimeDurationOneMonth,
	BlendTimeDurationThreeMonth,
	BlendTimeDurationYear,
}

var categoryRange = []blendCategory{
	BlendCategoryAlbum,
	BlendCategoryArtist,
	BlendCategoryTrack,
}

type Blend struct {
	id    string
	users []userid
}

type BlendResponse struct {
	ID               string             `json:"id"`
	Users            []userid           `json:"Users"`
	BlendPercentages []IndividualBlends `json:"blendpercentages"`
}

type IndividualBlends struct {
	Type string `json:"type"`
}

type IndividualUserData struct {
	Type string `json:"type"`
}

type DuoBlend struct {
	Users           []string   `json:"usernames"`
	OverallBlendNum int        `json:"overall"`
	ArtistBlend     TypeBlend  `json:"artist"`
	AlbumBlend      TypeBlend  `json:"album"`
	TrackBlend      TypeBlend  `json:"track"`
	TopArtists      []TopEntry `json:"topartists"`
	TopAlbums       []TopEntry `json:"topalbums"`
	TopTracks       []TopEntry `json:"toptracks"`
}

type TypeBlend struct {
	ThreeMonth int `json:"3month"`
	OneMonth   int `json:"1month"`
	OneYear    int `json:"1year"`
}

type TopEntry struct {
	Name         string `json:"name"`
	Distribution []int  `json:"distribution"`
}
