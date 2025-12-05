package blend

import (
	"backend-lastfm/internal/utility"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const LFM_EXPIRY = time.Duration(time.Hour * 24 * 3) //Three days //TODO: Change this to env var
const INVITE_EXPIRY = time.Duration(time.Hour * 24)

type RedisStateStore struct {
	client      *redis.Client
	userPrefix  string
	musicPrefix string
	blendPrefix string
}

func (r *RedisStateStore) AddUsersToBlend(context context.Context, id blendId, userids []userid) error {
	key := fmt.Sprintf("%s:%s:%s", r.blendPrefix, id, "users")

	members := make([]interface{}, len(userids))
	for i, u := range userids {
		members[i] = string(u)
	}

	return r.client.SAdd(context, key, members...).Err()
}

func (r *RedisStateStore) GetUsersFromBlend(context context.Context, id blendId) ([]userid, error) {
	key := fmt.Sprintf("%s:%s:%s", r.blendPrefix, id, "users")
	res, err := r.client.SMembers(context, key).Result()
	if err != nil {
		return nil, fmt.Errorf(" could not get LRange of users for users from blend id %s: and err %w", id, err)
	}
	users := make([]userid, len(res))
	for i, v := range res {
		users[i] = userid(v)
	}
	return users, nil
}

func (r *RedisStateStore) IsUserInBlend(context context.Context, user userid, id blendId) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", r.blendPrefix, id, "users")
	res, err := r.client.SIsMember(context, key, string(user)).Result()
	if err != nil {
		return false, fmt.Errorf(" error during checking if member was in set, as got value: %b: %w", res, err)
	}

	return res, nil
}

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client:      client,
		userPrefix:  "user", //TODO is this the right way to connect to redis?
		musicPrefix: "music_data",
		blendPrefix: "blend_data",
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

func (r *RedisStateStore) CacheUserMusicData(context context.Context, resp response) error {
	key := fmt.Sprintf("%s:%s:%s:%s", r.musicPrefix, resp.user, categoryPrefix[resp.category], durationPrefix[resp.duration])

	jsonBytes, err := utility.MapToJSON(resp.chart)
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
