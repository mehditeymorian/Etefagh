package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"
)

type Redis struct {
	Client *redis.Client
	Tracer trace.Tracer
}

// PublishState async published events states
type PublishState string

const (
	WaitingAck PublishState = "WAITING_ACK"
	Published               = "PUBLISHED"
	Failed                  = "FAILED"
)

// Connect create a connection to redis server
func Connect(config Config, tracer trace.Tracer) Redis {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	return Redis{Client: client, Tracer: tracer}
}

// SetEventState set event state by ackId from stan
func (r Redis) SetEventState(ctx context.Context, ackId string, state PublishState) error {
	spanCtx, span := r.Tracer.Start(ctx, "redis.SetEventState")
	defer span.End()

	// TODO: Change This!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// no expiration,
	err := r.Client.Set(spanCtx, ackId, string(state), 0).Err()

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to set key-value=(%s,%s): %w", ackId, state, err)
	}

	return nil
}

// GetEventState get event state by ackId
func (r Redis) GetEventState(ctx context.Context, ackId string) (string, error) {
	ctx, span := r.Tracer.Start(ctx, "redis.GetEventState")
	defer span.End()

	result, err := r.Client.Get(ctx, ackId).Result()
	if err == redis.Nil {
		span.RecordError(err)
		return "", fmt.Errorf("no state setted for %s", ackId)
	}

	return result, nil
}
