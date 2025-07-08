package database

import (
	"context"
	"time"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitializeCollections creates collections and indexes if they don't exist
func InitializeCollections(client *mongo.Client, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(dbName)

	// Initialize models with their indexes
	models := []interface{}{
		&models.Account{},
		&models.Transaction{},
		&models.Balance{},
	}

	// Initialize each model's indexes
	for _, model := range models {
		if initializer, ok := model.(CollectionInitializer); ok {
			if err := initializer.EnsureIndexes(ctx, db); err != nil {
				return err
			}
		}
	}

	log.Info().Msg("All collections initialized successfully")
	return nil
}

// CollectionInitializer interface defines the contract for models that need to initialize their collections
type CollectionInitializer interface {
	EnsureIndexes(ctx context.Context, db *mongo.Database) error
}
