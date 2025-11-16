package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)
type AuthRepository interface {
	SetNewStateSid(ctx context.Context, stateToken, sessionID string, ttl time.Duration) error
	ConsumeStateSID(ctx context.Context, stateToken string) (string, error)
	SetSidWebSesssionKey(ctx context.Context, sessionID, webSessionKey string, ttl time.Duration) (error)
	GetSidKey(ctx context.Context, sessionID string) (string, error)
	DelSidKey(ctx context.Context, sessionID string) error

}

type RedisStateStore struct {
	client *redis.Client
	prefix string
}



func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client: client,
		prefix: "login_state:",	//TODO is this the right way to connect to redis?
	}
}

func (s *RedisStateStore) SetNewStateSid(ctx context.Context, stateToken, sessionID string, ttl time.Duration) error {
	key := s.prefix + stateToken
	err := s.client.Set(ctx, key, sessionID, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set state SID key= %q : %w", key, err)
	}
	return nil
}

func (s* RedisStateStore) GetSidKey(ctx context.Context, sessionID string) (string, error) {
	key := s.prefix + "sid_key:" + sessionID
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
	key := s.prefix + stateToken
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

func (s *RedisStateStore) SetSidWebSesssionKey(ctx context.Context, sessionID, webSessionKey string, ttl time.Duration) (error) {
	key := s.prefix + "sid_key:" + sessionID
	err := s.client.Set(ctx, key, webSessionKey, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set sid key= %q : %w", key, err)
	}
	return nil
}

func (s *RedisStateStore) DelSidKey(ctx context.Context, sessionID string) error {
	key := s.prefix + "sid_key:" + sessionID
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis del sid key= %q : %w", key, err)
	}
	return nil
}





