package blend

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisStateStore struct {
	client *redis.Client
	prefix string
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
