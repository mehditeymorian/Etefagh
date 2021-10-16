package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event model
type Event struct {
	EventType   string             `bson:"event_type"`
	Description string             `bson:"description,omitempty"`
	Priority    int                `bson:"priority"`
	Payload     string             `bson:"payload,omitempty"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
	Id          primitive.ObjectID `bson:"_id"`
}
