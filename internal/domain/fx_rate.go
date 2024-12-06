package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FXRate struct {
	ID           int
	FromCurrency Currency
	ToCurrency   Currency
	Rate         decimal.Decimal
	UpdatedAt    time.Time
}
