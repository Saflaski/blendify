package utility

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(input string) map[string]interface{} {

	return map[string]interface{}{
		"key": "value",
	}
}

// Function to decode JSON response from a http.Response
func Decode[T any](response *http.Response) (T, error) {
	var v T
	if err := json.NewDecoder(response.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decoding JSON: %w", err)
	}
	return v, nil

}

func StructToJSONBytes(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func DecodeRequest[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decoding JSON: %w", err)
	}
	return v, nil
}

func MapToJSON(m map[string]int) ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ObjectToJSON[T any](m T) ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func JSONToMapCatStats(data []byte) (map[string]CatalogueStats, error) {
	var out map[string]CatalogueStats
	err := json.Unmarshal(data, &out)
	return out, err
}

type CatalogueStats struct { //A catalogue can be an album, track or artist. The following is metadata for a catalogue
	Artist      string `json:"artist"`
	Count       int    `json:"count"`
	PlatformURL string `json:"platformurl"` //Catalogue URL
	Image       string `json:"imageurl"`    //Image URL
	PlatformID  string `json:"platformid"`  //Catalogue Platform ID
}

func JSONToMap(data []byte) (map[string]int, error) {
	var out map[string]int
	err := json.Unmarshal(data, &out)
	return out, err
}
