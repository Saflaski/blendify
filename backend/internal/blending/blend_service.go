package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type BlendService struct {
	repo           *RedisStateStore
	LastFMExternal *musicapi.LastFMAPIExternal
	// authRepo       *auth.AuthStateStore
}

// Deletes blend on backend for all users
func (s *BlendService) DeleteBlend(context context.Context, user userid, blendId blendId) error {

	blendUsers, err := s.repo.GetUsersFromBlend(context, blendId)
	if err != nil {
		return fmt.Errorf(" could not find blendusers from blendid %s : %w", blendId, err)
	}

	for _, userFromBlend := range blendUsers {
		err := s.repo.DeleteBlendByBlendId(context, userFromBlend, blendId)
		if err != nil {
			return fmt.Errorf(" could not delete blend: %w", err)
		}
	}

	return nil
}

func (s *BlendService) GetUserBlends(context context.Context, user userid) (Blends, error) {
	blendIds, err := s.repo.GetBlendsByUser(context, user)
	if err != nil {
		return Blends{}, fmt.Errorf(" could not find blends from userid %s : %w", user, err)
	}
	blendAccumulator := make([]Blend, len(blendIds))
	for i, v := range blendIds {
		blendAccumulator[i].BlendId = string(v)
		blendUsers, err := s.repo.GetUsersFromBlend(context, v)
		if err != nil {
			return Blends{}, fmt.Errorf(" could not find blendusers from blendid %s : %w", v, err)
		}
		blendPlatformUsernames := make([]platformid, len(blendUsers))
		for j, v_2 := range blendUsers {
			platformUser, err := s.repo.GetLFMByUserId(context, string(v_2))
			if err != nil {
				return Blends{}, fmt.Errorf(" could not extract platformid from userid %s : %w", v_2, err)
			}
			blendPlatformUsernames[j] = platformid(platformUser)

		}
		blendAccumulator[i].Users = blendPlatformUsernames

		//REMOVE this part when LIMIT>2
		if BLEND_USER_LIMIT == 2 {
			// res, err := s.GenerateBlendOfTwo(context, blendUsers[0], blendUsers[1])
			// if err != nil {
			// 	return Blends{}, fmt.Errorf(" could not generate blendnumber from blendid %s : %w", v, err)
			// }
			overallVal, err := s.repo.GetCachedOverallBlend(context, v)
			if err != nil {
				glog.Errorf(" could not get cached overallblend: %s", err)
			} else {
				blendAccumulator[i].Value = overallVal
			}
		}
		timeRes, err := s.repo.GetBlendTimeStamp(context, v)
		if err != nil {
			glog.Errorf(" could not get cached timestamp: %s", err)
		} else {
			blendAccumulator[i].CreatedAt = timeRes
		}

	}
	allBlends := Blends{Blends: blendAccumulator}
	return allBlends, nil
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
	username, err := s.repo.GetLFMByUserId(ctx, userID)
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

	artistBlend, err := s.buildArtistBlend(context, userA, userB)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" failed to get artist blend: %w", err)
	}
	albumBlend, err := s.buildAlbumBlend(context, userA, userB)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" failed to get album blend: %w", err)
	}
	trackBlend, err := s.buildTrackBlend(context, userA, userB)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" failed to get track blend: %w", err)
	}

	overallBlendNum, err := s.buildOverallBlend(artistBlend, albumBlend, trackBlend)
	if err != nil {
		return DuoBlend{}, fmt.Errorf("could not get overall blend with %s and %s: %w", userA, userB, err)
	}

	// blendedArtists, err := s.BuildBlendedEntries(
	// 	context,
	// 	userA,
	// 	userB,
	// 	BlendTimeDurationYear,
	// 	BlendCategoryArtist,
	// 	25,
	// )
	// _ = blendedArtists //TODO use this later

	duoBlend := DuoBlend{Users: []string{usernameA, usernameB},
		OverallBlendNum: overallBlendNum,
		ArtistBlend:     artistBlend,
		AlbumBlend:      albumBlend,
		TrackBlend:      trackBlend,
		//TODO More to be added
	}

	return duoBlend, nil
}

func (s *BlendService) GetBlendEntryByBlendId(context context.Context, blendId blendId, category blendCategory, duration blendTimeDuration) ([]TopEntry, error) {

	userids, err := s.repo.GetUsersFromBlend(context, blendId)
	if err != nil {
		return nil, fmt.Errorf(" error getting users from blend id: %s, err: %w", blendId, err)
	}

	if len(userids) != 2 {
		return nil, fmt.Errorf(" only blends of 2 users are supported currently")
	}

	blendedEntries, err := s.BuildBlendedEntries(context, userids[0], userids[1], duration, category, 25)
	if err != nil {
		return nil, fmt.Errorf(" could not build blended entries for blendid %s: %w", blendId, err)
	}

	return blendedEntries, nil
}

func (s *BlendService) BuildBlendedEntries(context context.Context, userA, userB userid, duration blendTimeDuration, category blendCategory, minimum int) ([]TopEntry, error) {

	//Get common and unique artists between both users
	aEntries, err := s.getTopX(context, userA, duration, category)
	if err != nil {
		return nil, fmt.Errorf(" could not get top entries for userA %s: %w", userA, err)
	}
	bEntries, err := s.getTopX(context, userB, duration, category)
	if err != nil {
		return nil, fmt.Errorf(" could not get top entries for userB %s: %w", userB, err)
	}

	entries := s.GetCommonAndSortedEntries(aEntries, bEntries)

	blendedEntries := make([]TopEntry, len(entries))
	for k, v := range entries {
		aStat := aEntries[v]
		bStat := bEntries[v]
		blendedEntries[k] = s.ConvertCatalogueStatsToEntry(aStat, v, aStat.Count, bStat.Count)
	}
	return blendedEntries, nil
}

func (s *BlendService) ConvertCatalogueStatsToEntry(aStat CatalogueStats, name string, aCount int, bCount int) TopEntry {
	countA := aCount
	countB := bCount

	entry := TopEntry{
		Name:           name,
		ImageURL:       aStat.Image,
		URL:            aStat.PlatformURL,
		ArtistName:     aStat.Artist.Name,
		ArtistURL:      aStat.Artist.URL,
		ArtistImageURL: s.getCatalogueImageURL(aStat.Artist.LFMImages),
		Playcounts:     []int{countA, countB},
	}
	return entry
}

func (s *BlendService) GetCommonAndSortedEntries(aEntries map[string]CatalogueStats, bEntries map[string]CatalogueStats) []string {
	commonKeys := FindIntersectKeys(aEntries, bEntries)

	// sortedCommonKeys := make([]string, len(commonKeys))
	// scoreMap := make(map[string]float64)
	// for _, commonKey := range commonKeys {
	// 	left := aEntries[commonKey].Count
	// 	right := aEntries[commonKey].Count
	// 	scoreMap[commonKey] = GetLogCombinedScore(left, right)
	// }

	//Generate a score for each item in commonKeys (itemA and itemB) using GetLogCombineScore
	//Which is essentially (log(X) + log(Y)) * log(Y)/log(X) iff Y > X
	//Then sort list in descending order, i.e higher score at the top
	sortFunc := func(itemA, itemB string) int {
		itemAScore := GetLogCombinedScore(aEntries[itemA].Count, bEntries[itemA].Count)
		itemBScore := GetLogCombinedScore(aEntries[itemB].Count, bEntries[itemB].Count)
		return -cmp.Compare(itemAScore, itemBScore)
	}
	slices.SortFunc(commonKeys, sortFunc)

	return commonKeys
}

func (s *BlendService) buildOverallBlend(artistBlend, albumBlend, trackBlend TypeBlend) (int, error) {
	// We need granular aggregation of the blends hence taking them individually rather than looping

	// artistOverall := calcOverModality(artistBlend, 10, 3, 8) //These numbers are kept like this for more granularity later
	// albumOverall := calcOverModality(albumBlend, 10, 3, 8)
	// trackOverall := calcOverModality(trackBlend, 10, 3, 8)

	artistOverall, err := combineNumbersWithWeights(
		artistBlend.OneMonth,
		artistBlend.ThreeMonth,
		artistBlend.OneYear,
		10, 10, 10)

	if err != nil {
		return 0, fmt.Errorf(" could not calc over artist: %w", err)
	}
	albumOverall, err := combineNumbersWithWeights(
		albumBlend.OneMonth,
		albumBlend.ThreeMonth,
		albumBlend.OneYear,
		10, 10, 10)
	if err != nil {
		return 0, fmt.Errorf(" could not calc over artist: %w", err)
	}

	trackOverall, err := combineNumbersWithWeights(
		trackBlend.OneMonth,
		trackBlend.ThreeMonth,
		trackBlend.OneYear,
		10, 10, 10)
	if err != nil {
		return 0, fmt.Errorf(" could not calc over tracks: %w", err)
	}

	overallBlend, err := combineNumbersWithWeights(artistOverall, albumOverall, trackOverall, 10, 10, 10)
	if err != nil {
		return 0, fmt.Errorf("could not combine overall blend: %w", err)
	}
	return overallBlend, nil
}

func (s *BlendService) buildArtistBlend(context context.Context, usernameA, usernameB userid) (TypeBlend, error) {
	var (
		b   TypeBlend
		err error
	)

	if b.OneMonth, err = s.getArtistBlend(context, usernameA, usernameB, BlendTimeDurationOneMonth); err != nil {
		glog.Errorf("Could not get 1-month artist blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 1-month artist blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.ThreeMonth, err = s.getArtistBlend(context, usernameA, usernameB, BlendTimeDurationThreeMonth); err != nil {
		glog.Errorf("Could not get 3-month artist blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 3-month artist blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.OneYear, err = s.getArtistBlend(context, usernameA, usernameB, BlendTimeDurationYear); err != nil {
		glog.Errorf("Could not get 12-month artist blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 12-month artist blend for %s, %s: %v", usernameA, usernameB, err)
	}

	return b, nil
}

func (s *BlendService) buildAlbumBlend(context context.Context, usernameA, usernameB userid) (TypeBlend, error) {
	var (
		b   TypeBlend
		err error
	)

	if b.OneMonth, err = s.getAlbumBlend(context, usernameA, usernameB, BlendTimeDurationOneMonth); err != nil {
		glog.Errorf("Could not get 1-month album blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 1-month album blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.ThreeMonth, err = s.getAlbumBlend(context, usernameA, usernameB, BlendTimeDurationThreeMonth); err != nil {
		glog.Errorf("Could not get 3-month album blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 3-month album blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.OneYear, err = s.getAlbumBlend(context, usernameA, usernameB, BlendTimeDurationYear); err != nil {
		glog.Errorf("Could not get 12-month album blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 12-month album blend for %s, %s: %v", usernameA, usernameB, err)
	}

	return b, nil
}

func (s *BlendService) buildTrackBlend(context context.Context, usernameA, usernameB userid) (TypeBlend, error) {
	var (
		b   TypeBlend
		err error
	)

	if b.OneMonth, err = s.getTrackBlend(context, usernameA, usernameB, BlendTimeDurationOneMonth); err != nil {
		glog.Errorf("Could not get 1-month track blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 1-month track blend for %s, %s: %v", usernameA, usernameB, err)

	}

	if b.ThreeMonth, err = s.getTrackBlend(context, usernameA, usernameB, BlendTimeDurationThreeMonth); err != nil {
		glog.Errorf("Could not get 3-month track blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 3-month track blend for %s, %s: %v", usernameA, usernameB, err)
	}

	if b.OneYear, err = s.getTrackBlend(context, usernameA, usernameB, BlendTimeDurationYear); err != nil {
		glog.Errorf("Could not get 12-month track blend for %s, %s: %v", usernameA, usernameB, err)
		return TypeBlend{}, fmt.Errorf(" could not get 12-month track blend for %s, %s: %v", usernameA, usernameB, err)
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
		//TODO : THIS IS MAKING PROBLEMS
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
			glog.Info("Tried to add too many users to blend")
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
			err := s.AddUsersToBlend(context, id, []userid{userA})
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

func (s *BlendService) AddUsersToBlend(context context.Context, id blendId, userids []userid) error {
	err := s.repo.AddUsersToBlend(context, id, userids)
	if err != nil {
		return fmt.Errorf(" could not add user to blend: %w", err)
	}

	//Time to cache the overall blend num for retrieval later
	//It does not need to be atomic
	err = s.RefreshOverallBlendInCache(context, id)
	if err != nil {
		return fmt.Errorf(" During refresh blend: %w", err)
	}
	return nil
}

func (s *BlendService) RefreshOverallBlendInCache(context context.Context, id blendId) error {

	duoblend, err := s.GetDuoBlendData(context, id)
	if err != nil {
		glog.Infof(" could not generate duoblend data during adding users to blend: %s: %w", id, err)
		return nil //This is not a fatal error and we can live without it happening but still need to log
	}
	err = s.repo.AssignOverallBlendToBlend(context, id, duoblend.OverallBlendNum)
	if err != nil {
		glog.Infof(" could not assign duoblend integer value during adding users to blend:%s: %w", id, err)
		return nil //This is not a fatal error and we can live without it happening but still need to log
	}
	return nil
}

func (s *BlendService) GenerateNewBlendId(context context.Context, userids []userid) (blendId, error) {
	id := blendId(uuid.New().String())
	err := s.AddUsersToBlend(context, id, userids)
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
		if err := s.PopulateUserData(context, user); err != nil {
			return fmt.Errorf(" error during populating user data for %s : %w", user, err)
		}
	}

	return nil
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

func (s *BlendService) GetNewDataForUser(ctx context.Context, user userid) error {

	platformUsername, err := s.repo.GetLFMByUserId(ctx, string(user))
	if err != nil {
		return fmt.Errorf("could not find user by userid when getting new data: %w", err)
	}

	requestSize := len(durationRange) * len(categoryRange)
	// respc := make(chan response, requestSize)
	respc := make(chan complexResponse, requestSize)
	var wg sync.WaitGroup

	for _, duration := range durationRange {
		for _, category := range categoryRange {
			d, c := duration, category
			wg.Add(1)
			go func() {
				defer wg.Done()
				respData, err := s.downloadTopX(ctx, platformUsername, d, c)
				if len(respData) == 0 {
					//wrap the error regardless of it is nil/not nil
					err = fmt.Errorf(" downloaded empty map from platform: %w", err)
				}
				resp := complexResponse{
					user:     user,
					data:     respData,
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
			return fmt.Errorf("error in downloading data asynchronously for duration %s and category %s : %w", resp.duration, resp.category, resp.err)
		}
		if err := s.cacheLFMData(ctx, resp); err != nil {
			return fmt.Errorf("could not cache data: %w", err)
		}
	}
	return nil
}

func (s *BlendService) cacheLFMData(ctx context.Context, resp complexResponse) error {
	err := s.repo.CacheUserMusicData(ctx, resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *BlendService) downloadLFMData(context context.Context, user string, timePeriod blendTimeDuration, category blendCategory) (map[string]CatalogueStats, error) {

	switch category {
	case BlendCategoryArtist:
		return s.downloadTopArtists(context, user, timePeriod)

	case BlendCategoryTrack:
		return s.downloadTopTracks(context, user, timePeriod)

	case BlendCategoryAlbum:
		return s.downloadTopAlbums(context, user, timePeriod)
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

func NewBlendService(blendStore RedisStateStore, lfmAdapter musicapi.LastFMAPIExternal) *BlendService {
	return &BlendService{&blendStore, &lfmAdapter}
}

func (s *BlendService) GetBlend(context context.Context, userA userid, userB userid, category blendCategory, timeDuration blendTimeDuration) (int, error) {
	//Implement logic to calculate blend percentage based on user data, category, and time duration

	// //Get the username from the UUID of the given user that's sending the request
	// userNameA, err := s.repo.GetUser(userA)
	// if err != nil {
	// 	return 0, fmt.Errorf("could not extract username from UUID of user with ID: %s, %w", userA, err)
	// }

	glog.Info("Calculating blend for users: ", userA, " + ", userB, " category: ",
		category, " timeDuration: ", timeDuration)

	switch category {
	case BlendCategoryArtist:
		return s.getArtistBlend(context, userA, userB, timeDuration)
	case BlendCategoryTrack:
		return s.getTrackBlend(context, userA, userB, timeDuration)
	case BlendCategoryAlbum:
		return s.getAlbumBlend(context, userA, userB, timeDuration)
	default:
		return 0, fmt.Errorf("category does not match any of the required categories")
	}

}

// Converts catalogueKey -> Catalogue => catalogueKey -> playcount int
func (s *BlendService) extractPlayCount(input map[string]CatalogueStats) map[string]int {
	output := make(map[string]int)
	for k, v := range input {
		output[k] = v.Count
	}
	return output
}

// ========== Artist Blend ==========
func (s *BlendService) getArtistBlend(context context.Context, userA, userB userid, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopX(context, userA, timeDuration, BlendCategoryArtist)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}

	listenHistoryB, err := s.getTopX(context, userB, timeDuration, BlendCategoryArtist)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userB, err)
	}
	if len(listenHistoryA) == 0 || len(listenHistoryB) == 0 {
		return 0, fmt.Errorf("inappropriate listen history ranges, userA: %d , userB: %d", len(listenHistoryA), len(listenHistoryB))
	}

	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.8, s.extractPlayCount(listenHistoryA), s.extractPlayCount(listenHistoryB))
	return blendNumber, nil
}

func (s *BlendService) getTopX(context context.Context, userid userid, timeDuration blendTimeDuration, category blendCategory) (map[string]CatalogueStats, error) {
	// totalEntries := 250
	dbResp, err := s.repo.GetFromCacheTopX(context, string(userid), timeDuration, category)
	if err != nil {
		glog.Errorf(" Cache error during getting topX for user %s with duration %s and category %s that needs to be checked: %w", userid, timeDuration, category, err)
	}

	if len(dbResp) != 0 { //Cache hit
		// if len(dbResp) >= totalEntries {
		// 	return dbResp, nil
		// }
		return dbResp, nil
	} else { //Cache miss
		platformUsername, err := s.getLFM(context, string(userid))
		if err != nil {
			return nil, fmt.Errorf(" could get platform username from given userid:%s with err: %w", userid, err)
		}
		glog.Info("Cache miss, so returning downloadTopX")
		lfmResp, err := s.downloadTopX(context, platformUsername, timeDuration, category)
		if err != nil {
			return nil, fmt.Errorf(" did not download %s %s properly: %w", timeDuration, category, err)
		}

		if len(lfmResp) == 0 {
			return nil, fmt.Errorf(" downloaded empty map for %s %s : %w", timeDuration, category, err)
		}

		return lfmResp, nil

	}
}

func (s *BlendService) downloadTopX(context context.Context, userName string, timeDuration blendTimeDuration, category blendCategory) (map[string]CatalogueStats, error) {
	switch category {
	case BlendCategoryAlbum:
		return s.downloadTopAlbums(context, userName, timeDuration)
	case BlendCategoryArtist:
		return s.downloadTopArtists(context, userName, timeDuration)
	case BlendCategoryTrack:
		return s.downloadTopTracks(context, userName, timeDuration)
	default:
		return nil, nil
	}
}

func (s *BlendService) downloadTopArtists(context context.Context, userName string, timeDuration blendTimeDuration) (map[string]CatalogueStats, error) {

	artistToPlaybacks := make(map[string]CatalogueStats)
	topArtist, err := s.LastFMExternal.GetUserTopArtists(
		context,
		userName,
		string(timeDuration),
		6,
		50,
	)

	if err != nil {
		return artistToPlaybacks, fmt.Errorf("could not extract TopArtists object from lastfm adapter, %w", err)
	}
	// for _, v := range topArtist.TopArtists.Artist {
	// 	playcount, err := strconv.Atoi(v.Playcount)
	// 	if err != nil {
	// 		return artistToPlaybacks, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
	// 	}
	// 	artistToPlaybacks[v.Name] = playcount
	// }

	for _, v := range topArtist.TopArtists.Artist {
		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return artistToPlaybacks, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}

		imageURL := s.getCatalogueImageURL(v.LFMImages) //Selects a good pic out of the ones given

		catStat := CatalogueStats{
			Artist:      v,
			Count:       playcount,
			PlatformURL: v.URL,
			Image:       imageURL,
			PlatformID:  v.MBID,
		}
		artistToPlaybacks[v.Name] = catStat

	}

	return artistToPlaybacks, nil
}

// ========== Album Blend ==========
func (s *BlendService) getAlbumBlend(context context.Context, userA, userB userid, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopX(context, userA, timeDuration, BlendCategoryAlbum)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top albums for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.getTopX(context, userB, timeDuration, BlendCategoryAlbum)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top albums for %s as it returned error: %w", userB, err)
	}
	if len(listenHistoryA) == 0 || len(listenHistoryB) == 0 {
		return 0, fmt.Errorf("inappropriate listen history ranges, userA: %d , userB: %d", len(listenHistoryA), len(listenHistoryB))
	}
	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.8, s.extractPlayCount(listenHistoryA), s.extractPlayCount(listenHistoryB))
	return blendNumber, nil
}

func (s *BlendService) downloadTopAlbums(context context.Context, userName string, timeDuration blendTimeDuration) (map[string]CatalogueStats, error) {
	albumToPlays := make(map[string]CatalogueStats, 50)
	topAlbums, err := s.LastFMExternal.GetUserTopAlbums(
		context,
		userName,
		string(timeDuration),
		3,
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

		imageURL := s.getCatalogueImageURL(v.LFMImages) //Selects a good pic out of the ones given

		catStat := CatalogueStats{
			Artist:      v.Artist,
			Count:       playcount,
			PlatformURL: v.URL,
			Image:       imageURL,
			PlatformID:  v.MBID,
		}
		albumToPlays[v.Name] = catStat

	}

	return albumToPlays, nil
}

// ========== Track Blend ==========

func (s *BlendService) getTrackBlend(context context.Context, userA, userB userid, timeDuration blendTimeDuration) (int, error) {

	listenHistoryA, err := s.getTopX(context, userA, timeDuration, BlendCategoryTrack)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userA, err)
	}
	listenHistoryB, err := s.getTopX(context, userB, timeDuration, BlendCategoryTrack)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve top artists for %s as it returned error: %w", userB, err)
	}

	if len(listenHistoryA) == 0 || len(listenHistoryB) == 0 {
		return 0, fmt.Errorf("inappropriate listen history ranges, userA: %d , userB: %d", len(listenHistoryA), len(listenHistoryB))
	}
	//Using Log Weighted Cosine Similarity
	blendNumber := CalculateLWCS(0.8, s.extractPlayCount(listenHistoryA), s.extractPlayCount(listenHistoryB))
	return blendNumber, nil
}

func (s *BlendService) downloadTopTracks(context context.Context, userName string, timeDuration blendTimeDuration) (map[string]CatalogueStats, error) {
	trackToPlays := make(map[string]CatalogueStats)

	// TIMEDURATION DOESNT WORK
	// ------------------------
	topTracks, err := s.LastFMExternal.GetUserTopTracks(
		context,
		userName,
		string(timeDuration),
		5,
		50,
	)

	if len(topTracks.TopTracks.Track) == 0 {
		return trackToPlays, fmt.Errorf("downloaded empty toptrack, %w", err)
	}

	if err != nil {
		return trackToPlays, fmt.Errorf("could not extract TopTracks object from lastfm adapter, %w", err)
	}
	// for _, v := range topTracks.TopTracks.Track {
	// 	playcount, err := strconv.Atoi(v.Playcount)
	// 	if err != nil {
	// 		return trackToPlays, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
	// 	}
	// 	trackToPlays[v.Name] = playcount
	// }

	for _, v := range topTracks.TopTracks.Track {
		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return trackToPlays, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}

		imageURL := s.getCatalogueImageURL(v.LFMImages) //Selects a good pic out of the ones given

		catStat := CatalogueStats{
			Artist:      v.Artist,
			Count:       playcount,
			PlatformURL: v.URL,
			Image:       imageURL,
			PlatformID:  v.MBID,
		}
		trackToPlays[v.Name] = catStat

	}

	return trackToPlays, nil
}

func (s *BlendService) getCatalogueImageURL(images []musicapi.LFMImage) string {
	for _, img := range images {
		if img.Size == "large" {
			return img.URL
		}
	}
	//If no large image found, return the first available image
	if len(images) > 0 {
		return images[0].URL
	}
	return ""
}
