package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const connectionTimeout = 10 * time.Second

// Connect creates a new mongodb connection
func Connect(config Config) (*mongo.Database, error) {
	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Uri))
	if err != nil {
		return nil, fmt.Errorf("db new client error: %w", err)
	}

	// connect to the mongodb
	{
		ctx, done := context.WithTimeout(context.Background(), connectionTimeout)
		defer done()

		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("db connection error: %w", err)
		}
	}
	// ping the mongodb
	{
		ctx, done := context.WithTimeout(context.Background(), connectionTimeout)
		defer done()

		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, fmt.Errorf("db ping error: %w", err)
		}
	}

	return client.Database(config.Name), nil
}
