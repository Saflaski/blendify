package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"context"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

type BlendService struct {
	repo           *RedisStateStore
	LastFMExternal *musicapi.LastFMAPIExternal
}

func (s *BlendService) GenerateNewLinkAndAssignToUser(context context.Context, userA UUID) (any, error) {

	//Generate a linkId to be returned that won't hash collide
	// newInviteId := uuid.New()
	//Store the association

	// s.repo.SetUserToLink(userA, newInviteId)

	return "", nil
}

func (s *BlendService) AddBlendFromInvite(context context.Context, userA UUID, blendLinkValue string) error {

	return nil
}

type blendCategory string
type blendTimeDuration string

var durationMap = map[blendTimeDuration]musicapi.Period{
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

	glog.Info("Calculating blend for users: ", userNameA, " + ", userB, " category: ",
		category, " timeDuration: ", timeDuration)

	switch {
	case category == BlendCategoryArtist:
		return s.getArtistBlend(userNameA, userB, timeDuration)
	case category == BlendCategoryTrack:
		return s.getTrackBlend(userNameA, userB, timeDuration)
	case category == BlendCategoryAlbum:
		return s.getAlbumBlend(userNameA, userB, timeDuration)
	default:
		return 0, fmt.Errorf("category does not match any of the required categories")
	}

}

// ========== Artist Blend ==========
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
	blendNumber := CalculateLWCS(0.8, listenHistoryA, listenHistoryB)
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

// ========== Album Blend ==========
func (s *BlendService) getAlbumBlend(userA, userB string, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopAlbums(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top albums for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.getTopAlbums(userB, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top albums for %s as it returned error: %w", userB, err)
	}
	if len(listenHistoryA) == 0 || len(listenHistoryB) == 0 {
		return 0, fmt.Errorf("inappropriate listen history ranges, userA: %d , userB: %d", len(listenHistoryA), len(listenHistoryB))
	}
	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.8, listenHistoryA, listenHistoryB)
	return blendNumber, nil
}

func (s *BlendService) getTopAlbums(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
	albumToPlays := make(map[string]int)
	topAlbums, err := s.LastFMExternal.GetUserTopAlbums(
		userName,
		durationMap[timeDuration],
		1,
		50,
	)

	if err != nil {
		return albumToPlays, fmt.Errorf("could not extract TopAlbums object from lastfm adapter, %w", err)
	}
	for _, v := range topAlbums.TopAlbums.Album {
		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return albumToPlays, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}
		albumToPlays[v.Name] = playcount
	}

	return albumToPlays, nil
}

// ========== Track Blend ==========

func (s *BlendService) getTrackBlend(userA, userB string, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopTracks(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.getTopTracks(userB, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userB, err)
	}

	if len(listenHistoryA) == 0 || len(listenHistoryB) == 0 {
		return 0, fmt.Errorf("inappropriate listen history ranges, userA: %d , userB: %d", len(listenHistoryA), len(listenHistoryB))
	}
	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.8, listenHistoryA, listenHistoryB)
	return blendNumber, nil
}

func (s *BlendService) getTopTracks(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
	trackToPlays := make(map[string]int)
	topTracks, err := s.LastFMExternal.GetUserTopTracks(
		userName,
		durationMap[timeDuration],
		1,
		50,
	)

	if err != nil {
		return trackToPlays, fmt.Errorf("could not extract TopTracks object from lastfm adapter, %w", err)
	}
	for _, v := range topTracks.TopTracks.Track {
		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return trackToPlays, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}
		trackToPlays[v.Name] = playcount
	}

	return trackToPlays, nil
}
