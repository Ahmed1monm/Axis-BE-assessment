package models

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Balance struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID primitive.ObjectID `bson:"account_id" json:"account_id" validate:"required"`
	Amount    float64           `bson:"amount" json:"amount"`
	Currency  string            `bson:"currency" json:"currency" validate:"required,len=3"` // ISO 4217
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// Collection related constants
const (
	BalanceCollection = "balances"
)

// EnsureIndexes creates the required indexes for the Balance collection
func (b *Balance) EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "account_id", Value: 1},
			{Key: "currency", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	col := db.Collection(BalanceCollection)
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Error().Err(err).Str("collection", BalanceCollection).Msg("Failed to create indexes")
		return err
	}

	log.Info().Str("collection", BalanceCollection).Msg("Indexes created successfully")
	return nil
}
