package models

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID       primitive.ObjectID `bson:"account_id" json:"account_id" validate:"required"`
	Type            TransactionType    `bson:"type" json:"type" validate:"required"`
	Amount          float64            `bson:"amount" json:"amount" validate:"required,gt=0"`
	Currency        string             `bson:"currency" json:"currency" validate:"required,len=3"` // ISO 4217
	Status          TransactionStatus  `bson:"status" json:"status"`
	Reference       string             `bson:"reference" json:"reference"`
	Description     string             `bson:"description" json:"description"`
	TransactionDate time.Time          `bson:"transaction_date" json:"transaction_date"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "debit"
	TransactionTypeCredit TransactionType = "credit"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// Collection related constants
const (
	TransactionCollection = "transactions"
)

// EnsureIndexes creates the required indexes for the Transaction collection
func (t *Transaction) EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "account_id", Value: 1},
				{Key: "transaction_date", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	}

	col := db.Collection(TransactionCollection)
	_, err := col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Error().Err(err).Str("collection", TransactionCollection).Msg("Failed to create indexes")
		return err
	}

	log.Info().Str("collection", TransactionCollection).Msg("Indexes created successfully")
	return nil
}
