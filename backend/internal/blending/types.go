package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"time"
)

type blendId string
type blendLinkValue string

type userid string

type blendCategory string
type blendTimeDuration string

type responseStruct struct {
	Value string `json:"value"`
}

type deleteStruct struct {
	BlendId blendId `json:"blendId"`
}

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

type platformid string

type Blend struct {
	BlendId   string       `json:"blendid"`
	Value     int          `json:"value"`
	Users     []platformid `json:"user"`
	CreatedAt time.Time    `json:"timestamp"`
}

type Blends struct {
	Blends []Blend `json:"blends"`
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
	Users           []string   `json:"Usernames"`
	OverallBlendNum int        `json:"OverallBlendNum"`
	ArtistBlend     TypeBlend  `json:"ArtistBlend"`
	AlbumBlend      TypeBlend  `json:"AlbumBlend"`
	TrackBlend      TypeBlend  `json:"TrackBlend"`
	TopArtists      []TopEntry `json:"topartists"`
	TopAlbums       []TopEntry `json:"topalbums"`
	TopTracks       []TopEntry `json:"toptracks"`
}

type TypeBlend struct {
	ThreeMonth int `json:"ThreeMonth"`
	OneMonth   int `json:"OneMonth"`
	OneYear    int `json:"OneYear"`
}

type TopEntry struct {
	Name         string `json:"name"`
	Artist       string `json:"artist,omitempty"` //Not needed when it's an artist
	Distribution []int  `json:"distribution"`
}

type CatalogueStats = musicapi.CatalogueStats

type complexResponse struct {
	user     userid
	data     map[string]CatalogueStats // Album/Track/Artist -> stats. this used to be map [string]int
	duration blendTimeDuration
	category blendCategory
	err      error
}

// type MapCatStats map[string]CatalogueStats
