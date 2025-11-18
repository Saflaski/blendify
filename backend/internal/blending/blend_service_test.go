package blend

import (
	"fmt"
	"testing"
	"time"
)


type StubBlendService struct {

}

func TestGetBlend(t *testing.T) {
	blendService := NewBlendService()

	//Mock Data
	t.Run("Set User - Top Listened to Artists", func(t *testing.T) {
		uuid := "test-uuid-1234"
		artists := []string{"Artist1", "Artist2", "Artist3"}
		err := blendService.SetUserTopArtists(uuid, artists)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Get User - Top Listened to Artists in last week", func(t *testing.T) {
		weekDuration := time.Duration(7*24) * time.Hour
		uuid := "test-uuid-1234"
		response, err := blendService.GetUserTopArtists(uuid, weekDuration)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		fmt.Printf("Response: %v\n", response)
	})

}

// func NewStubBlendService() *StubBlendService {
// 	return &StubBlendService{

// 	}
// }