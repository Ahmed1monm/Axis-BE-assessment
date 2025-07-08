package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
)

type AccountRepository interface {
	Create(ctx context.Context, dto *dtos.CreateAccountDTO) (*models.Account, error)
	FindByEmail(ctx context.Context, email string) (*models.Account, error)
}

type accountRepository struct {
	db *mongo.Database
}

func NewAccountRepository(db *mongo.Database) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, dto *dtos.CreateAccountDTO) (*models.Account, error) {
	account := &models.Account{
		Name:        dto.Name,
		Email:       dto.Email,
		PhoneNumber: dto.PhoneNumber,
		Password:    dto.Password,
		Status:      models.AccountStatus(dto.Status),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}

	col := r.db.Collection(models.AccountCollection)
	_, err := col.InsertOne(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*models.Account, error) {
	col := r.db.Collection(models.AccountCollection)
	account := &models.Account{}
	err := col.FindOne(ctx, bson.M{"email": email}).Decode(account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return account, nil
}
