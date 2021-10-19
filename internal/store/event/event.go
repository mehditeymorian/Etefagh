package store

import (
	"context"
	"github.com/mehditeymorian/etefagh/internal/model"
)

// Event mongo operations
type Event interface {
	Create(ctx context.Context, event model.Event) (interface{}, error)
	Retrieve(ctx context.Context, eventId string) (*model.Event, error)
	RetrieveAll(ctx context.Context) ([]model.Event, error)
	Delete(ctx context.Context, eventId string) error
	UpdateAckId(ctx context.Context, eventId string, ackId string) error
}
