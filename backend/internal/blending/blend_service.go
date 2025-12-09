package blend

import (
	"backend-lastfm/internal/auth"
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type BlendService struct {
	repo           *RedisStateStore
	LastFMExternal *musicapi.LastFMAPIExternal
	authRepo       *auth.AuthStateStore
}

const BLEND_USER_LIMIT = 2

func (s *BlendService) GetDuoBlendData(context context.Context, blendId blendId) (DuoBlend, error) {

	//Get json data for percentage data of 3x3 data
	//Get json data for percentage data of distribution of art/alb/tra

	userids, err := s.repo.GetUsersFromBlend(context, blendId)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" error getting users from blend id: %s, err: %w", blendId, err)
	}

	err = s.PopulateUsersByBlend(context, blendId)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" Could not populate user data: %w", err)
	}

	// individualUserData := make([]IndividualUserData, len(userids))
	// errSum := 0
	// size := len(userids)
	// totalPairs := int(size * (size - 1) / 2) //1
	// allBlends := make([]DuoBlend, 0, totalPairs)
	// for i := 0; i < size; i++ {
	// 	for j := i + 1; j < size; j++ {
	duoBlend, err := s.GenerateBlendOfTwo(context, userids[0], userids[1])
	// 		allBlends = append(allBlends, duoBlend)
	// 		errSum += errAddition
	// 	}
	// }

	if err != nil {
		return DuoBlend{}, fmt.Errorf(" Could not generate blend due to one or more reasons with getting data from db/platform")
	}

	return duoBlend, nil

}

func (s *BlendService) getLFM(ctx context.Context, userID string) (string, error) {
	username, err := s.authRepo.GetLFMByUserId(ctx, userID)
	if err != nil {
		glog.Errorf("Could not extract platform username from userid: %s", userID)
		return "", err
	}
	return username, nil
}

func (s *BlendService) GenerateBlendOfTwo(context context.Context, userA userid, userB userid) (DuoBlend, error) {

	usernameA, err := s.getLFM(context, string(userA))
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" could not get username %s", userA)
	}

	usernameB, err := s.getLFM(context, string(userB))
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" could not get username %s", userB)
	}

	artistBlend, err := s.buildArtistBlend(usernameA, usernameB)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" failed to get artist blend: %w", err)
	}
	albumBlend, err := s.buildAlbumBlend(usernameA, usernameB)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" failed to get album blend: %w", err)
	}
	trackBlend, err := s.buildTrackBlend(usernameA, usernameB)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" failed to get track blend: %w", err)
	}
	duoBlend := DuoBlend{Users: []string{usernameA, usernameB},
		ArtistBlend: artistBlend,
		AlbumBlend:  albumBlend,
		TrackBlend:  trackBlend,
		//TODO More to be added
	}

	return duoBlend, nil
}

func (s *BlendService) buildArtistBlend(usernameA, usernameB string) (TypeBlend, error) {
	var (
		b   TypeBlend
		err error
	)

	if b.ThreeMonth, err = s.getArtistBlend(usernameA, usernameB, BlendTimeDurationOneMonth); err != nil {
		glog.Errorf("Could not get 1-month artist blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 1-month artist blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.SixMonth, err = s.getArtistBlend(usernameA, usernameB, BlendTimeDurationThreeMonth); err != nil {
		glog.Errorf("Could not get 3-month artist blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 3-month artist blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.OneYear, err = s.getArtistBlend(usernameA, usernameB, BlendTimeDurationYear); err != nil {
		glog.Errorf("Could not get 12-month artist blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 12-month artist blend for %s, %s: %v", usernameA, usernameB, err)
	}

	return b, nil
}

func (s *BlendService) buildAlbumBlend(usernameA, usernameB string) (TypeBlend, error) {
	var (
		b   TypeBlend
		err error
	)

	if b.ThreeMonth, err = s.getAlbumBlend(usernameA, usernameB, BlendTimeDurationOneMonth); err != nil {
		glog.Errorf("Could not get 1-month album blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 1-month album blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.SixMonth, err = s.getAlbumBlend(usernameA, usernameB, BlendTimeDurationThreeMonth); err != nil {
		glog.Errorf("Could not get 3-month album blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 3-month album blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.OneYear, err = s.getAlbumBlend(usernameA, usernameB, BlendTimeDurationYear); err != nil {
		glog.Errorf("Could not get 12-month album blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 12-month album blend for %s, %s: %v", usernameA, usernameB, err)
	}

	return b, nil
}

func (s *BlendService) buildTrackBlend(usernameA, usernameB string) (TypeBlend, error) {
	var (
		b   TypeBlend
		err error
	)

	if b.ThreeMonth, err = s.getTrackBlend(usernameA, usernameB, BlendTimeDurationOneMonth); err != nil {
		glog.Errorf("Could not get 1-month track blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 1-month track blend for %s, %s: %v", usernameA, usernameB, err)

	}

	if b.SixMonth, err = s.getTrackBlend(usernameA, usernameB, BlendTimeDurationThreeMonth); err != nil {
		glog.Errorf("Could not get 3-month track blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 3-month track blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.OneYear, err = s.getTrackBlend(usernameA, usernameB, BlendTimeDurationYear); err != nil {
		glog.Errorf("Could not get 12-month track blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf("Could not get 12-month track blend for %s, %s: %v", usernameA, usernameB, err)
	}

	return b, nil
}

func (s *BlendService) AuthoriseBlend(context context.Context, blendId blendId, id userid) (bool, error) {
	ok, err := s.repo.IsUserInBlend(context, id, blendId)
	if err != nil {
		return false, fmt.Errorf(" could not check if user is in blend: %w", err)
	}

	return ok, err
}

func (s *BlendService) AddOrMakeBlendFromLink(context context.Context, userA userid, link blendLinkValue) (blendId, error) {

	//First check if this is an existing invite link
	//Then if it is a new link, create a new blendId object
	//MakeNewBlend([]users{userA, userB})

	id, err := s.repo.IsExistingBlendFromLink(context, link)
	if err != nil {
		return "", fmt.Errorf(" error during checking if blendlink existed: %w", err)

	}
	glog.Infof("Found blend from link: %s", link)
	if id == "" { //No link found
		userB, err := s.repo.GetLinkCreator(context, link) //Fetch user who created link
		if err != nil {
			return "", fmt.Errorf(" error during getting user (creator) from link : %w", err)
		}
		glog.Infof("Blend created by: %s", userB)

		//Safety net to make sure userA != userB
		if userB == userA {
			glog.Info("Same user nvm")
			return "0", nil //0 is code for consuming user being the same user as creating user
		}

		// TEMPORARY LIMIT FOR NUM USERS WHO CAN BE IN A BLEND
		userids, err := s.repo.GetUsersFromBlend(context, id)
		if err != nil {
			return "", fmt.Errorf(" error getting users from blend id: %s, err: %w", id, err)
		}
		if len(userids)+1 > BLEND_USER_LIMIT {
			glog.Info("Same user nvm")
			return "-1", nil
		}

		// This fetched user + the user who resulted in this function being called are
		// now the first two users of this new blend

		// We only create the blend id for now as generating a whole new blend might time out
		// and if we tried to make it async, then there will be a race condition between
		// frontend loading blend page + backend trying to hydrate the blend
		id, err = s.GenerateNewBlendId(context, []userid{userA, userB}) //Should this make the whole blend or?
		if err != nil {
			return "", fmt.Errorf(" error during making a blend with users %s and %s: %w", userA, userB, err)
		}
		glog.Infof("Generating new blend: %s", id)
		return id, nil
	} else {
		//Check if user is already in this blend
		ok, err := s.repo.IsUserInBlend(context, userA, id)
		if err != nil {
			return "", fmt.Errorf(" could not check if user is in blend: %w", err)
		}
		if !ok {
			err := s.repo.AddUsersToBlend(context, id, []userid{userA})
			if err != nil {
				return "", fmt.Errorf(" could not add user to blend: %w", err)
			}
			glog.Infof("User does not exist in blend. Adding %s", userA)
			return id, nil
		} else {
			glog.Infof("User already exists in blend")
			//Nothing to see here, just return the existing blend id
			return id, nil
		}
	}

}

func (s *BlendService) GenerateNewBlendId(context context.Context, userids []userid) (blendId, error) {
	id := blendId(uuid.New().String())
	err := s.repo.AddUsersToBlend(context, id, userids)
	if err != nil {
		return "", fmt.Errorf(" error during inserting new blend frame: %w", err)
	}

	return id, err
}

func (s *BlendService) PopulateUsersByBlend(context context.Context, id blendId) error {
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

	userids, err := s.repo.GetUsersFromBlend(context, id)
	if err != nil {
		return fmt.Errorf(" error getting users from blend id: %s, err: %w", id, err)
	}

	//Check which userids have expired user data and populate them
	for _, user := range userids {
		ok, err := s.repo.UserHasAnyMusicData(context, user)
		if err != nil {
			return fmt.Errorf(" error during checking if user: %s has any music data: %w", id, err)
		}
		if !ok {
			err := s.GetNewDataForUser(context, user)
			if err != nil {
				return fmt.Errorf(" in PopulateBlend, could not get new data for user: %s with err: %w", user, err)
			}
		}
	}

	Blend, err := s.MakeNewBlend(context, userids) //Make a blend struct from all the userids
	if err != nil {
		return fmt.Errorf(" error during making new blend from userids: %w", err)
	}
	if err := s.CacheBlend(context, &Blend); err != nil {
		return fmt.Errorf(" error during caching blend: %s with err: %w", Blend.id, err)
	}

	return nil
}

func (s *BlendService) CacheBlend(context context.Context, blend *Blend) error {
	return nil
}

func (s *BlendService) MakeNewBlend(context context.Context, userids []userid) (Blend, error) {
	return Blend{}, nil
}

func (s *BlendService) PopulateUserData(context context.Context, user userid) error {

	//This function needs to be looked at again in the future for addition of granular
	//cache entry checking and granular cache hydration
	//Particularly, a secondary check if UserHasAnyMusicData -> true which is a
	//a different check function such as GetEachExpiredCacheEntryByUser which
	//needs to be used and then those keys need to be filled

	ok, err := s.repo.UserHasAnyMusicData(context, user)
	if err != nil {
		return fmt.Errorf(" error during checking if user has any music data: %w", err)
	}

	if !ok {
		err := s.GetNewDataForUser(context, user)
		if err != nil {
			return fmt.Errorf(" error during getting full data for user: %w", err)
		}
		return nil
	} else {
		return nil //Do nothing for now.
	}

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
		return s.downloadTopArtists(user, timePeriod)

	case BlendCategoryTrack:
		return s.downloadTopTracks(user, timePeriod)

	case BlendCategoryAlbum:
		return s.downloadTopAlbums(user, timePeriod)
	default:
		return nil, fmt.Errorf("invalid category")
	}
}

// func (s *BlendService) downloadData(context context.Context, userid userid, duration blendTimeDuration)

func (s *BlendService) GenerateNewLinkAndAssignToUser(context context.Context, userA userid) (blendLinkValue, error) {

	//Generate a linkId to be returned that won't hash collide
	newInviteValue := blendLinkValue(uuid.New().String())
	err := s.repo.SetUserToLink(context, userA, newInviteValue)
	if err != nil {
		return "", fmt.Errorf(" could not set user to link: %w", err)
	}

	return newInviteValue, nil
}

func (s *BlendService) AddBlendFromInvite(context context.Context, userA userid, blendLinkValue string) error {

	return nil
}

func NewBlendService(blendStore RedisStateStore, lfmAdapter musicapi.LastFMAPIExternal, authStore auth.AuthStateStore) *BlendService {
	return &BlendService{&blendStore, &lfmAdapter, &authStore}
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

	listenHistoryA, err := s.downloadTopArtists(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.downloadTopArtists(userB, timeDuration)
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

func (s *BlendService) downloadTopArtists(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
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

	listenHistoryA, err := s.downloadTopAlbums(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top albums for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.downloadTopAlbums(userB, timeDuration)
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

func (s *BlendService) downloadTopAlbums(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
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

	listenHistoryA, err := s.downloadTopTracks(userA, timeDuration)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.downloadTopTracks(userB, timeDuration)
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

func (s *BlendService) downloadTopTracks(userName string, timeDuration blendTimeDuration) (map[string]int, error) {
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
