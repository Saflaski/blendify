package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

type BlendService struct {
	repo           *RedisStateStore
	LastFMExternal *musicapi.LastFMAPIExternal
}

type blendCategory string
type blendTimeDuration string

var durationMap = map[blendTimeDuration]musicapi.Period{ //TODO Delete this shit
	BlendTimeDurationOneMonth:   musicapi.ONE_MONTH,
	BlendTimeDurationThreeMonth: musicapi.THREE_MONTHS,
	BlendTimeDurationYear:       musicapi.YEAR,
}

const (
	BlendCategoryArtist blendCategory = "artist"
	BlendCategoryTrack  blendCategory = "track"
	BlendCategoryAlbum  blendCategory = "album"
)

const (
	BlendTimeDurationOneMonth   blendTimeDuration = "1month"
	BlendTimeDurationThreeMonth blendTimeDuration = "3month"
	// BlendTimeDurationSixMonth   blendTimeDuration = "6month"
	BlendTimeDurationYear blendTimeDuration = "12month"
	// BlendTimeDurationAllTime    blendTimeDuration = "alltime"
)

func NewBlendService(blendStore RedisStateStore, lfmAdapter musicapi.LastFMAPIExternal) *BlendService {
	return &BlendService{&blendStore, &lfmAdapter}
}

func (s *BlendService) GetBlend(userA UUID, userB string, category blendCategory, timeDuration blendTimeDuration) (int, error) {
	//Implement logic to calculate blend percentage based on user data, category, and time duration

	//Get the username from the UUID of the given user that's sending the request
	userNameA, err := s.repo.GetUser(userA)
	if err != nil {
		return 0, fmt.Errorf("could not extract username from UUID of user with ID: %s, %w", userA, err)
	}

	glog.Info("Calculating blend for users: ", userA, " + ", userB, " category: ",
		category, " timeDuration: ", timeDuration)

	switch {
	case category == BlendCategoryArtist:
		glog.Info("Switched to getArtistBlend")
		return s.getArtistBlend(userNameA, userB, timeDuration)
	case category == BlendCategoryTrack:
		return 404, nil
	case category == BlendCategoryAlbum:
		return 404, nil
	default:
		return 0, fmt.Errorf("category does not match any of the required categories")
	}

}

func (s *BlendService) getArtistBlend(userA, userB string, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopArtists(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.getTopArtists(userB, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userB, err)
	}
	if len(listenHistoryA) == 0 || len(listenHistoryB) == 0 {
		return 0, fmt.Errorf("inappropriate listen history ranges, userA: %d , userB: %d", len(listenHistoryA), len(listenHistoryB))
	}

	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.5, listenHistoryA, listenHistoryB)
	return blendNumber, nil
}

func (s *BlendService) getTrackBlend(userA, userB string, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopTracks(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.getTopTracks(userB, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userB, err)
	}
	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.5, listenHistoryA, listenHistoryB)
	return blendNumber, nil
}

func (s *BlendService) getTopArtists(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
	artistToPlaybacks := make(map[string]int)
	topArtist, err := s.LastFMExternal.GetUserTopArtists(
		userName,
		durationMap[timeDuration],
		1,
		50,
	)

	if err != nil {
		return artistToPlaybacks, fmt.Errorf("could not extract TopArtists object from lastfm adapter, %w", err)
	}
	for _, v := range topArtist.TopArtists.Artist {
		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return artistToPlaybacks, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}
		artistToPlaybacks[v.Name] = playcount
	}

	return artistToPlaybacks, nil
}

func (s *BlendService) getTopTracks(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
	//TODO UNIMPLEMENTED
	panic("unimplemented")
	// tracksToPlaybacks := make(map[string]int)
	// topTracks, err := s.LastFMExternal.GetUserTopTracks(
	// 	userName,
	// 	durationMap[timeDuration],
	// 	1,
	// 	50,
	// )

	// if err != nil {
	// 	return tracksToPlaybacks, fmt.Errorf("could not extract TopTracks object from lastfm adapter, %w", err)
	// }
	// for _, v := range topTracks.TopTracks.Artist {
	// 	tracksToPlaybacks[v.Name] = v.Playcount
	// }

	// return tracksToPlaybacks, nil
}
