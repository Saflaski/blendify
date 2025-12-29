package musicapi

import "encoding/json"

// user.getinfo
type UserInfo struct {
	User struct {
		Name            string `json:"name"`
		Age             string `json:"age"`
		Subscribers     string `json:"subscriber"`
		RealName        string `json:"realname"`
		Bootstraped     string `json:"bootstrap"`
		Playcount       string `json:"playcount"`
		Artist_count    string `json:"artist_count"`
		Track_count     string `json:"track_count"`
		Album_count     string `json:"album_count"`
		Registered_unix string `json:"registered"`
		URL             string `json:"url"`
	} `json:"user"`
}

// user.getweeklyartistchart
type UserWeeklyArtistList struct {
	WeeklyArtistChart struct {
		Artist []Artist `json:"artist"`
	} `json:"weeklyartistchart"`
}

type Artist struct {
	MBID       string         `json:"mbid,omitempty"`
	URL        string         `json:"url,omitempty"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"@attr ,omitempty"`
	Playcount  string         `json:"playcount,omitempty"`
	LFMImages  []LFMImage     `json:"image,omitempty"`
}

type LFMImage struct {
	Size string `json:"size"`
	URL  string `json:"#text"`
}

// user.getweeklyalbumchart
type AlbumArtist struct {
	MBID string `json:"mbid"`
	Name string `json:"#text"`
}

type Album struct {
	Artist     Artist         `json:"artist"`
	MBID       string         `json:"mbid"`
	URL        string         `json:"url"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"@attr"`
	Playcount  string         `json:"playcount"`
	LFMImages  []LFMImage     `json:"image"`
}

type UserWeeklyAlbumList struct {
	WeeklyAlbumChart struct {
		Album []Album `json:"album"`
	} `json:"weeklyalbumchart"`
}

// user.getweeklytrackchart
type TrackArtist struct {
	MBID string `json:"mbid"`
	Name string `json:"#text"`
}

type Track struct {
	Artist     Artist         `json:"artist"`
	MBID       string         `json:"mbid"`
	URL        string         `json:"url"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"@attr"`
	Playcount  string         `json:"playcount"`
	LFMImages  []LFMImage     `json:"image"`
}

type UserWeeklyTrackList struct {
	WeeklyTrackChart struct {
		Track []Track `json:"track"`
	} `json:"weeklytrackchart"`
}

type ErrorResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

// type TopArtists struct {
// 	Artist     []Topartist_artist `json:"artist"`
// 	Attributes map[string]any     `json:"@attr"`
// }

type Topartist_artist struct {
	MBID       string         `json:"mbid"`
	URL        string         `json:"url"`
	Playcount  string         `json:"playcount"`
	Attributes map[string]any `json:"attr"`
	Name       string         `json:"name"`
	LFMImages  []LFMImage     `json:"image"`
}

// type LFMImages []LFMImage

// user.gettopalbums
type UserTopAlbums struct {
	TopAlbums  TopAlbums      `json:"topalbums"`
	Attributes map[string]any `json:"@attr"`
}

// user.gettoptracks
type TopTracks struct {
	Track []Track `json:"track"`
}

type TopArtists struct {
	Artist     []Artist       `json:"artist"`
	Attributes map[string]any `json:"@attr"`
}

type TopAlbums struct {
	Album      []Album        `json:"album"`
	Attributes map[string]any `json:"@attr"`
}

type UserTopTracks struct {
	TopTracks  TopTracks      `json:"toptracks"`
	Attributes map[string]any `json:"@attr"`
}

// user.getttopartists
type UserTopArtists struct {
	TopArtists TopArtists `json:"topartists"`
}

type CatalogueStats struct { //A catalogue can be an album, track or artist. The following is metadata for a catalogue
	Artist      Artist `json:"artist"`
	Count       int    `json:"count"`
	PlatformURL string `json:"platformurl"` //Catalogue URL
	Image       string `json:"imageurl"`    //Image URL
	PlatformID  string `json:"platformid"`  //Catalogue Platform ID
}

func JSONToMapCatStats(data []byte) (map[string]CatalogueStats, error) {
	var out map[string]CatalogueStats
	err := json.Unmarshal(data, &out)
	return out, err
}
