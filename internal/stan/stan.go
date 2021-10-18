package stan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mehditeymorian/etefagh/internal/model"
	"github.com/mehditeymorian/etefagh/internal/redis"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"github.com/nats-io/stan.go"
)

type Stan struct {
	Connection stan.Conn
	Redis      redis.Redis
	Store      *store.MongoEvent
}

type PublishType string

const (
	Sync  PublishType = "sync"
	Async             = "async"
)

func Connect(config Config) (stan.Conn, error) {

	stanConn, err := stan.Connect(config.ClusterName, config.ClientId, stan.NatsURL(config.Url))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to STAN %v", err)
	}

	return stanConn, nil
}

func (s Stan) Publish(ctx context.Context, publishType PublishType, subject string, event model.Event) error {
	// publish synchronously
	if publishType == Sync {

		err := s.publishSync(subject, event)
		if err != nil {
			return fmt.Errorf("failed to publish synchronously: %v", err)
		}

		return nil
	}

	// publish asynchronously
	ackHandler := stan.AckHandler(func(ackId string, err error) {
		var publishState redis.PublishState
		if err != nil {
			publishState = redis.Failed
		} else {
			publishState = redis.Published
		}
		err = s.Redis.SetEventState(ctx, ackId, publishState)

		if err != nil {
			//TODO: log
		}
	})

	ackId, err := s.publishAsync(subject, event, ackHandler)
	if err != nil {
		return fmt.Errorf("failed to publish asynchronously: %v", err)
	}

	// save ackId to event model
	err = s.Store.UpdateAckId(ctx, event.Id.Hex(), ackId)
	if err != nil {
		return fmt.Errorf("failed to store ackId: %w", err)
	}

	// cache publish state
	err = s.Redis.SetEventState(ctx, ackId, redis.WaitingAck)
	if err != nil {
		return fmt.Errorf("failed to set publish state: %w", err)
	}

	return nil
}

func (s Stan) publishSync(subject string, event model.Event) error {
	if err := validateSubject(subject); err != nil {
		return err
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// publish event to stan
	return s.Connection.Publish(subject, bytes)
}

// PublishAsync publish an event asynchronously and return ackId
func (s Stan) publishAsync(subject string, event model.Event, handler stan.AckHandler) (string, error) {
	if err := validateSubject(subject); err != nil {
		return "", err
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	// publish event to stan
	return s.Connection.PublishAsync(subject, bytes, handler)
}

// validateSubject check if the subject is valid
func validateSubject(subject string) error {
	if subject == "" {
		return fmt.Errorf("failed to publish: subject is not valid")
	}

	return nil
}
