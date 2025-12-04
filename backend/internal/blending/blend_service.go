package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/golang/glog"
)

type BlendService struct {
	repo           *RedisStateStore
	LastFMExternal *musicapi.LastFMAPIExternal
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

var durationRange = []blendTimeDuration{
	BlendTimeDurationOneMonth,
	BlendTimeDurationThreeMonth,
	BlendTimeDurationYear,
}

var categoryRange = []blendCategory{
	BlendCategoryAlbum,
	BlendCategoryArtist,
	BlendCategoryTrack,
}

func (s BlendService) NewBlend(context context.Context, userA userid, link blendLinkValue) (blendId, error) {

	//First check if this is an existing invite link
	//Then if it is a new link, create a new blendId object
	//If it is an existing link, check if we need to refresh it
	//Return a new blendId
	return "blendId", nil
}

func (s *BlendService) IsExistingLink(context context.Context, link blendLinkValue) (blendId, error) {
	//if link exists, return blendId
	//else nil
	return "RANDOMID1001", nil
}

func (s *BlendService) RefreshLinkIfExpired(context context.Context, id blendId) (blendLinkValue, error) {
	//If the blend under the id has expired, then we need to repopulate the blend
	//We call a populate method on a blendId
	return "", nil
}

func (s *BlendService) PopulateBlend(context context.Context, id blendId) (blendId, error) {
	//Design decisions:
	//This method will be called when a blend ID is either new
	//or its contents have expired.
	//By expired, we mean that the underlying calculations have expired
	//which just means we need to repopulate the user data, then
	//then recalculate
	//Defining expired: when one of the user's data has expired.
	//Need to consider case of 1: User clicking on existing blend
	//2: User adding a new link
	//3: UserA

	//Measuring expired: TTL? Won't maintain consistency across blends then?
	//Just make it be based directly off of the user data TTL OR
	//Make it be floor(userA.Expiry, userB.Expiry..userN.Expiry)

	//TODO revisit this after writing the underlying functions

	return "", nil
}

func (s *BlendService) PopulateUserData(context context.Context, user userid) error {
	ok, err := s.HasUserDataExpired(context, user)
	if err != nil {
		return err
	}

	if ok {
		//Call the getXData functions
		return nil
	}
	return nil

}

func (s *BlendService) HasUserDataExpired(context context.Context, user userid) (bool, error) {
	//Design:
	//Score system or TTL for checking expiry?
	//If score system then we need to check for score > Now - ExpiryDuration
	//If TTL, then simpy if it exists or not. TTL may not be supported for the data we want.
	return true, nil
}

type response struct {
	user     userid
	chart    map[string]int    //We can get category of map from this
	duration blendTimeDuration //And category of blend duration from this
	category blendCategory
	err      error
}

func (s *BlendService) GetNewDataForUser(ctx context.Context, user userid) error {

	platformUsername, err := s.repo.GetUser(UUID(user))
	if err != nil {
		return fmt.Errorf("could not find user by userid when getting new data: %w", err)
	}

	requestSize := len(durationRange) * len(categoryRange)
	respc := make(chan response, requestSize)
	var wg sync.WaitGroup

	for _, duration := range durationRange {
		for _, category := range categoryRange {
			d, c := duration, category
			wg.Add(1)
			go func() {
				defer wg.Done()
				respMap, err := s.downloadLFMData(ctx, platformUsername, d, c)
				if len(respMap) == 0 {
					//wrap the error regardless of it is nil/not nil
					err = fmt.Errorf(" downloaded empty map from platform: %w", err)
				}
				resp := response{
					user:     user,
					chart:    respMap,
					duration: d,
					category: c,
					err:      err,
				}
				respc <- resp
			}()
		}
	}
	go func() {
		wg.Wait()
		close(respc)
	}()

	for resp := range respc {
		// resp := <-respc
		if resp.err != nil {
			return fmt.Errorf("error in downloading data asynchronously: %w", resp.err)
		}
		if err := s.cacheLFMData(ctx, resp); err != nil {
			return fmt.Errorf("could not cache data: %w", err)
		}
	}
	return nil
}

func (s *BlendService) cacheLFMData(ctx context.Context, resp response) error {
	fmt.Println(resp.category)
	fmt.Println(resp.duration)
	fmt.Println(resp.chart)
	fmt.Println("___________________________________________")

	err := s.repo.CacheUserMusicData(ctx, resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *BlendService) downloadLFMData(context context.Context, user string, timePeriod blendTimeDuration, category blendCategory) (map[string]int, error) {

	_ = context //TODO: Change the request methods to accept and use a context
	switch category {
	case BlendCategoryArtist:
		return s.getTopArtists(user, timePeriod)

	case BlendCategoryTrack:
		return s.getTopTracks(user, timePeriod)

	case BlendCategoryAlbum:
		return s.getTopAlbums(user, timePeriod)
	default:
		return nil, fmt.Errorf("invalid category")
	}
}

// func (s *BlendService) downloadData(context context.Context, userid userid, duration blendTimeDuration)

func (s *BlendService) GenerateNewLinkAndAssignToUser(context context.Context, userA userid) (any, error) {

	//Generate a linkId to be returned that won't hash collide
	// newInviteId := uuid.New()
	//Store the association

	// s.repo.SetUserToLink(userA, newInviteId)

	return "", nil
}

func (s *BlendService) AddBlendFromInvite(context context.Context, userA userid, blendLinkValue string) error {

	return nil
}

func NewBlendService(blendStore RedisStateStore, lfmAdapter musicapi.LastFMAPIExternal) *BlendService {
	return &BlendService{&blendStore, &lfmAdapter}
}

// TODO: Delete this function or change UUID to userid type
func (s *BlendService) GetBlend(userA UUID, userB string, category blendCategory, timeDuration blendTimeDuration) (int, error) {
	//Implement logic to calculate blend percentage based on user data, category, and time duration

	//Get the username from the UUID of the given user that's sending the request
	userNameA, err := s.repo.GetUser(userA)
	if err != nil {
		return 0, fmt.Errorf("could not extract username from UUID of user with ID: %s, %w", userA, err)
	}

	glog.Info("Calculating blend for users: ", userNameA, " + ", userB, " category: ",
		category, " timeDuration: ", timeDuration)

	switch category {
	case BlendCategoryArtist:
		return s.getArtistBlend(userNameA, userB, timeDuration)
	case BlendCategoryTrack:
		return s.getTrackBlend(userNameA, userB, timeDuration)
	case BlendCategoryAlbum:
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

func (s *BlendService) getTopArtistsNew(userName string, timeDuration blendTimeDuration) (artistChart, error) {
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

	return artistChart(artistToPlaybacks), nil
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

func (s *BlendService) getTopAlbumsNew(userName string, timeDuration blendTimeDuration) (albumChart, error) {
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

	return albumChart(albumToPlays), nil
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

func (s *BlendService) getTopTracksNew(userName string, timeDuration blendTimeDuration) (trackChart, error) {
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

	return trackChart(trackToPlays), nil
}
