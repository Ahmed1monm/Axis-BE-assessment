package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

// TransactionRequest represents the transaction request data
type TransactionRequest struct {
	AccountID string  `json:"account_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Currency  string  `json:"currency" validate:"required,len=3"`
}

// CreateTransactionDTO represents the data needed to create a transaction
type CreateTransactionDTO struct {
	AccountID primitive.ObjectID
	Amount    float64
	Currency  string
	Type      string
}

// TransactionResponse represents the transaction response data
type TransactionResponse struct {
	TransactionID string `json:"transaction_id"`
}
