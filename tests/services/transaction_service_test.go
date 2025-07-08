package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"
	"github.com/Ahmed1monm/Axis-BE-assessment/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a custom service for testing that uses mocked repositories
type testTransactionService struct {
	*services.TransactionService
	mockDB              *mocks.MockMongoDatabase
	mockClient          *mocks.MockMongoClient
	mockSession         *mocks.MockMongoSession
	mockTransactionRepo *mocks.MockTransactionRepository
	mockBalanceRepo     *mocks.MockBalanceRepository
}

func setupTestService() *testTransactionService {
	mockDB := new(mocks.MockMongoDatabase)
	mockClient := new(mocks.MockMongoClient)
	mockSession := new(mocks.MockMongoSession)
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockBalanceRepo := new(mocks.MockBalanceRepository)

	// Setup the chain of mocks
	mockDB.On("Client").Return(mockClient)

	// Create the service with the mocked database
	service := services.NewTransactionService(mockDB)

	// Replace the repositories with mocks
	testService := &testTransactionService{
		TransactionService: service,
		mockDB:             mockDB,
		mockClient:         mockClient,
		mockSession:        mockSession,
		mockTransactionRepo: mockTransactionRepo,
		mockBalanceRepo:    mockBalanceRepo,
	}

	// Use reflection to set the mocked repositories
	// This is a bit of a hack, but necessary for testing with mocks
	serviceValue := reflect.ValueOf(service).Elem()
	
	transactionRepoField := serviceValue.FieldByName("transactionRepo")
	transactionRepoField = reflect.NewAt(transactionRepoField.Type(), unsafe.Pointer(transactionRepoField.UnsafeAddr())).Elem()
	transactionRepoField.Set(reflect.ValueOf(mockTransactionRepo))
	
	balanceRepoField := serviceValue.FieldByName("balanceRepo")
	balanceRepoField = reflect.NewAt(balanceRepoField.Type(), unsafe.Pointer(balanceRepoField.UnsafeAddr())).Elem()
	balanceRepoField.Set(reflect.ValueOf(mockBalanceRepo))

	return testService
}

func TestTransactionService_Deposit(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful Deposit", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 100.0
		currency := "USD"
		transactionID := primitive.NewObjectID()

		// Mock session setup
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("CommitTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock repository calls
		testService.mockBalanceRepo.On("UpdateBalance", mock.Anything, accountID, amount, currency).Return(nil)
		
		expectedCreateDTO := &dtos.CreateTransactionDTO{
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      string(models.TransactionTypeCredit),
		}
		
		testService.mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(dto *dtos.CreateTransactionDTO) bool {
			return dto.AccountID == expectedCreateDTO.AccountID &&
				dto.Amount == expectedCreateDTO.Amount &&
				dto.Currency == expectedCreateDTO.Currency &&
				dto.Type == expectedCreateDTO.Type
		})).Return(&models.Transaction{
			ID:        transactionID,
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      models.TransactionTypeCredit,
		}, nil)

		// Execute
		result, err := testService.TransactionService.Deposit(ctx, accountID, amount, currency)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, transactionID.Hex(), result)
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
		testService.mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Invalid Amount", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		
		// Test with zero amount
		result, err := testService.TransactionService.Deposit(ctx, accountID, 0, "USD")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrInvalidAmount, err)
		assert.Empty(t, result)
		
		// Test with negative amount
		result, err = testService.TransactionService.Deposit(ctx, accountID, -10, "USD")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrInvalidAmount, err)
		assert.Empty(t, result)
	})

	t.Run("Session Start Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 100.0
		currency := "USD"

		// Mock session error
		testService.mockClient.On("StartSession").Return(nil, errors.New("session error"))

		// Execute
		result, err := testService.TransactionService.Deposit(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "session error", err.Error())
		testService.mockClient.AssertExpectations(t)
	})

	t.Run("Transaction Start Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 100.0
		currency := "USD"

		// Mock session setup with transaction error
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(errors.New("transaction error"))
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Execute
		result, err := testService.TransactionService.Deposit(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "transaction error", err.Error())
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
	})

	t.Run("Update Balance Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 100.0
		currency := "USD"

		// Mock session setup
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("AbortTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock repository error
		testService.mockBalanceRepo.On("UpdateBalance", mock.Anything, accountID, amount, currency).Return(errors.New("balance update error"))

		// Execute
		result, err := testService.TransactionService.Deposit(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "balance update error", err.Error())
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
	})

	t.Run("Create Transaction Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 100.0
		currency := "USD"

		// Mock session setup
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("AbortTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock repository calls
		testService.mockBalanceRepo.On("UpdateBalance", mock.Anything, accountID, amount, currency).Return(nil)
		testService.mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil, errors.New("transaction creation error"))

		// Execute
		result, err := testService.TransactionService.Deposit(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "transaction creation error", err.Error())
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
		testService.mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Commit Transaction Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 100.0
		currency := "USD"
		transactionID := primitive.NewObjectID()

		// Mock session setup with commit error
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("CommitTransaction", mock.Anything).Return(errors.New("commit error"))
		testService.mockSession.On("AbortTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock repository calls
		testService.mockBalanceRepo.On("UpdateBalance", mock.Anything, accountID, amount, currency).Return(nil)
		testService.mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(&models.Transaction{
			ID:        transactionID,
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      models.TransactionTypeCredit,
		}, nil)

		// Execute
		result, err := testService.TransactionService.Deposit(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "commit error", err.Error())
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
		testService.mockTransactionRepo.AssertExpectations(t)
	})
}

func TestTransactionService_Withdraw(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful Withdrawal", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 50.0
		currency := "USD"
		transactionID := primitive.NewObjectID()

		// Mock session setup
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("CommitTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock repository calls
		testService.mockBalanceRepo.On("CheckAndDeductBalance", mock.Anything, accountID, amount, currency).Return(nil)
		
		expectedCreateDTO := &dtos.CreateTransactionDTO{
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      string(models.TransactionTypeDebit),
		}
		
		testService.mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(dto *dtos.CreateTransactionDTO) bool {
			return dto.AccountID == expectedCreateDTO.AccountID &&
				dto.Amount == expectedCreateDTO.Amount &&
				dto.Currency == expectedCreateDTO.Currency &&
				dto.Type == expectedCreateDTO.Type
		})).Return(&models.Transaction{
			ID:        transactionID,
			AccountID: accountID,
			Amount:    amount,
			Currency:  currency,
			Type:      models.TransactionTypeDebit,
		}, nil)

		// Execute
		result, err := testService.TransactionService.Withdraw(ctx, accountID, amount, currency)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, transactionID.Hex(), result)
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
		testService.mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Invalid Amount", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		
		// Test with zero amount
		result, err := testService.TransactionService.Withdraw(ctx, accountID, 0, "USD")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrInvalidAmount, err)
		assert.Empty(t, result)
		
		// Test with negative amount
		result, err = testService.TransactionService.Withdraw(ctx, accountID, -10, "USD")
		assert.Error(t, err)
		assert.Equal(t, utils.ErrInvalidAmount, err)
		assert.Empty(t, result)
	})

	t.Run("Insufficient Funds", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 1000.0
		currency := "USD"

		// Mock session setup
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("AbortTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock insufficient funds error
		insufficientFundsErr := utils.NewError(http.StatusBadRequest, "insufficient funds")
		testService.mockBalanceRepo.On("CheckAndDeductBalance", mock.Anything, accountID, amount, currency).Return(insufficientFundsErr)

		// Execute
		result, err := testService.TransactionService.Withdraw(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, insufficientFundsErr, err)
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
	})

	t.Run("Create Transaction Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		amount := 50.0
		currency := "USD"

		// Mock session setup
		testService.mockClient.On("StartSession").Return(testService.mockSession, nil)
		testService.mockSession.On("StartTransaction").Return(nil)
		testService.mockSession.On("AbortTransaction", mock.Anything).Return(nil)
		testService.mockSession.On("EndSession", mock.Anything).Return()

		// Mock repository calls
		testService.mockBalanceRepo.On("CheckAndDeductBalance", mock.Anything, accountID, amount, currency).Return(nil)
		testService.mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil, errors.New("transaction creation error"))

		// Execute
		result, err := testService.TransactionService.Withdraw(ctx, accountID, amount, currency)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "transaction creation error", err.Error())
		testService.mockClient.AssertExpectations(t)
		testService.mockSession.AssertExpectations(t)
		testService.mockBalanceRepo.AssertExpectations(t)
		testService.mockTransactionRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetBalances(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful Get Balances", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		
		balances := []models.Balance{
			{AccountID: accountID, Currency: "USD", Amount: 100.0},
			{AccountID: accountID, Currency: "EUR", Amount: 50.0},
		}
		
		testService.mockBalanceRepo.On("GetBalances", ctx, accountID).Return(balances, nil)

		// Execute
		result, err := testService.TransactionService.GetBalances(ctx, accountID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Balances))
		assert.Equal(t, "USD", result.Balances[0].Currency)
		assert.Equal(t, 100.0, result.Balances[0].Amount)
		assert.Equal(t, "EUR", result.Balances[1].Currency)
		assert.Equal(t, 50.0, result.Balances[1].Amount)
		testService.mockBalanceRepo.AssertExpectations(t)
	})

	t.Run("No Balances Found", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		
		var emptyBalances []models.Balance
		testService.mockBalanceRepo.On("GetBalances", ctx, accountID).Return(emptyBalances, nil)

		// Execute
		result, err := testService.TransactionService.GetBalances(ctx, accountID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Balances)
		testService.mockBalanceRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		// Setup
		testService := setupTestService()
		accountID := primitive.NewObjectID()
		
		testService.mockBalanceRepo.On("GetBalances", ctx, accountID).Return(nil, errors.New("database error"))

		// Execute
		result, err := testService.TransactionService.GetBalances(ctx, accountID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		testService.mockBalanceRepo.AssertExpectations(t)
	})
}
