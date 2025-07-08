package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/services"
	"github.com/Ahmed1monm/Axis-BE-assessment/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{}
	authService := services.NewAuthService(mockAccountRepo)
	ctx := context.Background()

	t.Run("Successful Registration", func(t *testing.T) {
		input := dtos.RegisterRequest{
			Name:        "John Doe",
			Email:       "john@example.com",
			Password:    "password123",
			PhoneNumber: "+1234567890",
		}

		// Mock FindByEmail to return nil (no existing user)
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(nil, nil)

		// Mock Create to return a new account
		mockAccountRepo.On("Create", ctx, mock.AnythingOfType("*dtos.CreateAccountDTO")).Return(&models.Account{
			ID:          primitive.NewObjectID(),
			Name:        input.Name,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
			Status:      models.AccountStatusActive,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)

		response, err := authService.Register(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, input.Email, response.User.Email)
		assert.Equal(t, input.Name, response.User.Name)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		input := dtos.RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		// Mock FindByEmail to return an existing account
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(&models.Account{
			Email: input.Email,
		}, nil)

		response, err := authService.Register(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, services.ErrEmailExists, err)
		assert.Nil(t, response)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		input := dtos.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		// Mock FindByEmail to return an error
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(nil, assert.AnError)

		response, err := authService.Register(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, response)
		mockAccountRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{}
	authService := services.NewAuthService(mockAccountRepo)
	ctx := context.Background()

	t.Run("Successful Login", func(t *testing.T) {
		input := dtos.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		// Hash the password as it would be stored in the database
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		// Mock FindByEmail to return an account with the hashed password
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(&models.Account{
			ID:       primitive.NewObjectID(),
			Email:    input.Email,
			Password: string(hashedPassword),
		}, nil)

		response, err := authService.Login(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, input.Email, response.User.Email)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Invalid Credentials - User Not Found", func(t *testing.T) {
		input := dtos.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		// Mock FindByEmail to return nil (no user found)
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(nil, nil)

		response, err := authService.Login(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, services.ErrInvalidCredentials, err)
		assert.Nil(t, response)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Invalid Credentials - Wrong Password", func(t *testing.T) {
		input := dtos.LoginRequest{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		// Hash a different password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

		// Mock FindByEmail to return an account with a different password
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(&models.Account{
			Email:    input.Email,
			Password: string(hashedPassword),
		}, nil)

		response, err := authService.Login(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, services.ErrInvalidCredentials, err)
		assert.Nil(t, response)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		input := dtos.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		// Mock FindByEmail to return an error
		mockAccountRepo.On("FindByEmail", ctx, input.Email).Return(nil, assert.AnError)

		response, err := authService.Login(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		assert.Nil(t, response)
		mockAccountRepo.AssertExpectations(t)
	})
}
