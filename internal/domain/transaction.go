package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

// Transaction represents a currency transaction
type Transaction struct {
	ID          int
	ReferenceID int
	FXRate      *FXRate
	Type        TransactionType
	Amount      decimal.Decimal
	Revenue     decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
