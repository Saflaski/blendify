package musicapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func BenchmarkUserTopTracks(b *testing.B) {
	lfm_adapter := NewLastFMExternalAdapter(
		os.Getenv("LASTFM_API_KEY"),
		"https://ws.audioscrobbler.com/2.0/",
		true,
		10,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := lfm_adapter.GetUserTopTracks(
			b.Context(),
			"saflas",
			"3month",
			1,
			50,
		)
		if err != nil {
			b.Errorf("Expected no error, got %q", err)
		}
	}
}

func BenchmarkAsyncUserTopTracks(b *testing.B) {
	lfm_adapter := NewLastFMExternalAdapter(
		os.Getenv("LASTFM_API_KEY"),
		"https://ws.audioscrobbler.com/2.0/",
		true,
		10,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := lfm_adapter.GetUserTopTracksAsync(
			b.Context(),
			"saflas",
			"3month",
			1,
			50,
		)
		if err != nil {
			b.Errorf("Expected no error, got %q", err)
		}
	}
}

func TestEndpoints(t *testing.T) {

	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal("godotenv.Load failed")
	}
	var apiKey = os.Getenv("LASTFM_API_KEY")
	var lastFMURL = "http://ws.audioscrobbler.com/2.0/"

	apiHandler := NewLastFMExternalAdapter(
		apiKey,
		lastFMURL,
		true,
		200,
	)

	t.Run("Get User Weekly Chart List", func(t *testing.T) {
		extraURLParams := map[string]string{
			"method": "user.getweeklychartlist",
			"user":   "saflas",
			"from":   "1751198400",
			"to":     "1751803200",
		}

		resp, err := checkResponseOK(apiHandler.MakeRequest(t.Context(), extraURLParams))

		if err != nil {
			t.Errorf("Error: %q", err)
		}
		defer resp.Body.Close()

		if err := checkResponseBody(resp); err != nil {
			t.Errorf("Response error: %q", err)
		}

	})

	t.Run("Get User Weekly Artist", func(t *testing.T) {

		extraURLParams := map[string]string{
			"method": "user.getweeklyartistchart",
			"user":   "saflas",
			"from":   "1749988800",
			"to":     "1750593600",
		}

		resp, err := checkResponseOK(apiHandler.MakeRequest(t.Context(), extraURLParams))

		if err != nil {
			t.Errorf("Error: %q", err)
		}
		defer resp.Body.Close()

		if err := checkResponseBody(resp); err != nil {
			t.Errorf("Response error: %q", err)
		}

	})

	t.Run("Get User Weekly Albums", func(t *testing.T) {

		extraURLParams := map[string]string{
			"method": "user.getweeklyalbumchart",
			"user":   "saflas",
			"from":   "1749988800",
			"to":     "1750593600",
		}

		resp, err := checkResponseOK(apiHandler.MakeRequest(t.Context(), extraURLParams))

		if err != nil {
			t.Errorf("Error: %q", err)
		}
		defer resp.Body.Close()

		if err := checkResponseBody(resp); err != nil {
			t.Errorf("Response error: %q", err)
		}

	})

	t.Run("Get User Weekly Tracks", func(t *testing.T) {

		extraURLParams := map[string]string{
			"method": "user.getweeklytrackchart",
			"user":   "saflas",
			"from":   "1749988800",
			"to":     "1750593600",
		}

		resp, err := checkResponseOK(apiHandler.MakeRequest(t.Context(), extraURLParams))

		if err != nil {
			t.Errorf("Error: %q", err)
		}
		defer resp.Body.Close()

		if err := checkResponseBody(resp); err != nil {
			t.Errorf("Response error: %q", err)
		}

	})

	t.Run("Get User Top Artists", func(t *testing.T) {

		extraURLParams := map[string]string{
			"method": "user.gettopartists",
			"user":   "saflas",
			"page":   "1",
			"limit":  "50",
		}

		resp, err := checkResponseOK(apiHandler.MakeRequest(t.Context(), extraURLParams))

		if err != nil {
			t.Errorf("Error: %q", err)
		}
		t.Error()
		defer resp.Body.Close()

		if err := checkResponseBody(resp); err != nil {
			t.Errorf("Response error: %q", err)
		}

	})

	t.Run("Get User - Top 50 Listened to Artists -- 3 months", func(t *testing.T) {
		userName := "test2002"
		// response, err := blendService.getTopArists(userName, BlendTimeDurationThreeMonth)
		TopArtistResponse, err := apiHandler.GetUserTopArtists(
			t.Context(),
			userName,
			"3month",
			1,
			50,
		)
		// glog.Error(topArtist, err)

		if err != nil {
			t.Errorf("Expected no error, got %q", err)
		}

		mostListenedToArtist := TopArtistResponse.TopArtists.Artist[0].Name
		if mostListenedToArtist == "" {
			t.Errorf("Got empty string after processing LastFM response: %q", mostListenedToArtist)

		}
		if len(TopArtistResponse.TopArtists.Artist) == 0 {
			t.Errorf("empty error from processing topartists")
		}
		// t.Error(fmt.Sprint(len(TopArtistResponse.TopArtists.Artist)))

	})

	t.Run("Get UserTopTracks Async", func(t *testing.T) {
		res, err := apiHandler.GetUserTopTracksAsync(
			t.Context(),
			"saflas",
			"3month",
			7,
			50,
		)

		assert.NoError(t, err, "Expected no error, got %q", err)
		t.Log(res)
	})

	t.Run("Get UserTopTracks", func(t *testing.T) {
		res, err := apiHandler.GetUserTopTracks(
			t.Context(),
			"saflas",
			"3month",
			7,
			50,
		)

		assert.NoError(t, err, "Expected no error, got %q", err)
		t.Log(res)
	})

	t.Run("See if async and seq UserTopTracks are equal", func(t *testing.T) {
		res, err := apiHandler.GetUserTopTracksAsync(
			t.Context(),
			"saflas",
			"3month",
			3,
			50,
		)

		res2, err2 := apiHandler.GetUserTopTracks(
			t.Context(),
			"saflas",
			"3month",
			3,
			50,
		)

		assert.NoError(t, err, "Expected no error, got %q", err)
		assert.NoError(t, err2, "Expected no error, got %q", err2)
		assert.Equal(t, res, res2, "Expected async and sequential results to be equal")
	})

}

func checkResponseOK(resp *http.Response, err error) (*http.Response, error) {
	if resp.StatusCode != 200 {
		return resp, fmt.Errorf("Non-Ok Reponse from API: %q", resp.Status)
	}

	return resp, nil
}

func checkResponseBody(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error in decoding response body: %q", err)
	}

	if printResponse := false; printResponse {
		var obj map[string]interface{}
		if err := json.Unmarshal(body, &obj); err != nil {
			return (err)
		}

		count := 0
		limited := make(map[string]interface{})

		for k, v := range obj {
			limited[k] = v
			count++
			if count == 1 {
				break
			}
		}

		partial, _ := json.MarshalIndent(limited, "", "  ")
		fmt.Println(string(partial))
	}

	return nil
}
