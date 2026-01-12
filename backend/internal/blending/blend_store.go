package blend

import (
	musicapi "backend-lastfm/internal/music_api/lastfm"
	"backend-lastfm/internal/utility"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/redis/go-redis/v9"
)

const LFM_EXPIRY = time.Duration(time.Hour * 24 * 3) //Three days //TODO: Change this to env var
const INVITE_EXPIRY = time.Duration(time.Hour * 24)

type RedisStateStore struct {
	client           *redis.Client
	userPrefix       string
	lfmPrefix        string
	musicPrefix      string
	blendPrefix      string
	blendIndexPrefix string
}

func (r *RedisStateStore) DeleteMusicData(context context.Context, user string) error {
	pattern := fmt.Sprintf("%s:%s:*", r.musicPrefix, user)
	// r.client.Unlink(context, )
	keys, err := r.client.Keys(context, pattern).Result()
	if err != nil {
		return fmt.Errorf("could not get keys for deletion: %w", err)
	}
	if len(keys) == 0 {
		//No music_data to delete
		return nil
	}

	err = r.client.Unlink(context, keys...).Err()
	if err != nil {
		return fmt.Errorf("could not unlink keys during deleting music data: %w", err)
	}
	return nil
}

func (r *RedisStateStore) DeleteBlendByBlendId(context context.Context, user userid, blendId blendId) error {
	keyByUser := fmt.Sprintf("%s:%s:%s", r.userPrefix, "blends", string(user))
	keyByBlend := fmt.Sprintf("%s:%s:%s", r.blendPrefix, string(blendId), "users")

	pipe := r.client.TxPipeline()
	pipe.SRem(context, keyByUser, string(blendId)).Err()
	pipe.Del(context, keyByBlend).Err()

	_, err := pipe.Exec(context)
	return err
}

// Returns nil, nil for cache miss, else MapCatStats, nil or nil, error for error
func (r *RedisStateStore) GetFromCacheTopX(context context.Context, userName string, timeDuration blendTimeDuration, category blendCategory) (map[string]CatalogueStats, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", r.musicPrefix, userName, categoryPrefix[category], durationPrefix[timeDuration])

	Result, err := r.client.Get(context, key).Result()
	if err == redis.Nil {
		glog.Infof("Cache Miss: %s - %s", timeDuration, category)
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf(" during extracting cache from db, could not get json map from db:%w", err)
	}

	// respMap, err := utility.JSONToMap([]byte(Result))
	// if err != nil {
	// 	return nil, fmt.Errorf(" during extracting cache db, error in decoding from json: %w", err)
	// }

	respMap, err := musicapi.JSONToMapCatStats([]byte(Result))
	if err != nil {

		return nil, fmt.Errorf(" during extracting cache db, error in decoding from json: %w", err)
	}
	glog.Infof("Cache Hit: %s - %s", timeDuration, category)

	return respMap, nil

}

func (r *RedisStateStore) GetLFMByUserId(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("%s:%s", r.userPrefix, userID)
	result, err := r.client.HGet(ctx, key, "LFM Username").Result()
	return result, err
}

func (r *RedisStateStore) GetUserIdByLFMId(ctx context.Context, lfmuserid string) (string, error) {
	key := fmt.Sprintf("%s:%s", r.lfmPrefix, lfmuserid)
	result, err := r.client.Get(ctx, key).Result()
	return result, err
}

func (r *RedisStateStore) GetCachedOverallBlend(context context.Context, blendid blendId) (int, error) {
	key := fmt.Sprintf("%s:%s", r.blendPrefix, string(blendid))
	res, err := r.client.HGet(context, key, "Overall").Result()
	if err != nil {
		return -1, fmt.Errorf(" could not set overallblend num to blend: %w", err)
	}
	num, err := strconv.Atoi(res)
	if err != nil {
		return -1, fmt.Errorf(" could not convert cache value to num: %w", err)
	}
	return num, nil
}

func (r *RedisStateStore) GetBlendTimeStamp(context context.Context, id blendId) (time.Time, error) {
	key_timestamp := fmt.Sprintf("%s:%s", r.blendPrefix, string(id))
	res, err := r.client.HGet(context, key_timestamp, "Created At").Result()
	if err != nil {
		return time.Time{}, fmt.Errorf(" could not get blend timestamp of blend: %s", err)
	}
	num, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf(" could not convert time value to num: %w", err)
	}
	return time.Unix(num, 0), nil

}

func (r *RedisStateStore) AssignOverallBlendToBlend(context context.Context, id blendId, blendNum int) error {
	key := fmt.Sprintf("%s:%s", r.blendPrefix, string(id))
	// key_overall := fmt.Sprintf("%s:%s:%s", r.blendPrefix, string(id), "Overall")
	// key_timestamp := fmt.Sprintf("%s:%s:%s", r.blendPrefix, string(id), "Created At")

	// err := r.client.Set(context, key, blendNum, 0).Err()
	err := r.client.HSet(context, key, "Overall", blendNum, "Created At", time.Now().Unix()).Err()
	if err != nil {
		return fmt.Errorf(" could not set overallblend num to blend, with blendNum %d and Created at %d : %s", blendNum, time.Now().Unix(), err)
	}
	return nil
}
func (r *RedisStateStore) AddUsersToBlend(context context.Context, id blendId, userids []userid) error {

	pipe := r.client.TxPipeline() //Execute redis commands with atomicity

	key := fmt.Sprintf("%s:%s:%s", r.blendPrefix, id, "users")
	members := make([]interface{}, len(userids))
	for i, u := range userids {
		members[i] = string(u)

		//For secondary indexing= userId -> blendId
		s_index_key := fmt.Sprintf("%s:%s:%s", r.userPrefix, "blends", string(u))
		// pipe.ZAdd(context, s_index_key, redis.Z{
		// 	Score:  0.0,
		// 	Member: string(id),
		// })
		pipe.SAdd(context, s_index_key, string(id))
	}
	pipe.SAdd(context, key, members...).Err()

	_, err := pipe.Exec(context)

	return err
}

func (r *RedisStateStore) GetBlendsByUser(context context.Context, user userid) ([]blendId, error) {
	key := fmt.Sprintf("%s:%s:%s", r.userPrefix, "blends", string(user))
	ress, err := r.client.SMembers(context, key).Result()
	// ress, err := r.client.ZRange(context, key, -1, 999).Result()
	if err != nil {
		return nil, fmt.Errorf(" could not get Blends of user from user id %s: and err %w", user, err)
	}

	if len(ress) == 0 {
		return nil, nil //Empty?
	}

	blends := make([]blendId, len(ress))
	for i, res := range ress {
		blends[i] = blendId(res)
	}

	return blends, nil
}

func (r *RedisStateStore) GetUsersFromBlend(context context.Context, id blendId) ([]userid, error) {
	key := fmt.Sprintf("%s:%s:%s", r.blendPrefix, id, "users")
	res, err := r.client.SMembers(context, key).Result()
	if err != nil {
		return nil, fmt.Errorf(" could not get Members of users for users from blend id %s: and err %w", id, err)
	}
	if len(res) != 0 {
		users := make([]userid, len(res))
		for i, v := range res {
			users[i] = userid(v)
		}
		return users, nil
	} else {
		return nil, nil
	}

}

func (r *RedisStateStore) IsUserInBlend(context context.Context, user userid, id blendId) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", r.blendPrefix, id, "users")
	res, err := r.client.SIsMember(context, key, string(user)).Result()
	if err != nil {
		return false, fmt.Errorf(" error during checking if member was in set, as got value: %t: %w", res, err)
	}

	return res, nil
}

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client:           client,
		userPrefix:       "user", //TODO is this the right way to connect to redis?
		lfmPrefix:        "lfm_users",
		musicPrefix:      "music_data",
		blendPrefix:      "blend_data",
		blendIndexPrefix: "blend_data:index:",
	}

}

func (r *RedisStateStore) IsExistingBlendFromLink(context context.Context, linkValue blendLinkValue) (blendId, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", r.blendPrefix, "invite", linkValue, "id")
	res, err := r.client.Get(context, key).Result()
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf(" could not fetch blend's id from link in redis: %w", err)
	} else if err == redis.Nil {
		return "", nil
	} else {
		return blendId(res), nil
	}
}

func (r *RedisStateStore) GetLinkCreator(context context.Context, linkValue blendLinkValue) (userid, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", r.blendPrefix, "invite", linkValue, "creator")
	res, err := r.client.Get(context, key).Result()
	if err != nil {
		return "", fmt.Errorf(" could not fetch blend's user from link in redis: %w", err)
	} else {
		return userid(res), nil
	}
}

type Key struct {
	cat blendCategory
	dur blendTimeDuration
	// Expired bool
}

// Checks to see if any music data under this user exists. Returns true if anything exists
func (r *RedisStateStore) UserHasAnyMusicData(context context.Context, user userid) (bool, error) {
	pattern := fmt.Sprintf("%s:%s:*", r.musicPrefix, user)
	iter := r.client.Scan(context, 0, pattern, 1).Iterator()
	for iter.Next(context) {
		return true, nil
	}

	err := iter.Err()
	if err != nil {
		return false, fmt.Errorf(" error during checking full cache expiry: %w", err)
	}

	return false, nil

}

// Individually checks for each possible key in cache and returns the ones that are expired
func (r *RedisStateStore) GetEachExpiredCacheEntryByUser(context context.Context, user userid) ([]Key, error) {
	Keys := make([]Key, 0)
	for _, v1 := range categoryRange {
		for _, v2 := range durationRange {
			key := fmt.Sprintf("%s:%s:%s:%s", r.musicPrefix, user, v1, v2)
			err := r.client.Get(context, key).Err()
			if err == redis.Nil {
				Keys = append(Keys, Key{
					cat: v1,
					dur: v2,
					// Expired: true,
				})

			} else {
				return Keys, fmt.Errorf(" non-nil error during checking for data entry in cache:%w", err)
			}
		}
	}
	return Keys, nil
}

var categoryPrefix = map[blendCategory]string{
	BlendCategoryAlbum:  "album",
	BlendCategoryArtist: "artist",
	BlendCategoryTrack:  "track",
}

var durationPrefix = map[blendTimeDuration]string{
	BlendTimeDurationOneMonth:   "one_month",
	BlendTimeDurationThreeMonth: "three_month",
	BlendTimeDurationYear:       "one_year",
}

func (r *RedisStateStore) CacheUserMusicData(context context.Context, resp complexResponse) error {
	key := fmt.Sprintf("%s:%s:%s:%s", r.musicPrefix, resp.user, categoryPrefix[resp.category], durationPrefix[resp.duration])

	jsonBytes, err := utility.ObjectToJSON(resp.data)
	if err != nil {
		return fmt.Errorf(" during caching to db, error in encoding to json: %w", err)
	}
	err = r.client.Set(context, key, jsonBytes, LFM_EXPIRY).Err()
	if err != nil {
		return fmt.Errorf(" during caching to db, could not set json map in db: %w", err)
	}
	return nil
}

func (r *RedisStateStore) GetUser(userA UUID) (string, error) {
	return "saflas", nil
}

type UserListenHistory struct {
	// Define fields for user listen history
}

func (r *RedisStateStore) SetUserToLink(context context.Context, userA userid, linkValue blendLinkValue) error {
	key := fmt.Sprintf("%s:%s:%s:%s", r.blendPrefix, "invite", linkValue, "creator")
	err := r.client.Set(context, key, string(userA), INVITE_EXPIRY).Err()
	if err != nil {
		return fmt.Errorf(" could not set blend's user from link into redis: %w", err)
	} else {
		return nil
	}

}

func (r *RedisStateStore) GetUserListenHistory(userID string) (UserListenHistory, error) {
	return UserListenHistory{}, nil
}
