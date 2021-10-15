package store

import (
	"context"
	"fmt"
	"github.com/mehditeymorian/etefagh/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// MongoEvent hand
type MongoEvent struct {
	DB *mongo.Database
}

// collection name
const collection = "events"

// NewMongoEvent create mongoEvent
func NewMongoEvent(db *mongo.Database) *MongoEvent {
	return &MongoEvent{DB: db}
}

func (receiver *MongoEvent) Create(ctx context.Context, event model.Event) (interface{}, error) {

	// add created datetime
	event.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	one, err := receiver.DB.Collection(collection).InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("Erro while storing event: %w", err)
	}

	return one.InsertedID, nil
}

func (receiver *MongoEvent) Retrieve(ctx context.Context, eventId string) (*model.Event, error) {

	// check id format
	objectID, _ := primitive.ObjectIDFromHex(eventId)
	result := receiver.DB.Collection(collection).FindOne(ctx, bson.D{{"_id", objectID}})

	var event model.Event

	if err := result.Decode(&event); err != nil {
		return nil, fmt.Errorf("error while decoding document: %w", err)
	}

	return &event, nil
}

func (receiver *MongoEvent) RetrieveAll(ctx context.Context) ([]model.Event, error) {

	cursor, err := receiver.DB.Collection(collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro while finding documents: %w", err)
	}

	var events []model.Event

	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("error while deconding docuements: %w", err)
	}

	return events, nil
}

func (receiver *MongoEvent) Delete(ctx context.Context, eventId string) error {

	// check id format
	objectID, _ := primitive.ObjectIDFromHex(eventId)
	_, err := receiver.DB.Collection(collection).DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("Error while deleting document: %w", err)
	}
	return nil
}
