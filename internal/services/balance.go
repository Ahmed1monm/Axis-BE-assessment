package services

import (
	"context"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BalanceService struct {
	db         *mongo.Database
	repository repository.BalanceRepository
}

func NewBalanceService(db *mongo.Database) *BalanceService {
	return &BalanceService{
		db:         db,
		repository: repository.NewBalanceRepository(db),
	}
}

// GetBalances returns all balances for an account
func (s *BalanceService) GetBalances(ctx context.Context, accountID primitive.ObjectID) (*dtos.BalanceResponse, error) {
	balances, err := s.repository.GetBalances(ctx, accountID)
	if err != nil {
		return nil, err
	}

	response := &dtos.BalanceResponse{
		AccountID: accountID.Hex(),
		Balances:  make([]dtos.CurrencyBalance, len(balances)),
	}

	for i, balance := range balances {
		response.Balances[i] = dtos.CurrencyBalance{
			Currency: balance.Currency,
			Amount:   balance.Amount,
		}
	}

	return response, nil
}
