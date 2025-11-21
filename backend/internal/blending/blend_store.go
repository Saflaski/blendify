package blend

import "github.com/redis/go-redis/v9"

type RedisStateStore struct {
	client *redis.Client
	prefix string
}

func (s *RedisStateStore) GetUser(userA UUID) (string, error) {
	return "saflas", nil
}

type UserListenHistory struct {
	// Define fields for user listen history
}

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client: client,
		prefix: "blends:", //TODO is this the right way to connect to redis?
	}
}

func (s *RedisStateStore) GetUserListenHistory(userID string) (UserListenHistory, error) {
	return UserListenHistory{}, nil
}
