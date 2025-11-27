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

func DecodeRequest[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decoding JSON: %w", err)
	}
	return v, nil
}
