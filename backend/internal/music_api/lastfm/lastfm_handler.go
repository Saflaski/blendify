package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ===================================================================
// Go file containing implementations of all relevant LastFM API calls
// ===================================================================

type ApiKey string
type LastFMURL string

type LastFMHAPIHandler struct {
	ApiKey
	LastFMURL
}

func NewLastFMHAPIHandler(apiKey ApiKey, lastFMURL LastFMURL) *LastFMHAPIHandler {
	return &LastFMHAPIHandler{
		apiKey, lastFMURL,
	}
}

func (h *LastFMHAPIHandler) makeRequest(extraURLParams map[string]string) (*http.Response, error) {
	q := url.Values{}
	for paramName, paramValue := range extraURLParams {
		q.Set(paramName, paramValue)

	}
	q.Set("api_key", string(h.ApiKey))

	q.Set("format", "json") //JSON RESPONSE

	resp, err := http.Post(
		string(h.LastFMURL),
		"application/x-www-form-urlencoded",
		strings.NewReader(q.Encode()),
	)

	if err != nil {
		return nil, fmt.Errorf("makeRequest Post Error: %w", err)
	}

	// defer resp.Body.Close()

	return resp, err

}

// User Calls

// user.getInfo
func (h *LastFMHAPIHandler) GetUserInfo(userName string) {

}
