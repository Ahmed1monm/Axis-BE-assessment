package repository

import (
	"context"
	"time"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalanceRepository interface {
	GetBalances(ctx context.Context, accountID primitive.ObjectID) ([]models.Balance, error)
	UpdateBalance(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) error
	CheckAndDeductBalance(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) error
}

type balanceRepository struct {
	db *mongo.Database
}

func NewBalanceRepository(db *mongo.Database) BalanceRepository {
	return &balanceRepository{db: db}
}

func (r *balanceRepository) GetBalances(ctx context.Context, accountID primitive.ObjectID) ([]models.Balance, error) {
	collection := r.db.Collection(models.BalanceCollection)

	filter := bson.M{"account_id": accountID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, utils.DatabaseError("getting balances", err)
	}
	defer cursor.Close(ctx)

	var balances []models.Balance
	if err := cursor.All(ctx, &balances); err != nil {
		return nil, utils.DatabaseError("decoding balances", err)
	}

	if len(balances) == 0 {
		return []models.Balance{}, nil
	}

	return balances, nil
}

func (r *balanceRepository) UpdateBalance(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) error {
	collection := r.db.Collection(models.BalanceCollection)

	filter := bson.M{
		"account_id": accountID,
		"currency":   currency,
	}

	update := bson.M{
		"$inc": bson.M{"amount": amount},
		"$set": bson.M{"updated_at": time.Now()},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return utils.DatabaseError("updating balance", err)
	}

	return nil
}

func (r *balanceRepository) CheckAndDeductBalance(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) error {
	collection := r.db.Collection(models.BalanceCollection)

	filter := bson.M{
		"account_id": accountID,
		"currency":   currency,
		"amount":     bson.M{"$gte": amount},
	}
	update := bson.M{
		"$inc": bson.M{"amount": -amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return utils.DatabaseError("checking and deducting balance", err)
	}
	if result.MatchedCount == 0 {
		return utils.ErrInsufficientBalance
	}
	return nil
}
