package stan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mehditeymorian/etefagh/internal/model"
	"github.com/mehditeymorian/etefagh/internal/redis"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"github.com/nats-io/stan.go"
	"go.opentelemetry.io/otel/trace"
)

type Stan struct {
	Connection stan.Conn
	Redis      redis.Redis
	Store      *store.MongoEvent
	Tracer     trace.Tracer
}

// PublishType publish type of event
type PublishType string

const (
	Sync  PublishType = "sync"
	Async             = "async"
)

// Connect to nats streaming
func Connect(config Config) (stan.Conn, error) {

	// connect to nats using cluster name, client id, and nats address
	stanConn, err := stan.Connect(config.ClusterName, config.ClientId, stan.NatsURL(config.Url))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to STAN %v", err)
	}

	return stanConn, nil
}

// Publish an event to stan
func (s Stan) Publish(ctx context.Context, publishType PublishType, subject string, event model.Event) error {
	spanCtx, span := s.Tracer.Start(ctx, "stan.publish")
	defer span.End()

	// publish synchronously
	if publishType == Sync {
		err := s.publishSync(subject, event)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to publish synchronously: %w", err)
		}

		return nil
	}

	// publish asynchronously

	// event ack handler
	ackHandler := stan.AckHandler(func(ackId string, err error) {
		// determine final state of ack
		var publishState redis.PublishState
		if err != nil {
			publishState = redis.Failed
		} else {
			publishState = redis.Published
		}
		err = s.Redis.SetEventState(context.Background(), ackId, publishState)

		if err != nil {
			//TODO: log
		}
	})

	// publish event async
	ackId, err := s.publishAsync(subject, event, ackHandler)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to publish asynchronously: %w", err)
	}

	// save ackId to event model
	err = s.Store.UpdateAckId(spanCtx, event.Id.Hex(), ackId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to store ackId: %w", err)
	}

	// cache publish state
	err = s.Redis.SetEventState(spanCtx, ackId, redis.WaitingAck)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to set publish state: %w", err)
	}

	return nil
}

// publishSync publish event synchronously
func (s Stan) publishSync(subject string, event model.Event) error {
	if err := validateSubject(subject); err != nil {
		return err
	}

	// convert events to bytes
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// publish event to stan
	return s.Connection.Publish(subject, bytes)
}

// publishAsync publish an event asynchronously and return ackId
func (s Stan) publishAsync(subject string, event model.Event, handler stan.AckHandler) (string, error) {
	if err := validateSubject(subject); err != nil {
		return "", err
	}

	// convert events to bytes
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
