package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

// PublishState async published events states
type PublishState string

const (
	WaitingAck PublishState = "WAITING_ACK"
	Published               = "PUBLISHED"
	Failed                  = "FAILED"
)

// Connect create a connection to redis server
func Connect(config Config) Redis {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	return Redis{Client: client}
}

// SetEventState set event state by ackId from stan
func (r Redis) SetEventState(ctx context.Context, ackId string, state PublishState) error {
	// TODO: Change This!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// no expiration,
	err := r.Client.Set(ctx, ackId, state, 0).Err()

	if err != nil {
		return fmt.Errorf("failed to set key-value=(%s,%s): %v", ackId, state, err)
	}

	return nil
}

// GetEventState get event state by ackId
func (r Redis) GetEventState(ctx context.Context, ackId string) (string, error) {
	result, err := r.Client.Get(ctx, ackId).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("no state setted for %s", ackId)
	}

	return result, nil
}
