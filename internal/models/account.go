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

type Account struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string            `bson:"name" json:"name" validate:"required"`
	Email       string            `bson:"email" json:"email" validate:"required,email"`
	PhoneNumber string            `bson:"phone_number" json:"phone_number" validate:"required"`
	Password    string            `bson:"password" json:"-"` // Password is never returned in JSON
	Status      AccountStatus     `bson:"status" json:"status"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusBlocked  AccountStatus = "blocked"
)

// Collection related constants
const (
	AccountCollection = "accounts"
)

// EnsureIndexes creates the required indexes for the Account collection
func (a *Account) EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "phone_number", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	col := db.Collection(AccountCollection)
	_, err := col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Error().Err(err).Str("collection", AccountCollection).Msg("Failed to create indexes")
		return err
	}

	log.Info().Str("collection", AccountCollection).Msg("Indexes created successfully")
	return nil
}
