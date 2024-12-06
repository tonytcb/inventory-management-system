package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type FXRate struct {
	ID           int
	FromCurrency Currency
	ToCurrency   Currency
	Rate         decimal.Decimal
	UpdatedAt    time.Time
}
