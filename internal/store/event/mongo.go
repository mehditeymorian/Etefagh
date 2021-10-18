package store

import (
	"context"
	"fmt"
	"github.com/mehditeymorian/etefagh/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

// MongoEvent hand
type MongoEvent struct {
	DB     *mongo.Database
	Tracer trace.Tracer
}

// collection name
const collection = "events"

// NewMongoEvent create mongoEvent
func NewMongoEvent(db *mongo.Database, tracer trace.Tracer) *MongoEvent {
	return &MongoEvent{DB: db, Tracer: tracer}
}

func (m *MongoEvent) Create(ctx context.Context, event model.Event) (interface{}, error) {
	// initialize a span
	ctx, span := m.Tracer.Start(ctx, "store.events.Create")
	defer span.End()

	one, err := m.DB.Collection(collection).InsertOne(ctx, event)
	if err != nil {
		span.RecordError(err)

		return nil, fmt.Errorf("Erro while storing event: %w", err)
	}

	return one.InsertedID, nil
}

func (m *MongoEvent) Retrieve(ctx context.Context, eventId string) (*model.Event, error) {
	ctx, span := m.Tracer.Start(ctx, "store.events.Retrieve")
	defer span.End()

	// check id format
	objectID, _ := primitive.ObjectIDFromHex(eventId)
	result := m.DB.Collection(collection).FindOne(ctx, bson.D{{"_id", objectID}})

	var event model.Event

	if err := result.Decode(&event); err != nil {
		span.RecordError(err)

		return nil, fmt.Errorf("error while decoding document: %w", err)
	}

	return &event, nil
}

func (m *MongoEvent) RetrieveAll(ctx context.Context) ([]model.Event, error) {
	ctx, span := m.Tracer.Start(ctx, "store.events.RetrieveAll")
	defer span.End()

	cursor, err := m.DB.Collection(collection).Find(ctx, bson.M{})
	if err != nil {
		span.RecordError(err)

		return nil, fmt.Errorf("erro while finding documents: %w", err)
	}

	var events []model.Event

	if err := cursor.All(ctx, &events); err != nil {
		span.RecordError(err)

		return nil, fmt.Errorf("error while deconding docuements: %w", err)
	}

	return events, nil
}

func (m *MongoEvent) Delete(ctx context.Context, eventId string) error {
	ctx, span := m.Tracer.Start(ctx, "store.events.Delete")
	defer span.End()

	// check id format
	objectID, _ := primitive.ObjectIDFromHex(eventId)
	_, err := m.DB.Collection(collection).DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		span.RecordError(err)

		return fmt.Errorf("error while deleting document: %w", err)
	}

	return nil
}

func (m *MongoEvent) UpdateAckId(ctx context.Context, eventId string, ackId string) error {
	ctx, span := m.Tracer.Start(ctx, "store.events.UpdateAckId")
	defer span.End()

	objectID, _ := primitive.ObjectIDFromHex(eventId)
	_, err := m.DB.Collection(collection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"%set", bson.D{{"ack_id", ackId}}},
		})

	if err != nil {
		span.RecordError(err)

		return fmt.Errorf("error while updating ackId: %w", err)
	}

	return nil
}
