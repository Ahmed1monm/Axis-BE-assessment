package services

import (
	"context"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/repository"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidAmount = utils.ErrInvalidAmount
)

type TransactionService struct {
	db              *mongo.Database
	transactionRepo repository.TransactionRepository
	balanceRepo     repository.BalanceRepository
}

func NewTransactionService(db *mongo.Database) *TransactionService {
	return &TransactionService{
		db:              db,
		transactionRepo: repository.NewTransactionRepository(db),
		balanceRepo:     repository.NewBalanceRepository(db),
	}
}

func (s *TransactionService) Deposit(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) (string, error) {
	if amount <= 0 {
		return "", utils.ErrInvalidAmount
	}

	session, err := s.db.Client().StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(ctx)

	var transaction *models.Transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Update balance
		if err := s.balanceRepo.UpdateBalance(sc, accountID, amount, currency); err != nil {
			return err
		}

		// Create transaction record
		createDTO := &dtos.CreateTransactionDTO{
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      string(models.TransactionTypeCredit),
		}

		var err error
		transaction, err = s.transactionRepo.CreateTransaction(sc, createDTO)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return "", abortErr
		}
		return "", err
	}

	return transaction.ID.Hex(), nil
}

func (s *TransactionService) GetBalances(ctx context.Context, accountID primitive.ObjectID) (*dtos.BalancesResponse, error) {
	balances, err := s.balanceRepo.GetBalances(ctx, accountID)
	if err != nil {
		return nil, err
	}

	currencyBalances := make([]dtos.CurrencyBalance, len(balances))
	for i, balance := range balances {
		currencyBalances[i] = dtos.CurrencyBalance{
			Currency: balance.Currency,
			Amount:   balance.Amount,
		}
	}

	return &dtos.BalancesResponse{
		Balances: currencyBalances,
	}, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID primitive.ObjectID, amount float64, currency string) (string, error) {
	if amount <= 0 {
		return "", utils.ErrInvalidAmount
	}

	session, err := s.db.Client().StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(ctx)

	var transaction *models.Transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Check and update balance
		if err := s.balanceRepo.CheckAndDeductBalance(sc, accountID, amount, currency); err != nil {
			return err
		}

		// Create transaction record
		createDTO := &dtos.CreateTransactionDTO{
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      string(models.TransactionTypeDebit),
		}

		var err error
		transaction, err = s.transactionRepo.CreateTransaction(sc, createDTO)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return "", abortErr
		}
		return "", err
	}

	return transaction.ID.Hex(), nil
}
