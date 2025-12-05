package blend

import (
	"backend-lastfm/internal/utility"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const LFM_EXPIRY = time.Duration(time.Hour * 24 * 3) //Three days //TODO: Change this to env var

type RedisStateStore struct {
	client      *redis.Client
	userPrefix  string
	musicPrefix string
}

func (r *RedisStateStore) GetUsersFromBlend(id blendId) ([]userid, error) {
	panic("unimplemented")
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

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client:      client,
		userPrefix:  "user", //TODO is this the right way to connect to redis?
		musicPrefix: "music_data",
	}
}
func (r *RedisStateStore) SetUserToLink(ctx context.Context, userA UUID, newInviteId uuid.UUID) {
	//userIDString := string(userA)
	//Ideally we set via user id and not straight user
	//user:USERID:
	// queryString := r.prefix + ""
	// r.client.HSet(ctx)

}

func (r *RedisStateStore) GetUserListenHistory(userID string) (UserListenHistory, error) {
	return UserListenHistory{}, nil
}
