package usecases

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type FXRateProvider interface {
	GetLatestRate(ctx context.Context, from domain.Currency, to domain.Currency) (*domain.FXRate, error)
}

type CurrencyPoolRepository interface {
	Debit(context.Context, domain.Currency, decimal.Decimal) error
}

type TransferRepository interface {
	Save(context.Context, *domain.Transfer) error
}

type TransactionRepository interface {
	Save(context.Context, *domain.Transaction) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}

type CreateTransfer struct {
	rateProvider     FXRateProvider
	currencyPoolRepo CurrencyPoolRepository
	transferRepo     TransferRepository
	transactionRepo  TransactionRepository
	txManager        TxManager
}

func (c *CreateTransfer) Create(
	ctx context.Context,
	transfer *domain.Transfer,
	margin decimal.Decimal,
) (*domain.Transfer, error) {
	rate, err := c.rateProvider.GetLatestRate(ctx, transfer.From.Currency, transfer.To.Currency)
	if err != nil {
		return nil, errors.Wrap(err, "error to get fx rate")
	}

	var (
		convertedAmount = transfer.OriginalAmount.Mul(rate.Rate)
		marginAmount    = convertedAmount.Mul(margin)
		finalAmount     = convertedAmount.Sub(marginAmount)
	)

	err = c.txManager.Do(ctx, func(ctx context.Context) error {
		if err = c.currencyPoolRepo.Debit(ctx, transfer.To.Currency, finalAmount); err != nil {
			return errors.Wrap(err, "error debiting amount from currency pool")
		}

		transfer.ConvertedAmount = convertedAmount
		transfer.FinalAmount = finalAmount
		transfer.Status = domain.TransferStatusPending
		transfer.CreatedAt = time.Now().UTC()

		if err = c.transferRepo.Save(ctx, transfer); err != nil {
			return errors.Wrap(err, "error saving transfer")
		}

		transaction := &domain.Transaction{
			Type:        domain.TransactionTypeTransfer,
			ReferenceID: transfer.ID,
			Amount:      transfer.ConvertedAmount,
			FXRate:      rate,
			Revenue:     marginAmount,
			CreatedAt:   time.Now().UTC(),
		}

		if err = c.transactionRepo.Save(ctx, transaction); err != nil {
			return errors.Wrap(err, "error saving transaction")
		}

		return nil
	})

	return transfer, nil
}
