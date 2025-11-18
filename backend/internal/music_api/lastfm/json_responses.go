package musicapi

// user.getinfo
type UserInfo struct {
	User struct {
		Name            string `json:"name"`
		Age             int    `json:"age"`
		Subscribers     int    `json:"subscriber"`
		RealName        string `json:"realname"`
		Bootstraped     string `json:"bootstrap"`
		Playcount       int    `json:"playcount"`
		Artist_count    int    `json:"artist_count"`
		Track_count     int    `json:"track_count"`
		Album_count     int    `json:"album_count"`
		Registered_unix int    `json:"registered"`
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
	MBID       string                 `json:"mbid"`
	URL        string                 `json:"url"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"@attr"`
	Playcount  int                    `json:"playcount"`
}

// user.getweeklyalbumchart
type AlbumArtist struct {
	MBID string `json:"mbid"`
	Name string `json:"#text"`
}

type Album struct {
	Artist     AlbumArtist            `json:"artist"`
	MBID       string                 `json:"mbid"`
	URL        string                 `json:"url"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"@attr"`
	Playcount  int                    `json:"playcount"`
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
	Artist     TrackArtist            `json:"artist"`
	MBID       string                 `json:"mbid"`
	URL        string                 `json:"url"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"@attr"`
	Playcount  int                    `json:"playcount"`
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
