package dtos

import "time"

// CreateAccountDTO represents the data needed to create an account in the repository
type CreateAccountDTO struct {
	Name        string
	Email       string
	PhoneNumber string
	Password    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
