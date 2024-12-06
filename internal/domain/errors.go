package domain

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrRateNotFound        = errors.New("rate not found")
	ErrCurrencyNotFound    = errors.New("currency not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrTransferNotFound    = errors.New("transfer not found")
)
