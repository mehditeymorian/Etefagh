package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event model
type Event struct {
	EventType   string             `bson:"event_type"  json:"event_type,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Priority    int                `bson:"priority" json:"priority,omitempty"`
	Payload     string             `bson:"payload,omitempty" json:"payload,omitempty"`
	AckId       string             `bson:"ack_id,omitempty" json:"ack_id,omitempty"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at,omitempty"`
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}
