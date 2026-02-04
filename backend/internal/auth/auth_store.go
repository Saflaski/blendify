package auth

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/golang/glog"
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
	GetUserByAnySessionID(context context.Context, sid string) (string, error)
	GetUserByValidSessionID(context context.Context, sid string, expiryDuration time.Duration) (string, error)
	DeleteUser(ctx context.Context, userid string) error
	DeleteSingularSessionID(context context.Context, sid string) error
	GetUserIdByLFM(context context.Context, lfmName string) (string, error)
	AddUserIdToLFMIndex(context context.Context, userid, lfmName string) error
	AddNewSidToExistingUser(ctx context.Context, userid uuid.UUID, validationSid string) error
	GetLFMByUserId(context context.Context, id string) (string, error)
}

type AuthStateStore struct {
	client            *redis.Client
	prefixAuth        string
	prefixUser        string
	prefixSidToUser   string
	prefixSidList     string
	prefixLFMToUser   string
	sidExpirationTime time.Duration
}

func NewRedisStateStore(client *redis.Client, sidExpiryTime time.Duration) *AuthStateStore {
	return &AuthStateStore{
		client:            client,
		prefixAuth:        "login_state:", //TODO is this the right way to connect to redis?
		prefixUser:        "user",
		prefixSidToUser:   "user_sids",
		prefixSidList:     "sidlist",
		prefixLFMToUser:   "lfm_users",
		sidExpirationTime: sidExpiryTime,
	}
}

func (r *AuthStateStore) MakeNewUser(ctx context.Context, validationSid string, userName string, userid uuid.UUID) error {

	key := fmt.Sprintf("%s:%s", r.prefixUser, userid.String())

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()

		pipe.HSet(ctx, key,
			"LFM Username", userName,
			"Created At", time.Now().UTC(),
		)

		r.queueAddNewSid(pipe, userid.String(), validationSid)
		r.queueAddLFMUserIndex(pipe, ctx, userid.String(), userName)

		_, err := pipe.Exec(ctx)
		return err
	}, key)

	if err == redis.TxFailedErr {
		// retry logic
		return r.MakeNewUser(ctx, validationSid, userName, userid)
	}

	return err

}

func (r *AuthStateStore) AddNewSidToExistingUser(ctx context.Context, userid uuid.UUID, validationSid string) error {
	pipe := r.client.TxPipeline()
	r.queueAddNewSid(pipe, userid.String(), validationSid)
	_, err := pipe.Exec(ctx)
	return err
}

// Adds newsid atomically.
func (r *AuthStateStore) queueAddNewSid(pipe redis.Pipeliner, userid string, sid string) error {

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

func (r *AuthStateStore) DeleteSingularSessionID(ctx context.Context, sid string) error {
	commandContext := context.Background()

	//Get User whose owns this sid

	u, err := r.GetUserByAnySessionID(commandContext, sid)
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

func (r *AuthStateStore) DeleteUser(ctx context.Context, userid string) error {

	//Deletes 3 things
	// the sidlist
	// the sid -> user index maps
	// the user HashSet

	// pipe := r.client.TxPipeline()
	commandContext := context.Background()
	pipe := r.client.Pipeline()

	//Delete LFM connection
	platform_name, err := r.GetLFMByUserId(ctx, userid)
	// glog.Info(userid)
	if err != nil {
		glog.Errorf(" during deletion of user auth level, could not find lfm by userid")
	}
	lfmDelKey := fmt.Sprintf("%s:%s", r.prefixLFMToUser, platform_name)

	r.client.Del(ctx, lfmDelKey)

	//Get key to sid list of user
	keyToSidList := fmt.Sprintf("%s:%s", r.prefixSidList, userid)

	//Collect all the sids that belong to a user from sorted set user:[sid1, sid2...]
	sids, err := r.client.ZRangeByScore(commandContext, keyToSidList, &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	if err != nil {
		return fmt.Errorf("error in getting sorted set range, %s", err)
	}

	//Delete all the sid->user index

	for _, sid := range sids {
		pipe.Del(ctx, fmt.Sprintf("%s:%s", r.prefixSidToUser, sid))
	}
	pipe.Del(ctx, keyToSidList)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	r.client.Del(commandContext, keyToSidList)

	// if deleted != len(sids) && len(sids) != 0 {
	// 	return fmt.Errorf("could not delete all sids, operation cancelled")
	// }

	//Delete sidlist
	if _, err := r.client.Del(ctx, keyToSidList).Result(); err != nil {
		return err
	}

	//Delete User
	userKey := fmt.Sprintf("%s:%s", r.prefixUser, userid)
	_, err = r.client.Del(commandContext, userKey).Result()
	// glog.Infof("DEBUG - GDPR Auth level erasure complete: \nuser_id: %s \nsid_count: %d\ntimestamp (UTC): %s", userid, len(sids), time.Now().UTC())
	glog.Infof("GDPR Auth level erasure complete")
	return err
}

func (r *AuthStateStore) GetValidSidByUser(context context.Context, userid string, lastValidTime time.Time) ([]string, error) {
	key := fmt.Sprintf("%s:%s", r.prefixSidList, userid)

	vals, err := r.client.ZRevRangeByScore(context, key, &redis.ZRangeBy{
		Min: strconv.Itoa(int(lastValidTime.Unix())),
		Max: "+inf",
	}).Result()

	return vals, err
}

func (r *AuthStateStore) GetUserByAnySessionID(context context.Context, sid string) (string, error) {

	keyUserSid := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid)
	val, err := r.client.Get(context, keyUserSid).Result()
	if err == redis.Nil {
		return "", nil //We want to return entry string if we don't find any
	}
	return val, err
}

func (r *AuthStateStore) GetUserByValidSessionID(context context.Context, sid string, expiryDuration time.Duration) (string, error) {
	keyUserSid := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid)

	//Get User from SID
	userid, err := r.client.Get(context, keyUserSid).Result()
	if err == redis.Nil {
		glog.Errorf(" could not find sid upon validation: %s", sid)
		return "", nil //We want to return entry string if we don't find any
	}

	lastValidTime := time.Now().Add(-expiryDuration)

	//Eg. If expiryDuration is 24 hours, then it will only return session ids that were made in the last 24 hours
	sids, err := r.GetValidSidByUser(context, userid, lastValidTime)
	// glog.Info(sid)
	// glog.Info("-----")
	// glog.Info(sids)
	//Check if given SID is within these sids and return
	if err != nil {
		return "", err
	}
	if slices.Contains(sids, sid) {
		return userid, nil
	} else {
		return "", nil
	}
}

func (r *AuthStateStore) queueAddNewSidUserIndexCmd(pipe redis.Pipeliner, userid string, sid string) error {
	keyUserSid := fmt.Sprintf("%s:%s", r.prefixSidToUser, sid) // user_sids:SID1 = USER1001	//This can have a EXPR
	return pipe.Set(context.Background(), keyUserSid, userid, r.sidExpirationTime).Err()

}

func (r *AuthStateStore) queueAddNewSidListWithScoreCmd(pipe redis.Pipeliner, userid string, sid string) error {
	key := fmt.Sprintf("%s:%s", r.prefixSidList, userid)

	return pipe.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), //Score = Current unix time which we will later use for expiry purposes
		Member: sid,
	}).Err()

}

func (r *AuthStateStore) AddUserIdToLFMIndex(context context.Context, userid, lfmName string) error {
	key := fmt.Sprintf("%s:%s", r.prefixLFMToUser, lfmName)
	// glog.Info("Set lfm username")
	return r.client.Set(context, key, userid, 0).Err()
}

func (r *AuthStateStore) queueAddLFMUserIndex(pipe redis.Pipeliner, context context.Context, userid, lfmName string) error {
	key := fmt.Sprintf("%s:%s", r.prefixLFMToUser, lfmName)
	return pipe.Set(context, key, userid, 0).Err()
}

func (r *AuthStateStore) GetUserIdByLFM(context context.Context, lfmName string) (string, error) {
	key := fmt.Sprintf("%s:%s", r.prefixLFMToUser, lfmName)
	result, err := r.client.Get(context, key).Result()
	return result, err
}

func (r *AuthStateStore) GetLFMByUserId(context context.Context, id string) (string, error) {
	key := fmt.Sprintf("%s:%s", r.prefixUser, id)
	result, err := r.client.HGet(context, key, "LFM Username").Result()
	return result, err
}

func (s *AuthStateStore) SetNewStateSid(ctx context.Context, stateToken, sessionID string, ttl time.Duration) error {
	key := s.prefixAuth + stateToken
	err := s.client.Set(ctx, key, sessionID, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set state SID key= %q : %w", key, err)
	}
	return nil
}

func (s *AuthStateStore) GetSidKey(ctx context.Context, sessionID string) (string, error) {
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

func (s *AuthStateStore) ConsumeStateSID(ctx context.Context, stateToken string) (string, error) {
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

func (s *AuthStateStore) SetSidWebSesssionKey(ctx context.Context, sessionID, webSessionKey string, ttl time.Duration) error {
	key := s.prefixAuth + "sid_key:" + sessionID
	err := s.client.Set(ctx, key, webSessionKey, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set sid key= %q : %w", key, err)
	}
	return nil
}

func (s *AuthStateStore) DelSidKey(ctx context.Context, sessionID string) error {
	key := s.prefixAuth + "sid_key:" + sessionID
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis del sid key= %q : %w", key, err)
	}
	return nil
}
