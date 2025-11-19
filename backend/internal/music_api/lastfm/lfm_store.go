package musicapi

import "github.com/redis/go-redis/v9"

type RedisStateStore struct {
	client *redis.Client
	prefix string
}

func NewRedisStateStore(client *redis.Client) *RedisStateStore {
	return &RedisStateStore{
		client: client,
		prefix: "lfm_cache:",
	}
}
