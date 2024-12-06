package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferStatus string

const (
	TransferStatusPending TransferStatus = "pending"
	TransferStatusSuccess                = "succeeded"
	TransferStatusFailed                 = "failed"
)

type Transfer struct {
	ID              int
	ConvertedAmount decimal.Decimal
	FinalAmount     decimal.Decimal
	OriginalAmount  decimal.Decimal
	Description     string
	Status          TransferStatus
	From            Account
	To              Account
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
