package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Event model
type Event struct {
	EventType   string `json:"event_type"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Payload     string `json:"payload"`
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
