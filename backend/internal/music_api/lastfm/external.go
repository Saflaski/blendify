package musicapi

import (
	"backend-lastfm/internal/utility"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/golang/glog"
	"golang.org/x/time/rate"
)

type LastFMAPIExternal struct {
	apiKey               string
	lastFMURL            string
	setJson              bool
	rateLimitMinimumTime int
}

type Period string

var requestLimiter = rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

const (
	YEAR         Period = "12month"
	SIX_MONTHS   Period = "6month"
	THREE_MONTHS Period = "3month"
	ONE_MONTH    Period = "1month"
	WEEK         Period = "7day"
)

type websessionKey string

func NewLastFMExternalAdapter(apiKey, lastFMURL string, setJson bool, rateLimitMinimumTime int) *LastFMAPIExternal {

	return &LastFMAPIExternal{
		apiKey, lastFMURL, setJson, rateLimitMinimumTime,
	}
}

func (h *LastFMAPIExternal) GetAPISignature(wsKey websessionKey, methodName string) {

}

func (h *LastFMAPIExternal) GetUserInfo(ctx context.Context, userName string) (userInfo UserInfo, err error) {

	extraURLParams := map[string]string{
		"method": "user.getinfo",
	}
	if userName != "" {
		extraURLParams["user"] = userName
	}

	resp, err := h.MakeRequest(ctx, extraURLParams)
	if err != nil {
		return userInfo, fmt.Errorf("GetUserInfo makeRequest Error: %v", err)

	}
	defer resp.Body.Close()

	userInfo, err = utility.Decode[UserInfo](resp)
	if err != nil {
		return userInfo, fmt.Errorf("GetUserInfo decode Error: %v", err)
	}

	return userInfo, nil
}

func (h *LastFMAPIExternal) GetUserWeeklyArtists(ctx context.Context, userName string, from time.Time, to time.Time) (weeklyArtists UserWeeklyTrackList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklyartistchart",
		"user":   userName,
		"from":   strconv.FormatInt(from.Unix(), 10),
		"to":     strconv.FormatInt(to.Unix(), 10),
	}

	resp, err := h.MakeRequest(ctx, extraURLParams)
	if err != nil {
		return UserWeeklyTrackList{}, fmt.Errorf("GetUserWeeklyArtists makeRequest Error: %v", err)
	}
	defer resp.Body.Close()

	weeklyArtists, err = utility.Decode[UserWeeklyTrackList](resp)
	if err != nil {
		return weeklyArtists, fmt.Errorf("GetUserWeeklyArtists decode Error: %v", err)
	}

	return weeklyArtists, nil
}

func (h *LastFMAPIExternal) GetUserWeeklyAlbums(context context.Context, userName string, from time.Time, to time.Time) (weeklyAlbums UserWeeklyAlbumList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklyalbumchart",
		"user":   userName,
		"from":   strconv.FormatInt(from.Unix(), 10),
		"to":     strconv.FormatInt(to.Unix(), 10),
	}

	resp, err := h.MakeRequest(context, extraURLParams)
	if err != nil {
		return UserWeeklyAlbumList{}, fmt.Errorf("GetUserWeeklyArtists makeRequest Error: %v", err)
	}
	defer resp.Body.Close()

	weeklyAlbums, err = utility.Decode[UserWeeklyAlbumList](resp)
	if err != nil {
		return weeklyAlbums, fmt.Errorf("GetUserWeeklyArtists decode Error: %v", err)
	}

	return weeklyAlbums, nil
}

func (h *LastFMAPIExternal) GetUserWeeklyTracks(context context.Context, userName string, from time.Time, to time.Time) (weeklyTracks UserWeeklyTrackList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklytrackchart",
		"user":   userName,
		"from":   strconv.FormatInt(from.Unix(), 10),
		"to":     strconv.FormatInt(to.Unix(), 10),
	}

	resp, err := h.MakeRequest(context, extraURLParams)
	if err != nil {
		return UserWeeklyTrackList{}, fmt.Errorf("GetUserWeeklyArtists makeRequest Error: %v", err)
	}
	defer resp.Body.Close()

	weeklyTracks, err = utility.Decode[UserWeeklyTrackList](resp)
	if err != nil {
		return weeklyTracks, fmt.Errorf("GetUserWeeklyArtists decode Error: %v", err)
	}

	return weeklyTracks, nil
}

func (h *LastFMAPIExternal) GetUserTopArtists(context context.Context, userName string, period string, maxPages int, limit int) (UserTopArtists, error) {

	if maxPages == 0 {
		maxPages = 1
	}
	if limit == 0 {
		limit = 50
	}
	topArtists := make([]UserTopArtists, 0, maxPages)
	completeArtists := make([]Artist, 0, maxPages*limit)

	for page := 1; page <= maxPages; page++ {
		//Default LFM API limit per page
		extraURLParams := map[string]string{
			"method": "user.gettopartists",
			"user":   userName,
			"period": string(period),
			"page":   strconv.Itoa(page),
			"limit":  strconv.Itoa(limit),
		}

		resp, err := h.MakeRequest(context, extraURLParams)
		if err != nil {
			return UserTopArtists{}, fmt.Errorf("UserTopArtists makeRequest Error: %v", err)
		}

		nextTopArtists, err := utility.Decode[UserTopArtists](resp)
		if err != nil {
			return UserTopArtists{}, fmt.Errorf("UserTopArtists decode Error: %v", err)
		}
		topArtists = append(topArtists, nextTopArtists)
		artists := nextTopArtists.TopArtists.Artist
		if len(artists) == 0 {
			break
		}
		completeArtists = append(completeArtists, artists...)

		resp.Body.Close()

	}
	return UserTopArtists{
		TopArtists: TopArtists{
			Artist: completeArtists,
		},
	}, nil
}

func (h *LastFMAPIExternal) GetUserTopAlbums(context context.Context, userName string, period string, maxPages int, limit int) (UserTopAlbums, error) {

	if maxPages == 0 {
		maxPages = 1
	}
	if limit == 0 {
		limit = 50
	}
	topAlbums := make([]UserTopAlbums, 0, maxPages)
	completeAlbums := make([]Album, 0, maxPages*limit)

	for page := 1; page <= maxPages; page++ {
		//Default LFM API limit per page
		extraURLParams := map[string]string{
			"method": "user.gettopalbums",
			"user":   userName,
			"period": string(period),
			"page":   strconv.Itoa(page),
			"limit":  strconv.Itoa(limit),
		}

		resp, err := h.MakeRequest(context, extraURLParams)
		if err != nil {
			return UserTopAlbums{}, fmt.Errorf("UserTopAlbums makeRequest Error: %v", err)
		}

		nextTopAlbums, err := utility.Decode[UserTopAlbums](resp)
		if err != nil {
			return UserTopAlbums{}, fmt.Errorf("UserTopAlbums decode Error: %v", err)
		}
		topAlbums = append(topAlbums, nextTopAlbums)
		albums := nextTopAlbums.TopAlbums.Album
		if len(albums) == 0 {
			break
		}
		completeAlbums = append(completeAlbums, albums...)

		resp.Body.Close()

	}
	return UserTopAlbums{
		TopAlbums: TopAlbums{
			Album: completeAlbums,
		},
	}, nil
}

func (h *LastFMAPIExternal) GetUserTopTracks(context context.Context, userName string, period string, maxPages int, limit int) (UserTopTracks, error) {
	if maxPages == 0 {
		maxPages = 1
	}
	if limit == 0 {
		limit = 50
	}
	topTracks := make([]UserTopTracks, 0, maxPages)
	completeTracks := make([]Track, 0, maxPages*limit)

	for page := 1; page <= maxPages; page++ {
		//Default LFM API limit per page
		extraURLParams := map[string]string{
			"method": "user.gettoptracks",
			"user":   userName,
			"period": string(period),
			"page":   strconv.Itoa(page),
			"limit":  strconv.Itoa(limit),
		}

		resp, err := h.MakeRequest(context, extraURLParams)
		if err != nil {
			return UserTopTracks{}, fmt.Errorf("UserTopTracks makeRequest Error: %v", err)
		}

		nextTopTracks, err := utility.Decode[UserTopTracks](resp)
		if err != nil {
			return UserTopTracks{}, fmt.Errorf("UserTopTracks decode Error: %v", err)
		}
		topTracks = append(topTracks, nextTopTracks)
		tracks := nextTopTracks.TopTracks.Track
		if len(tracks) == 0 {
			break
		}
		completeTracks = append(completeTracks, tracks...)

		resp.Body.Close()

		// if len(topTracks[i].TopTracks.Track) == 0 {
		// 	break
		// }

	}

	// completeUserTopTracks := make([]Track, 0)
	// for _, val := range topTracks {
	// 	completeUserTopTracks = append(completeUserTopTracks, val.TopTracks.Track...)
	// }
	// FinalUserTopTracks := UserTopTracks{
	// 	TopTracks: TopTracks{
	// 		Track: completeUserTopTracks, //dc: why is this empty?
	// 	},
	// }

	// return FinalUserTopTracks, nil

	return UserTopTracks{
		TopTracks: TopTracks{
			Track: completeTracks,
		},
	}, nil
}

func (h *LastFMAPIExternal) MakeRequest(ctx context.Context, extraURLParams map[string]string) (*http.Response, error) {
	// q := url.Values{}
	// for paramName, paramValue := range extraURLParams {
	// 	q.Set(paramName, paramValue)
	// }
	// q.Set("api_key", string(h.apiKey))

	requestLimiter.Wait(context.Background()) //For delaying every request made by the backend at minimum 200ms from each other.
	glog.Infof("Making LFM API Request: %s", extraURLParams)
	u, err := url.Parse(h.lastFMURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	q := u.Query()
	for paramName, paramValue := range extraURLParams {
		q.Set(paramName, paramValue)
	}
	q.Set("api_key", string(h.apiKey))

	if h.setJson {
		q.Set("format", "json") //JSON RESPONSE
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(" during makeRequest, could not make new request with Error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(" during makeRequest, GET Error: %w", err)
	}

	// resp, err := http.Get(u.String())
	// if err != nil {
	// 	return nil, fmt.Errorf(" during makeRequest, GET Error: %w", err)
	// }

	// resp, err := http.Post(
	// 	string(h.lastFMURL),
	// 	"application/x-www-form-urlencoded",
	// 	strings.NewReader(q.Encode()),
	// )

	// if err != nil {
	// 	return nil, fmt.Errorf("makeRequest Post Error: %w", err)
	// }

	return resp, err

}
