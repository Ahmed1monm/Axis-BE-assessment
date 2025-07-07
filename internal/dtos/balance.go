package dtos

// BalanceResponse represents the response for getting account balances
type BalanceResponse struct {
	AccountID string            `json:"account_id"`
	Balances  []CurrencyBalance `json:"balances"`
}

// CurrencyBalance represents a balance for a specific currency
type CurrencyBalance struct {
	Currency string  `json:"currency" validate:"required,len=3"`
	Amount   float64 `json:"amount" validate:"required,gte=0"`
}
