package model

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event model
type Event struct {
	EventType   string             `bson:"event_type" json:"event_type"`
	Description string             `bson:"description,omitempty" json:"description"`
	Priority    int                `bson:"priority" json:"priority"`
	Payload     string             `bson:"payload,omitempty" json:"payload"`
	CreatedAt   primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}

// Validate event
func (e Event) Validate() error {
	if err := validation.ValidateStruct(&e,
		validation.Field(&e.EventType, validation.Required),
		validation.Field(&e.Priority, validation.Required, validation.Min(0), validation.Max(10)),
		validation.Field(&e.Payload, is.JSON),
	); err != nil {
		return fmt.Errorf("Event validation failed: %w", err)
	}

	return nil
}
