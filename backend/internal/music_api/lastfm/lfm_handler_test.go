package musicapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

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
	)

	t.Run("Get User Weekly Chart List", func(t *testing.T) {
		extraURLParams := map[string]string{
			"method": "user.getweeklychartlist",
			"user":   "saflas",
			"from":   "1751198400",
			"to":     "1751803200",
		}

		resp, err := checkResponseOK(apiHandler.MakeRequest(extraURLParams))

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

		resp, err := checkResponseOK(apiHandler.MakeRequest(extraURLParams))

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

		resp, err := checkResponseOK(apiHandler.MakeRequest(extraURLParams))

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

		resp, err := checkResponseOK(apiHandler.MakeRequest(extraURLParams))

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

		resp, err := checkResponseOK(apiHandler.MakeRequest(extraURLParams))

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
