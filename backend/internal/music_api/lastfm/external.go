package musicapi

import (
	"backend-lastfm/internal/utility"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ApiKey string
type LastFMURL string

type PERIOD string

const YEAR PERIOD = "12month"
const MONTHS_6 PERIOD = "6month"
const MONTHS_3 PERIOD = "3month"
const MONTH PERIOD = "1month"
const WEEK PERIOD = "7day"

type LastFMAPIExternal struct {
	ApiKey
	LastFMURL
	setJson bool
}

func NewLastFMExternalAdapter(apiKey ApiKey, lastFMURL LastFMURL, setJson bool) *LastFMAPIExternal {
	return &LastFMAPIExternal{
		apiKey, lastFMURL, setJson,
	}
}

func (h *LastFMAPIExternal) GetUserInfo(userName string) (userInfo UserInfo, err error) {

	extraURLParams := map[string]string{
		"method": "user.getinfo",
	}
	if userName != "" {
		extraURLParams["user"] = userName
	}
	resp, err := h.MakeRequest(extraURLParams)
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

func (h *LastFMAPIExternal) GetUserWeeklyArtists(userName string, from time.Time, to time.Time) (weeklyArtists UserWeeklyTrackList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklyartistchart",
		"user":   userName,
		"from":   strconv.FormatInt(from.Unix(), 10),
		"to":     strconv.FormatInt(to.Unix(), 10),
	}

	resp, err := h.MakeRequest(extraURLParams)
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

func (h *LastFMAPIExternal) GetUserWeeklyAlbums(userName string, from time.Time, to time.Time) (weeklyAlbums UserWeeklyAlbumList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklyalbumchart",
		"user":   userName,
		"from":   strconv.FormatInt(from.Unix(), 10),
		"to":     strconv.FormatInt(to.Unix(), 10),
	}

	resp, err := h.MakeRequest(extraURLParams)
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

func (h *LastFMAPIExternal) GetUserWeeklyTracks(userName string, from time.Time, to time.Time) (weeklyTracks UserWeeklyTrackList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklytrackchart",
		"user":   userName,
		"from":   strconv.FormatInt(from.Unix(), 10),
		"to":     strconv.FormatInt(to.Unix(), 10),
	}

	resp, err := h.MakeRequest(extraURLParams)
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

func (h *LastFMAPIExternal) GetUserTopArtists(userName string, period PERIOD, page int, limit int) (weeklyTracks UserWeeklyTrackList, err error) {

	extraURLParams := map[string]string{
		"method": "user.getweeklytrackchart",
		"user":   userName,
		"period": string(period),
		"page":   strconv.Itoa(page),
		"limit":  strconv.Itoa(limit),
	}

	resp, err := h.MakeRequest(extraURLParams)
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

func (h *LastFMAPIExternal) MakeRequest(extraURLParams map[string]string) (*http.Response, error) {
	q := url.Values{}
	for paramName, paramValue := range extraURLParams {
		q.Set(paramName, paramValue)
	}
	q.Set("api_key", string(h.ApiKey))

	if h.setJson {
		q.Set("format", "json") //JSON RESPONSE
	}

	resp, err := http.Post(
		string(h.LastFMURL),
		"application/x-www-form-urlencoded",
		strings.NewReader(q.Encode()),
	)

	if err != nil {
		return nil, fmt.Errorf("makeRequest Post Error: %w", err)
	}
	return resp, err

}
