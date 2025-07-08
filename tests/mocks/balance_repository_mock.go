package mocks

import (
	"context"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockBalanceRepository struct {
	mock.Mock
}

func (m *MockBalanceRepository) GetBalances(ctx context.Context, accountID primitive.ObjectID) ([]models.Balance, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Balance), args.Error(1)
}

func (m *MockBalanceRepository) UpdateBalance(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) error {
	args := m.Called(ctx, accountID, amount, currency)
	return args.Error(0)
}

func (m *MockBalanceRepository) CheckAndDeductBalance(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) error {
	args := m.Called(ctx, accountID, amount, currency)
	return args.Error(0)
}
