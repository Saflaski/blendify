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
	client *redis.Client
	prefix string
}

const allMusicPrefix = "music_data"

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
	key := fmt.Sprintf("%s:%s:%s:%s", allMusicPrefix, resp.user, categoryPrefix[resp.category], durationPrefix[resp.duration])

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
		client: client,
		prefix: "user:", //TODO is this the right way to connect to redis?
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
