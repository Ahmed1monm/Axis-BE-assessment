package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB establishes connection to MongoDB and initializes collections
func ConnectDB(uri, dbName string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Info().Msg("Connected to MongoDB")

	// Initialize collections and indexes
	if err = InitializeCollections(client, dbName); err != nil {
		log.Error().Err(err).Msg("Failed to initialize collections")
		return nil, err
	}

	return client, nil
}

// GetCollection returns a handle to a MongoDB collection
func GetCollection(client *mongo.Client, dbName, collName string) *mongo.Collection {
	return client.Database(dbName).Collection(collName)
}
