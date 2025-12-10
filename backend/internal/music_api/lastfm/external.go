package musicapi

import (
	"backend-lastfm/internal/utility"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Period string

const (
	YEAR         Period = "12month"
	SIX_MONTHS   Period = "6month"
	THREE_MONTHS Period = "3month"
	ONE_MONTH    Period = "1month"
	WEEK         Period = "7day"
)

type LastFMAPIExternal struct {
	apiKey    string
	lastFMURL string
	setJson   bool
}

type websessionKey string

func NewLastFMExternalAdapter(apiKey, lastFMURL string, setJson bool) *LastFMAPIExternal {
	return &LastFMAPIExternal{
		apiKey, lastFMURL, setJson,
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

func (h *LastFMAPIExternal) GetUserTopArtists(context context.Context, userName string, period Period, page int, limit int) (topArtists UserTopArtists, err error) {

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 50
	}

	extraURLParams := map[string]string{
		"method": "user.gettopartists",
		"user":   userName,
		"period": string(period),
		"page":   strconv.Itoa(page),
		"limit":  strconv.Itoa(limit),
	}

	resp, err := h.MakeRequest(context, extraURLParams)
	if err != nil {
		return UserTopArtists{}, fmt.Errorf("GetUserTopArtists makeRequest Error: %v", err)
	}
	defer resp.Body.Close()

	topArtists, err = utility.Decode[UserTopArtists](resp)
	if err != nil {
		return topArtists, fmt.Errorf("GetUserTopArtists decode Error: %v", err)
	}

	return topArtists, nil
}

func (h *LastFMAPIExternal) GetUserTopAlbums(context context.Context, userName string, period Period, page int, limit int) (topAlbums UserTopAlbums, err error) {

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 50
	}

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
	defer resp.Body.Close()

	topAlbums, err = utility.Decode[UserTopAlbums](resp)
	if err != nil {
		return topAlbums, fmt.Errorf("UserTopAlbums decode Error: %v", err)
	}

	return topAlbums, nil
}

func (h *LastFMAPIExternal) GetUserTopTracks(context context.Context, userName string, period Period, page int, limit int) (topTracks UserTopTracks, err error) {

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 50
	}

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
	defer resp.Body.Close()

	topTracks, err = utility.Decode[UserTopTracks](resp)
	if err != nil {
		return topTracks, fmt.Errorf("UserTopTracks decode Error: %v", err)
	}

	return topTracks, nil
}

func (h *LastFMAPIExternal) MakeRequest(ctx context.Context, extraURLParams map[string]string) (*http.Response, error) {
	// q := url.Values{}
	// for paramName, paramValue := range extraURLParams {
	// 	q.Set(paramName, paramValue)
	// }
	// q.Set("api_key", string(h.apiKey))

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
