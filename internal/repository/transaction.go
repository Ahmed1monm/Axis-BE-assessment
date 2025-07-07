package repository

import (
	"context"
	"time"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, dto *dtos.CreateTransactionDTO) (*models.Transaction, error)
}

type transactionRepository struct {
	db *mongo.Database
}

func NewTransactionRepository(db *mongo.Database) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, dto *dtos.CreateTransactionDTO) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:              primitive.NewObjectID(),
		AccountID:       dto.AccountID,
		Type:           models.TransactionType(dto.Type),
		Amount:         dto.Amount,
		Currency:       dto.Currency,
		Status:         models.TransactionStatusCompleted,
		TransactionDate: time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	collection := r.db.Collection(models.TransactionCollection)
	_, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
