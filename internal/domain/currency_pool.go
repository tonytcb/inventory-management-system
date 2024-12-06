package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// CurrencyPool represents a liquidity pool for a specific currency
type CurrencyPool struct {
	ID           int
	Currency     Currency
	Balance      decimal.Decimal
	ExchangeRate decimal.Decimal
	UpdatedAt    time.Time
}

//
//func (cp *CurrencyPool) Deposit(amount decimal.Decimal) {
//	cp.Balance.Add(amount)
//}
//
//func (cp *CurrencyPool) Withdraw(amount decimal.Decimal) error {
//	if cp.Balance.LessThan(amount) {
//		return ErrInsufficientBalance
//	}
//	cp.Balance.Add(amount.Neg())
//	return nil
//}
//
//func (cp *CurrencyPool) UpdateExchangeRate(rate decimal.Decimal) {
//	cp.ExchangeRate = rate
//}
