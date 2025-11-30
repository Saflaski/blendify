package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthRepository interface {
	SetNewStateSid(ctx context.Context, stateToken, sessionID string, ttl time.Duration) error
	ConsumeStateSID(ctx context.Context, stateToken string) (string, error)
	SetSidWebSesssionKey(ctx context.Context, sessionID, webSessionKey string, ttl time.Duration) error
	GetSidKey(ctx context.Context, sessionID string) (string, error)
	DelSidKey(ctx context.Context, sessionID string) error
	MakeNewUser(context context.Context, validationSid string, userName string, userid uuid.UUID) error
	GetUserBySessionID(context context.Context, sid string) (string, error)
	DeleteUser(ctx context.Context, userid string) error
	DeleteSingularSessionID(context context.Context, sid string) error
}

type RedisStateStore struct {
	client            *redis.Client
	prefixAuth        string
	prefixUser        string
	prefixSidToUser   string
	prefixSidList     string
	sidExpirationTime int64
}

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client:            client,
		prefixAuth:        "login_state:", //TODO is this the right way to connect to redis?
		prefixUser:        "user",
		prefixSidToUser:   "user_sids",
		prefixSidList:     "sidlist",
		sidExpirationTime: time.Duration(time.Hour * 24).Nanoseconds(),
	}
}

func (r *RedisStateStore) MakeNewUser(ctx context.Context, validationSid string, userName string, userid uuid.UUID) error {

	key := fmt.Sprintf("%s:%s", r.prefixUser, userid.String())

	pipe := r.client.TxPipeline() //Execute redis commands with atomicity

	pipe.HSet(context.Background(), key, "lfm", userName)
	r.queueAddNewSid(pipe, userid.String(), validationSid)

	_, err := pipe.Exec(ctx)

	return err
}

// Adds newsid atomically.
func (r *RedisStateStore) queueAddNewSid(pipe redis.Pipeliner, userid string, sid string) error {

	// pipe := r.client.TxPipeline()

	err := r.queueAddNewSidUserIndexCmd(pipe, userid, sid)
	if err != nil {
		return err
	}
	err = r.queueAddNewSidListWithScoreCmd(pipe, userid, sid)
	if err != nil {
		return err
	}

	// _, err = pipe.Exec(context)
	return nil
}

func (r *RedisStateStore) DeleteSingularSessionID(ctx context.Context, sid string) error {
	commandContext := context.Background()

	//Get User whose owns this sid

	u, err := r.GetUserBySessionID(commandContext, sid)
	if err != nil {
		return err
	}
	// keyToSidList := fmt.Sprintf("%s:%s", r.prefixSidList, u)

	//Delete sid from sid -> user index
	keySidUser := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid)
	_, err = r.client.Del(commandContext, keySidUser).Result()
	if err != nil {
		return err
	}

	//Delete entry from sidlist
	key := fmt.Sprintf("%s:%s", r.prefixSidList, u)
	_, err = r.client.ZRem(commandContext, key, sid).Result()
	if err != nil {
		return err
	}
	return nil

}

func (r *RedisStateStore) DeleteUser(ctx context.Context, userid string) error {

	//Deletes 3 things
	// the sidlist
	// the sid -> user index maps
	// the user HashSet

	// pipe := r.client.TxPipeline()
	commandContext := context.Background()

	//Get key to sid list of user
	keyToSidList := fmt.Sprintf("%s:%s", r.prefixSidList, userid)

	//Collect all the sids that belong to a user from sorted set user:[sid1, sid2...]
	sids, err := r.client.ZRangeByScore(commandContext, keyToSidList, &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	if err != nil {
		return fmt.Errorf("error in getting sorted set range, %s", err)
	}

	//Delete all the sid->user indices
	deleted := 0
	for _, sid := range sids {
		keySidUser := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid)
		val, err := r.client.Del(commandContext, keySidUser).Result()
		if err != nil {
			return err
		}
		deleted += int(val)
	}

	if deleted != len(sids) {
		return fmt.Errorf("could not delete all sids, operation cancelled")
	}

	//Delete sidlist
	val, err := r.client.Del(commandContext, keyToSidList).Result()
	fmt.Println("Deleted: ", val)
	if val > 1 || err != nil {
		return fmt.Errorf("illegal remove or invalid op. Removed %d and error: %w", val, err)
	}

	//Delete User
	userKey := fmt.Sprintf("%s:%s", r.prefixUser, userid)
	val, err = r.client.Del(commandContext, userKey).Result()
	if val == 0 {
		return err
	}
	// _, err = pipe.Exec(ctx)
	return err
}

func (r *RedisStateStore) GetValidSidByUser(context context.Context, userid string, lastValidTime time.Time) ([]string, error) {
	key := fmt.Sprintf("%s:%s", r.prefixSidList, userid)

	vals, err := r.client.ZRevRangeByScore(context, key, &redis.ZRangeBy{
		Min: strconv.Itoa(int(lastValidTime.Unix())),
		Max: "+inf",
	}).Result()

	return vals, err
}

func (r *RedisStateStore) GetUserBySessionID(context context.Context, sid string) (string, error) {
	keyUserSid := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid)
	val, err := r.client.Get(context, keyUserSid).Result()
	if err == redis.Nil {
		return "", nil //We want to return entry string if we don't find any
	}
	return val, err
}

func (r *RedisStateStore) queueAddNewSidUserIndexCmd(pipe redis.Pipeliner, userid string, sid string) error {
	keyUserSid := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid) // user_sids:SID1 = USER1001	//This can have a EXPR
	return pipe.Set(context.Background(), keyUserSid, userid, time.Duration(r.sidExpirationTime)).Err()

}

func (r *RedisStateStore) queueAddNewSidListWithScoreCmd(pipe redis.Pipeliner, userid string, sid string) error {
	key := fmt.Sprintf("%s:%s", r.prefixSidList, userid)

	return pipe.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), //Score = Current unix time which we will later use for expiry purposes
		Member: sid,
	}).Err()

}

func (s *RedisStateStore) SetNewStateSid(ctx context.Context, stateToken, sessionID string, ttl time.Duration) error {
	key := s.prefixAuth + stateToken
	err := s.client.Set(ctx, key, sessionID, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set state SID key= %q : %w", key, err)
	}
	return nil
}

func (s *RedisStateStore) GetSidKey(ctx context.Context, sessionID string) (string, error) {
	key := s.prefixAuth + "sid_key:" + sessionID
	value, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key does not exist
		} else {
			return "", fmt.Errorf("redis get sid key = %q : %w", key, err)
		}
	}
	return value, nil
}

func (s *RedisStateStore) ConsumeStateSID(ctx context.Context, stateToken string) (string, error) {
	key := s.prefixAuth + stateToken
	value, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis get state sid = %q : %w", key, err)
	}
	//Found value, now delete it.

	delErr := s.client.Del(ctx, key).Err()
	if delErr != nil {
		return "", fmt.Errorf("redis del state sid = %q : %w", key, delErr)
	}
	return value, nil
}

func (s *RedisStateStore) SetSidWebSesssionKey(ctx context.Context, sessionID, webSessionKey string, ttl time.Duration) error {
	key := s.prefixAuth + "sid_key:" + sessionID
	err := s.client.Set(ctx, key, webSessionKey, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set sid key= %q : %w", key, err)
	}
	return nil
}

func (s *RedisStateStore) DelSidKey(ctx context.Context, sessionID string) error {
	key := s.prefixAuth + "sid_key:" + sessionID
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis del sid key= %q : %w", key, err)
	}
	return nil
}
