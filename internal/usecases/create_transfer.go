package usecases

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type FXRateProvider interface {
	GetLatestRate(ctx context.Context, from domain.Currency, to domain.Currency) (*domain.FXRate, error)
}

type CurrencyPoolRepository interface {
	Debit(context.Context, domain.Currency, decimal.Decimal) error
	Credit(context.Context, domain.Currency, decimal.Decimal) error
	GetAvailableLiquidity(ctx context.Context, currency domain.Currency) (decimal.Decimal, error)
	Rebalance(
		ctx context.Context,
		fromCurrency domain.Currency,
		toCurrency domain.Currency,
		amount decimal.Decimal,
		rate *domain.FXRate,
	) (decimal.Decimal, decimal.Decimal, error)
}

type TransferRepository interface {
	GetByID(ctx context.Context, id int) (*domain.Transfer, error)
	Save(context.Context, *domain.Transfer) error
	UpdateStatus(ctx context.Context, id int, status domain.TransferStatus) error
}

type TransferNotifier interface {
	Created(context.Context, *domain.TransferCreatedEvent) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}

type CreateTransfer struct {
	rateProvider     FXRateProvider
	currencyPoolRepo CurrencyPoolRepository
	transferRepo     TransferRepository
	notifier         TransferNotifier
	txManager        TxManager
}

func NewCreateTransfer(
	rateProvider FXRateProvider,
	currencyPoolRepo CurrencyPoolRepository,
	transferRepo TransferRepository,
	notifier TransferNotifier,
	txManager TxManager,
) *CreateTransfer {
	return &CreateTransfer{
		rateProvider:     rateProvider,
		currencyPoolRepo: currencyPoolRepo,
		transferRepo:     transferRepo,
		notifier:         notifier,
		txManager:        txManager,
	}
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

	transfer.Status = domain.TransferStatusPending

	transfer.ConvertAmounts(rate.Rate, margin)

	err = c.txManager.Do(ctx, func(ctx context.Context) error {
		if err = c.currencyPoolRepo.Debit(ctx, transfer.From.Currency, transfer.OriginalAmount); err != nil {
			return errors.Wrap(err, "error debiting amount from currency pool")
		}

		if err = c.transferRepo.Save(ctx, transfer); err != nil {
			return errors.Wrap(err, "error saving transfer")
		}

		evt := domain.NewTransferCreatedEvent(transfer, rate)
		if err = c.notifier.Created(ctx, evt); err != nil {
			return errors.Wrap(err, "error dispatching transfer created event")
		}

		return nil
	})

	return transfer, err
}
