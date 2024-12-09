package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
)

type Transfer struct {
	ID              int
	ConvertedAmount decimal.Decimal `json:"converted_amount"`
	FinalAmount     decimal.Decimal `json:"final_amount"`
	OriginalAmount  decimal.Decimal `json:"original_amount"`
	Description     string          `json:"description"`
	Status          TransferStatus  `json:"-"`
	From            Account         `json:"from"`
	To              Account         `json:"to"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"-"`
}

func (t *Transfer) ConvertAmounts(rate decimal.Decimal, margin decimal.Decimal) {
	t.ConvertedAmount = t.OriginalAmount.Mul(rate)

	var marginAmount = t.ConvertedAmount.Mul(margin)

	t.FinalAmount = t.ConvertedAmount.Add(marginAmount)
}

func (t *Transfer) IsStatusPending() bool {
	return t.Status == TransferStatusPending
}

func (t *Transfer) Margin() decimal.Decimal {
	if m := t.FinalAmount.Sub(t.ConvertedAmount); m.GreaterThan(decimal.Zero) {
		return m
	}

	return decimal.Zero
}

type TransferCreatedEvent struct {
	Transfer *Transfer
	FxRate   *FXRate
}

func NewTransferCreatedEvent(transfer *Transfer, fxRate *FXRate) *TransferCreatedEvent {
	return &TransferCreatedEvent{
		Transfer: transfer,
		FxRate:   fxRate,
	}
}
