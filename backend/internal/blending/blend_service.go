package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"backend-lastfm/internal/musicbrainz"
	"cmp"
	"context"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type BlendService struct {
	repo           *BlendStore
	LastFMExternal *musicapi.LastFMAPIExternal
	MBService      *musicbrainz.MBService
}

func (s *BlendService) GetPermanentLinkForUser(context context.Context, userA userid) (permaLinkValue, error) {

	//Check if the user has a current permanent link
	link, err := s.repo.GetPermanentLinkByUser(context, userA)
	if err != nil {
		return "", fmt.Errorf(" could not get permanent link from userid %s due to repo error: %w", userA, err)
	}

	if link != "" {
		//Return existing link
		return link, nil
	} else {
		//Make new permanent link
		linkID, err := gonanoid.New(10)

		if err != nil {
			return "", fmt.Errorf(" could not generate new permanent link from userid %s due to nanoid error: %w", userA, err)
		}
		newLinkValue := permaLinkValue(linkID)
		err = s.repo.AssignPermanentLinkToUser(context, userA, newLinkValue)
		if err != nil {
			return "", fmt.Errorf(" could not assign new permanent link from userid %s due to repo error: %w", userA, err)
		}
		return newLinkValue, nil
	}
}

func (s *BlendService) GetBlendTopGenres(context context.Context, blendId blendId, userA userid, timeDuration blendTimeDuration) ([]string, error) {
	userids, err := s.repo.GetUsersFromBlend(context, blendId)
	if err != nil {
		return nil, fmt.Errorf(" error getting users from blend id: %s, err: %w", blendId, err)
	}

	//Collect top genres from all users
	allUserGenres := make([][]string, len(userids))

	for k, userid := range userids {
		userTopGenres, err := s.GetUserTopGenres(context, userid)
		if err != nil {
			return nil, fmt.Errorf(" could not get top genres from userid %s : %w", userid, err)
		}
		allUserGenres[k] = userTopGenres
	}

	//This is passed off to another function as in the future we can use a separate topgenre calculation method
	commonTopGenres, err := s.CalculateIntersectionOfStringSlices(allUserGenres)
	if err != nil {
		return nil, fmt.Errorf(" coult not calculate intersection of genres: %v", err)
	}

	return commonTopGenres, nil
}

func (s *BlendService) CalculateIntersectionOfStringSlices(megaslice [][]string) ([]string, error) {

	//I did not have a better word than megaslice ok?
	if len(megaslice) == 0 {
		return []string{}, fmt.Errorf(" Given length 0 for intersection calculation")
	}
	// "Rock", "Pop", "Jazz", "Classical", "Hip Hop", "Electronic"
	// "Rock", "Pop", "Country", "Classical", "Reggae", "Electronic"
	counts := make(map[string]int)
	for _, subSlice := range megaslice {
		visited := make(map[string]bool)
		for _, item := range subSlice {
			if !visited[item] {
				counts[item]++
				visited[item] = true
			}
		}
	}

	intersection := make([]string, 0)
	totalSlices := len(megaslice)

	seen := make(map[string]bool)
	for _, item := range megaslice[0] {
		if !seen[item] && counts[item] == totalSlices {
			intersection = append(intersection, item)
			seen[item] = true
		}
	}

	return intersection, nil
}

func (s *BlendService) GetUserTopGenres(context context.Context, user userid) ([]string, error) {
	topGenres, err := s.repo.GetCachedUserTopGenres(context, user)
	if err != nil {
		return nil, fmt.Errorf(" could not get top genres from userid %s : %w", user, err)
	}

	return topGenres, nil
}

func (s *BlendService) GetUserInfo(context context.Context, userid userid) (any, error) {
	username, err := s.repo.GetLFMByUserId(context, string(userid))
	if err != nil {
		return nil, fmt.Errorf(" could not get username from userid %s : %w", userid, err)
	}

	mapToReturn := make(map[string]string)
	mapToReturn["username"] = username
	userinfo, err := s.LastFMExternal.GetUserInfo(context, username)
	if err != nil {
		return nil, fmt.Errorf(" could not get userinfo from username %s : %w", username, err)
	}
	mapToReturn["playcount"] = userinfo.User.Playcount
	mapToReturn["artist"] = userinfo.User.Artist_count
	mapToReturn["track"] = userinfo.User.Track_count

	return mapToReturn, nil
}

func (s *BlendService) GetUserTopItems(context context.Context, blendId blendId, user userid, requestedUsername string, mode blendCategory, duration blendTimeDuration) (TopItems, error) {

	ok, err := s.AuthoriseBlend(context, blendId, user)
	if err != nil {
		return TopItems{}, fmt.Errorf("error during getting topitems' authorising blend: %w", err)
	}
	if !ok {
		return TopItems{}, fmt.Errorf(" user %s not authorised to access blend %s", user, blendId)
	}
	useridRequested, err := s.repo.GetUserIdByLFMId(context, requestedUsername)
	if err != nil {
		return TopItems{}, fmt.Errorf("error during getting usertopitem's useridbylfm: %s : %w", requestedUsername, err)
	}
	items, err := s.getTopX(context, userid(useridRequested), duration, mode)
	if err != nil {
		return TopItems{}, fmt.Errorf("could not get top x during GetUserTopItems: %w", err)
	}

	keys := s.GetSortedEntries(items)

	return TopItems{
		Items: keys[0:10],
	}, nil
}

func (s *BlendService) DeleteUserBlends(context context.Context, user string) error {
	blendIds, err := s.repo.GetBlendsByUser(context, userid(user))
	if err != nil {
		glog.Info("USER DELETE REQUEST FAILED")
		return fmt.Errorf(" could not find blends from userid %s : %w", user, err)
	}
	for _, v := range blendIds {
		err = s.DeleteBlend(context, userid(user), v)
		if err != nil {
			glog.Info("USER DELETE REQUEST FAILED")
			return fmt.Errorf(" could not delete blend %s for user: %s -- %w ", v, user, err)
		}
	}
	err = s.repo.DeleteMusicData(context, user)
	if err != nil {
		glog.Info("USER DELETE REQUEST FAILED")
		return fmt.Errorf("could not delete music data for user %s: %w", user, err)
	}
	return nil

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
			return Blends{}, fmt.Errorf(" could not find blendusers from blendid %s : %v", v, err)
		}
		blendPlatformUsernames := make([]platformid, len(blendUsers))
		for j, v_2 := range blendUsers {
			platformUser, err := s.repo.GetLFMByUserId(context, string(v_2))
			if err != nil {
				return Blends{}, fmt.Errorf(" could not extract platformid from userid %s for blendId: %s with requesting user: %s: %v", v_2, v, user, err)
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

func (s *BlendService) GetBlendAndRefreshCache(context context.Context, blendId blendId) (DuoBlend, error) {
	duoBlend, err := s.GetDuoBlendData(context, blendId)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" could not get duoblend data: %w", err)
	}

	// //Refresh cache after getting data
	// err = s.RefreshOverallBlendInCache(context, blendId)
	// if err != nil {
	// 	glog.Infof(" could not refresh overall blend in cache: %w", err)
	// 	//Not a fatal error
	// }

	return duoBlend, nil
}
func (s *BlendService) GetDuoBlendData(context context.Context, blendId blendId) (DuoBlend, error) {

	//Get json data for percentage data of 3x3 data
	//Get json data for percentage data of distribution of art/alb/tra

	userids, err := s.repo.GetUsersFromBlend(context, blendId)
	if err != nil {
		return DuoBlend{}, fmt.Errorf(" error getting users from blend id: %s, err: %w", blendId, err)
	}
	if len(userids) == 0 {
		return DuoBlend{}, nil
	}

	// err = s.PopulateUsersByBlend(context, blendId)
	// if err != nil {
	// 	return DuoBlend{}, fmt.Errorf(" Could not populate user data: %w", err)
	// }

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

	err = s.repo.AssignOverallBlendToBlend(context, blendId, duoBlend.OverallBlendNum)
	if err != nil {
		glog.Infof(" could not assign duoblend integer value during adding users to blend:%s: %w", blendId, err)
		//This is not a fatal error and we can live without it happening but still need to log
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
		Name:     name,
		ImageURL: aStat.Image,
		URL:      aStat.PlatformURL,
		// MBID:           aStat.MBID,
		ArtistName:     aStat.Artist.Name,
		ArtistURL:      aStat.Artist.URL,
		ArtistImageURL: s.getCatalogueImageURL(aStat.Artist.LFMImages),
		Playcounts:     []int{countA, countB},
	}

	if aStat.Genres != nil {
		entry.Genres = aStat.Genres
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

func (s *BlendService) GetSortedEntries(entries map[string]CatalogueStats) []string {
	keys := slices.Collect(maps.Keys(entries))

	slices.SortFunc(keys, func(a, b string) int {
		if entries[a].Count != entries[b].Count {
			return cmp.Compare(entries[b].Count, entries[a].Count)
		}
		return cmp.Compare(a, b)
	})

	return keys
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
		10, 7, 5)

	if err != nil {
		return 0, fmt.Errorf(" could not calc over artist: %w", err)
	}
	albumOverall, err := combineNumbersWithWeights(
		albumBlend.OneMonth,
		albumBlend.ThreeMonth,
		albumBlend.OneYear,
		10, 7, 5)
	if err != nil {
		return 0, fmt.Errorf(" could not calc over artist: %w", err)
	}

	trackOverall, err := combineNumbersWithWeights(
		trackBlend.OneMonth,
		trackBlend.ThreeMonth,
		trackBlend.OneYear,
		10, 7, 5)
	if err != nil {
		return 0, fmt.Errorf(" could not calc over tracks: %w", err)
	}

	overallBlend, err := combineNumbersWithWeights(artistOverall, albumOverall, trackOverall, 10, 6, 8)
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

func (s *BlendService) MakeBlendFromPermaLink(context context.Context, userA userid, link permaLinkValue) (blendId, error) {

	//First get blend link from permanent link
	userB, err := s.repo.GetUserByPermanentLink(context, link) //Fetch user who created link
	if err != nil {
		return "", fmt.Errorf(" error during getting user (creator) from link : %w", err)
	}
	glog.Infof("Blend created by: %s", userB)

	//Safety net to make sure userA != userB
	if userB == userA {
		glog.Info("Same user nvm")
		return "0", nil //0 is code for consuming user being the same user as creating user
	}

	id, err := s.GenerateNewBlendId(context, []userid{userA, userB}) //Should this make the whole blend or?
	if err != nil {
		return "", fmt.Errorf(" error during making a blend with users %s and %s: %w", userA, userB, err)
	}
	glog.Infof("Generating new blend: %s", id)

	return id, nil
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
	if id == "" { //It's not an existing blend. Go forth with making a new blend
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
		// userids, err := s.repo.GetUsersFromBlend(context, id)
		// if err != nil {
		// 	return "", fmt.Errorf(" error getting users from blend id: %s, err: %w", id, err)
		// }
		// if len(userids)+1 > BLEND_USER_LIMIT {
		// 	glog.Info("Tried to add too many users to blend")
		// 	return "-1", nil
		// }

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
			// glog.Infof("User does not exist in blend. Adding %s", userA)
			return id, nil
		} else {
			// glog.Infof("User already exists in blend")
			//Nothing to see here, just return the existing blend id
			return id, nil
		}
	}

}

func (s *BlendService) AddUsersToBlend(context context.Context, id blendId, userids []userid) error {
	// TEMPORARY LIMIT FOR NUM USERS WHO CAN BE IN A BLEND
	blendusers, err := s.repo.GetUsersFromBlend(context, id)
	if err != nil {
		return fmt.Errorf(" error getting users from blend id: %s, err: %w", id, err)
	}
	if len(blendusers)+len(userids) > BLEND_USER_LIMIT {
		glog.Info("Tried to add too many users to blend")
		return fmt.Errorf(" adding too many users to blend id: %s with new userids %s , err: %w", id, userids, err)
	}

	err = s.repo.AddUsersToBlend(context, id, userids)
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

				//It is possible that this request for new data came as a result of a mid-fresh load reload.
				//Therefore we use the cached data to reduce load up times on the n+1th load.
				dbResp, err := s.repo.GetFromCacheTopX(ctx, string(user), d, c)
				if err != nil {
					glog.Errorf(" Cache error during getting topX for user %s with duration %s and category %s that needs to be checked: %w", user, d, c, err)
				}
				resp := complexResponse{}
				if len(dbResp) != 0 || dbResp != nil { //Use cached data
					fmt.Println("Using cached data for user:", user, " duration:", d, " category:", c)
					resp = complexResponse{
						user:     user,
						data:     dbResp,
						duration: d,
						category: c,
						err:      err,
					}

				} else {
					respData, err := s.downloadTopX(ctx, platformUsername, d, c)
					if len(respData) == 0 {
						//wrap the error regardless of it is nil/not nil
						err = fmt.Errorf(" downloaded empty map from platform: %w", err)
					}
					resp = complexResponse{
						user:     user,
						data:     respData,
						duration: d,
						category: c,
						err:      err,
					}
					fmt.Println("Caching data for user:", resp.user, " duration:", resp.duration, " category:", resp.category)
					err = s.cacheLFMData(ctx, user, resp.category, resp.duration, resp.data)
					if err != nil {
						resp.err = fmt.Errorf(" error during caching lfm data: %w", err)
					}
				}
				// respData, err := s.downloadTopX(ctx, platformUsername, d, c)

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

		if resp.category == BlendCategoryTrack {
			topGenres := s.extractTopGenres(resp.data, 50)
			// if err := s.cacheTopGenres(ctx, resp.user, topGenres); err != nil {
			// 	return fmt.Errorf(" could not cache top genres: %v", err)
			// }
			// for _, genre := range topGenres {
			// 	glog.Infof("User: %s Top Genre: %s", resp.user, genre)
			// }
			cacheTopGenresErr := s.cacheTopGenres(ctx, resp.user, resp.data, topGenres)
			if cacheTopGenresErr != nil {
				return fmt.Errorf(" could not cache top genres: %v", cacheTopGenresErr)
			}
			//DEBUG
			topGenres, err := s.repo.GetCachedUserTopGenres(ctx, resp.user)
			if err != nil {
				return fmt.Errorf(" could not extract top genres from redis cache: %w", err)
			}
			// for _, genre := range topGenres {
			// 	glog.Infof("User: %s Top Genre from REDIS: %s", resp.user, genre)
			// }
		}
	}
	return nil
}

func (s *BlendService) GetCachedUserTopGenres(ctx context.Context, user userid) ([]string, error) {
	topGenres, err := s.repo.GetCachedUserTopGenres(ctx, user)
	if err != nil {
		return nil, fmt.Errorf(" could not extract top genres from cache: %w", err)
	}
	return topGenres, nil
}

func (s *BlendService) cacheTopGenres(ctx context.Context, userid userid, mcs map[string]CatalogueStats, topGenres []string) error {
	// For each genre in the top genres, cache it
	// And add those mbids to the genre
	s.repo.CacheUserTopGenres(ctx, userid, mcs, topGenres)
	return nil
}

func (s *BlendService) extractTopGenres(input map[string]CatalogueStats, topN int) []string {
	genreCount := make(map[string]int)
	for _, v := range input {
		for _, genre := range v.Genres {
			genreCount[genre] += v.Count
		}
	}

	type genrePair struct {
		Genre string
		Count int
	}

	genrePairs := make([]genrePair, 0, len(genreCount))
	for genre, count := range genreCount {
		genrePairs = append(genrePairs, genrePair{Genre: genre, Count: count})
	}

	sort.Slice(genrePairs, func(i, j int) bool {
		return genrePairs[i].Count > genrePairs[j].Count
	})

	topGenres := make([]string, 0, topN)
	for i := 0; i < topN && i < len(genrePairs); i++ {
		topGenres = append(topGenres, genrePairs[i].Genre)
	}

	return topGenres
}

func (s *BlendService) cacheLFMData(ctx context.Context, user userid, category blendCategory, duration blendTimeDuration, data map[string]CatalogueStats) error {

	cacheTime := time.Duration(time.Hour * 24 * 1) //1 day default
	switch duration {
	case BlendTimeDurationOneMonth:
		cacheTime *= 2
	case BlendTimeDurationThreeMonth:
		cacheTime *= 4
	case BlendTimeDurationYear:
		cacheTime *= 5
	}

	glog.Info("Caching data: \n", " duration:", duration, " category:", category, " with cache time (hours):", cacheTime.Hours())

	err := s.repo.CacheUserMusicDataV2(
		ctx,
		user,
		category,
		duration,
		data,
		cacheTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *BlendService) GenerateNewLinkAndAssignToUser(context context.Context, userA userid) (blendLinkValue, error) {

	//Generate a linkId to be returned that won't hash collide
	newInviteValue := blendLinkValue(uuid.New().String())
	err := s.repo.SetUserToLink(context, userA, newInviteValue)
	if err != nil {
		return "", fmt.Errorf(" could not set user to link: %w", err)
	}

	return newInviteValue, nil
}

func NewBlendService(blendStore BlendStore, lfmAdapter musicapi.LastFMAPIExternal, mbService musicbrainz.MBService) *BlendService {
	return &BlendService{&blendStore, &lfmAdapter, &mbService}
}

func (s *BlendService) GetBlend(context context.Context, userA userid, userB userid, category blendCategory, timeDuration blendTimeDuration) (int, error) {
	//Implement logic to calculate blend percentage based on user data, category, and time duration

	// //Get the username from the UUID of the given user that's sending the request
	// userNameA, err := s.repo.GetUser(userA)
	// if err != nil {
	// 	return 0, fmt.Errorf("could not extract username from UUID of user with ID: %s, %w", userA, err)
	// }

	// glog.Info("Calculating blend for users: ", userA, " + ", userB, " category: ",
	// category, " timeDuration: ", timeDuration)

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
	blendNumber := CalculateLWCS(1.0, s.extractPlayCount(listenHistoryA), s.extractPlayCount(listenHistoryB))
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
		// glog.Info("Cache miss on user %s", userid)
		platformUsername, err := s.getLFM(context, string(userid))
		if err != nil {
			return nil, fmt.Errorf(" could get platform username from given userid:%s with err: %w", userid, err)
		}
		glog.Info("Cache miss, so returning downloadTopX")
		lfmResp, err := s.downloadTopX(context, platformUsername, timeDuration, category)
		if err != nil {
			return nil, fmt.Errorf(" did not download %s %s properly: %w", timeDuration, category, err)
		}

		//Cache the downloaded data
		cacheErr := s.cacheLFMData(context, userid, category, timeDuration, lfmResp)
		if cacheErr != nil {
			return nil, fmt.Errorf(" error during caching lfm data: %w", cacheErr)
		}

		if len(lfmResp) == 0 {
			return nil, fmt.Errorf(" downloaded empty map for %s %s : %s", timeDuration, category, err)
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

	artistCatalogues := make(map[string]CatalogueStats)
	// topArtist, err := s.LastFMExternal.GetUserTopArtists(
	// 	context,
	// 	userName,
	// 	string(timeDuration),
	// 	4,
	// 	50,
	// )

	topArtist, err := s.LastFMExternal.GetUserTopArtistsAsync(
		context,
		userName,
		string(timeDuration),
		4,
		50,
	)

	if err != nil {
		return artistCatalogues, fmt.Errorf("could not extract TopArtists object from lastfm adapter, %w", err)
	}
	// for _, v := range topArtist.TopArtists.Artist {
	// 	playcount, err := strconv.Atoi(v.Playcount)
	// 	if err != nil {
	// 		return artistToPlaybacks, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
	// 	}
	// 	artistToPlaybacks[v.Name] = playcount
	// }

	for _, v := range topArtist.TopArtists.Artist {
		if v.MBID == "" {
			//skip
			continue
		}

		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return artistCatalogues, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}

		imageURL := s.getCatalogueImageURL(v.LFMImages) //Selects a good pic out of the ones given

		catStat := CatalogueStats{
			Artist:      v,
			Count:       playcount,
			PlatformURL: v.URL,
			Image:       imageURL,
			PlatformID:  v.MBID,
		}
		artistCatalogues[v.Name] = catStat

	}

	artistCatalogues, err = s.PopulateArtistMBIDs(context, artistCatalogues)
	if err != nil {
		return artistCatalogues, fmt.Errorf(" could not populate artist mbids: %w", err)
	}
	artistCatalogues, err = s.PopulateGenresForMapCatStats(context, artistCatalogues, "artist")
	if err != nil {
		return artistCatalogues, fmt.Errorf(" could not populate artist genres: %w", err)
	}

	return artistCatalogues, nil
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
	blendNumber := CalculateLWCS(1.0, s.extractPlayCount(listenHistoryA), s.extractPlayCount(listenHistoryB))
	return blendNumber, nil
}

func (s *BlendService) downloadTopAlbums(context context.Context, userName string, timeDuration blendTimeDuration) (map[string]CatalogueStats, error) {
	albumToPlays := make(map[string]CatalogueStats, 50)
	// topAlbums, err := s.LastFMExternal.GetUserTopAlbums(
	// 	context,
	// 	userName,
	// 	string(timeDuration),
	// 	2,
	// 	50,
	// )
	topAlbums, err := s.LastFMExternal.GetUserTopAlbumsAsync(
		context,
		userName,
		string(timeDuration),
		2,
		50,
	)

	if err != nil {
		return albumToPlays, fmt.Errorf("could not extract TopAlbums object from lastfm adapter, %w", err)
	}
	for _, v := range topAlbums.TopAlbums.Album {

		if v.MBID == "" {
			//Try to get it from album name + artist name
			//skip for now
			continue
		}

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
	blendNumber := CalculateLWCS(1.0, s.extractPlayCount(listenHistoryA), s.extractPlayCount(listenHistoryB))
	return blendNumber, nil
}

func (s *BlendService) downloadTopTracks(context context.Context, userName string, timeDuration blendTimeDuration) (map[string]CatalogueStats, error) {
	trackToPlays := make(map[string]CatalogueStats)

	// TIMEDURATION DOESNT WORK
	// ------------------------
	// topTracks, err := s.LastFMExternal.GetUserTopTracks(
	// 	context,
	// 	userName,
	// 	string(timeDuration),
	// 	6,
	// 	50,
	// )

	topTracks, err := s.LastFMExternal.GetUserTopTracksAsync(
		context,
		userName,
		string(timeDuration),
		6,
		50,
	)

	if len(topTracks.TopTracks.Track) == 0 {
		return trackToPlays, fmt.Errorf("downloaded empty toptrack, %w", err)
	}

	if err != nil {
		return trackToPlays, fmt.Errorf("could not extract TopTracks object from lastfm adapter, %w", err)
	}

	for _, v := range topTracks.TopTracks.Track {
		// exists := false
		// if v.MBID != "" {
		// 	exists, err = s.MBService.IsValidMBIDForRecord(context, v.MBID)
		// 	if err != nil {
		// 		return trackToPlays, fmt.Errorf("could not check if record exists through mbid: %w", err)
		// 	}
		// 	if !exists {
		// 		glog.Infof("Record does not exist for mbid: %s, track: %s by artist: %s", v.MBID, v.Name, v.Artist.Name)
		// 	}
		// }
		// _ = exists
		// if v.MBID == "" {

		// continue
		// //Get MBID from artistname + track name, if not found, skip
		// trackInfo, err := s.MBService.GetMBIDFromArtistAndTrackName(context, v.Artist.Name, v.Name)
		// if err != nil {
		// 	return trackToPlays, fmt.Errorf("could not perform mbidfromartistandtrackname after got empty mbid: %w", err)
		// }
		// if trackInfo.IsEmpty() {
		// 	continue //Could not find a suitable track using pre-given SQL settings
		// }

		// v.MBID = trackInfo.RecordingMBID
		// v.Artist.MBID = trackInfo.ArtistMBID
		// v.Artist.Name = trackInfo.ArtistName
		// v.Name = trackInfo.RecordingName

		// }
		playcount, err := strconv.Atoi(v.Playcount)
		if err != nil {
			return trackToPlays, fmt.Errorf("got unparseable string during string -> int conversation: %w", err)
		}

		imageURL := s.getCatalogueImageURL(v.LFMImages) //Selects a good pic out of the ones given
		// genreObjects, err := s.MBService.GetGenresByRecordingMBID(context, v.MBID)
		// if err != nil {
		// 	return trackToPlays, fmt.Errorf("could not get genres during gettoptracks: %w", err)

		// }
		// genres := make([]string, len(genreObjects))
		// for i, g := range genreObjects {
		// 	genres[i] = g.Name //Capitalize first letter of each genre
		// }
		catStat := CatalogueStats{
			Artist:      v.Artist,
			Count:       playcount,
			PlatformURL: v.URL,
			Image:       imageURL,
			PlatformID:  v.MBID,
		}
		trackToPlays[v.Name] = catStat

	}

	trackToPlays, err = s.PopulateTrackMBIDs(context, trackToPlays)
	if err != nil {
		return trackToPlays, fmt.Errorf(" could not populate mbids for map cat stats: %w", err)
	}
	trackToPlays, err = s.PopulateGenresForMapCatStats(context, trackToPlays, "recording")
	if err != nil {
		return trackToPlays, fmt.Errorf(" could not internally populate genres for map cat stats : %w", err)
	}
	//Prepare list of tracknames and artist names from CatalogueStats map

	return trackToPlays, nil
}

func (s *BlendService) PopulateArtistMBIDs(context context.Context, input map[string]CatalogueStats) (map[string]CatalogueStats, error) {
	artistNames := make([]string, len(input))
	i := 0
	for k, _ := range input {
		artistNames[i] = k
		i++
	}

	mbidList, err := s.MBService.GetMBIDsFromArtistNames(context, artistNames)
	if err != nil {
		return input, fmt.Errorf(" could not get mbid list from artist: %w", err)
	}

	for _, v := range mbidList {
		newMapCatStats, ok := input[v.ArtistName]
		if !ok {
			continue
		}

		newMapCatStats.PlatformID = v.ArtistMBID
		newMapCatStats.Artist.MBID = v.ArtistMBID
		newMapCatStats.Artist.Name = v.ArtistName
		input[v.ArtistName] = newMapCatStats

	}

	if len(mbidList) == 0 {
		return input, fmt.Errorf(" could not populate any mbids for map cat stats")
	}

	return input, nil
}
func (s *BlendService) PopulateTrackMBIDs(context context.Context, input map[string]CatalogueStats) (map[string]CatalogueStats, error) {

	//Whilst this does not fully populate it as it uses '=' operator on SQL,
	//we can assume that the data is good enough for our use case
	//the goal would now be to start a cron job later that updates this data fully in the background
	//ie, get mbids from trigram search on postgresql
	//and for genres, make request to
	//Prepare list of tracknames and artist names from CatalogueStats map
	trackNames := make([]string, len(input))
	artistNames := make([]string, len(input))
	i := 0
	for k, v := range input {
		trackNames[i] = k
		artistNames[i] = v.Artist.Name
		i++
	}

	mbidList, err := s.MBService.GetMBIDsFromArtistAndTrackNames(context, artistNames, trackNames)
	if err != nil {
		return input, fmt.Errorf(" could not get mbid list from artist and track names: %w", err)
	}

	for _, v := range mbidList {
		newMapCatStats, ok := input[v.RecordingName]
		if !ok {
			continue
		}
		newMapCatStats.PlatformID = v.RecordingMBID
		newMapCatStats.Artist.MBID = v.ArtistMBID
		newMapCatStats.Artist.Name = v.ArtistName
		input[v.RecordingName] = newMapCatStats
	}

	// for k := 0; k < len(input); k++ {

	// 	trackName := trackNames[k]
	// 	catStat := input[trackName]
	// 	mbidInfo := mbidList[k]
	// 	if !mbidInfo.IsEmpty() {
	// 		catStat.PlatformID = mbidInfo.RecordingMBID
	// 		catStat.Artist.MBID = mbidInfo.ArtistMBID
	// 		catStat.Artist.Name = mbidInfo.ArtistName
	// 		// genreObjects := mbidInfo.Genres
	// 		// genres := make([]string, len(genreObjects))
	// 		// for i, g := range genreObjects {
	// 		// 	genres[i] = g.Name //Capitalize first letter of each genre
	// 		// }
	// 		// catStat.Genres = genres
	// 	}
	// 	output[trackName] = catStat
	// }

	if len(mbidList) == 0 {
		return input, fmt.Errorf(" could not populate any mbids for map cat stats")
	}

	return input, nil
}

// Need MBIDs to get genres from musicbrainz or else will fail
func (s *BlendService) PopulateGenresForMapCatStats(context context.Context, input map[string]CatalogueStats, mode string) (map[string]CatalogueStats, error) {
	output := make(map[string]CatalogueStats)

	//We are going to pass a list of mbids to musicbrainz service to get genres
	mbids := make([]string, 0, len(input))
	// injectionMap := make(map[string]string) //mbid -> genres
	for k, v := range input {
		if input[k].PlatformID == "" {
			continue
		} else {
			mbids = append(mbids, v.PlatformID)
		}
	}

	genreList := make(map[string][]musicbrainz.Genre)
	if mode == "recording" {
		recordingsGenreList, err := s.MBService.GetGenresByRecordingMBIDs(context, mbids)
		genreList = recordingsGenreList
		if err != nil {
			return output, fmt.Errorf(" could not get genres by recording mbids: %w", err)
		}
	} else if mode == "artist" {
		artistGenresList, err := s.MBService.GetGenreByArtistMBIDs(context, mbids)
		genreList = artistGenresList
		if err != nil {
			return output, fmt.Errorf(" could not get genres by recording mbids: %w", err)
		}
	}

	for k, v := range input {
		if v.PlatformID == "" {
			output[k] = v
			continue
		}
		genreObjects, ok := genreList[v.PlatformID]
		if !ok {
			output[k] = v
			continue
		}
		genres := make([]string, len(genreObjects))
		for i, g := range genreObjects {
			genres[i] = g.Name //Capitalize first letter of each genre
		}
		v.Genres = genres
		output[k] = v
	}
	return output, nil
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
