package musicapi

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
	MBID       string         `json:"mbid"`
	URL        string         `json:"url"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"@attr"`
	Playcount  string         `json:"playcount"`
}

// user.getweeklyalbumchart
type AlbumArtist struct {
	MBID string `json:"mbid"`
	Name string `json:"#text"`
}

type Album struct {
	Artist     AlbumArtist    `json:"artist"`
	MBID       string         `json:"mbid"`
	URL        string         `json:"url"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"@attr"`
	Playcount  string         `json:"playcount"`
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
	Artist     TrackArtist    `json:"artist"`
	MBID       string         `json:"mbid"`
	URL        string         `json:"url"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"@attr"`
	Playcount  string         `json:"playcount"`
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

// user.getttopartists
type UserTopArtists struct {
	TopArtists struct {
		Artist     []Topartist_artist `json:"artist"`
		Attributes map[string]any     `json:"@attr"`
	} `json:"topartists"`
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
}
