package services

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Ahmed1monm/Axis-BE-assessment/internal/dtos"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/models"
	"github.com/Ahmed1monm/Axis-BE-assessment/internal/repository"
	"github.com/Ahmed1monm/Axis-BE-assessment/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

type AuthService interface {
	Register(ctx context.Context, input dtos.RegisterRequest) (*dtos.AuthResponse, error)
	Login(ctx context.Context, input dtos.LoginRequest) (*dtos.AuthResponse, error)
}

type authService struct {
	accountRepo repository.AccountRepository
}

func NewAuthService(accountRepo repository.AccountRepository) AuthService {
	return &authService{accountRepo: accountRepo}
}

func (s *authService) Register(ctx context.Context, input dtos.RegisterRequest) (*dtos.AuthResponse, error) {
	// Check if email exists
	existingAccount, err := s.accountRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingAccount != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createAccountDTO := &dtos.CreateAccountDTO{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    string(hashedPassword),
		Status:      string(models.AccountStatusActive),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	account, err := s.accountRepo.Create(ctx, createAccountDTO)
	if err != nil {
		return nil, err
	}

	// Generate JWT token using timestamp as uint
	token, err := jwt.GenerateToken(uint(account.ID.Timestamp().Unix()))
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		Token: token,
		User:  account,
	}, nil
}

func (s *authService) Login(ctx context.Context, input dtos.LoginRequest) (*dtos.AuthResponse, error) {
	account, err := s.accountRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrInvalidCredentials
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token using timestamp as uint
	token, err := jwt.GenerateToken(uint(account.ID.Timestamp().Unix()))
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		Token: token,
		User:  account,
	}, nil
}
