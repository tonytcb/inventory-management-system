package usecases

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransactionsVolumeRepository interface {
	Save(
		ctx context.Context,
		fromCurrency domain.Currency,
		toCurrency domain.Currency,
		amount decimal.Decimal,
	) error

	GetVolume(
		ctx context.Context,
		fromCurrency domain.Currency,
		toCurrency domain.Currency,
	) (decimal.Decimal, error)
}

type RebalanceCurrencyPool struct {
	log              *slog.Logger
	rateProvider     FXRateProvider
	currencyPoolRepo CurrencyPoolRepository
	volumeRepo       TransactionsVolumeRepository
	interval         time.Duration
	thresholdPercent decimal.Decimal
}

func NewRebalanceCurrencyPool(
	rateProvider FXRateProvider,
	currencyPoolRepo CurrencyPoolRepository,
	interval time.Duration,
	thresholdPercent int64,
) *RebalanceCurrencyPool {
	return &RebalanceCurrencyPool{
		log:              slog.Default(),
		rateProvider:     rateProvider,
		currencyPoolRepo: currencyPoolRepo,
		interval:         interval,
		thresholdPercent: decimal.NewFromInt(thresholdPercent),
	}
}

func (r *RebalanceCurrencyPool) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := r.CheckPairs(ctx); err != nil {
				r.log.Error("error checking pairs", "error", err)
			}
		}
	}
}

func (r *RebalanceCurrencyPool) CheckPairs(ctx context.Context) error {
	// @TODO: define currencies via config
	currencies := []domain.Currency{domain.USD, domain.EUR, domain.GBP, domain.JPY, domain.AUD}

	for _, fromCurrency := range currencies {
		for _, toCurrency := range currencies {
			if fromCurrency == toCurrency {
				continue
			}

			if _, err := r.RebalanceFromTo(ctx, fromCurrency, toCurrency); err != nil {
				return errors.Wrapf(err, "error rebalancing from %s to %s", fromCurrency, toCurrency)
			}
		}
	}

	return nil
}

func (r *RebalanceCurrencyPool) RebalanceFromTo(ctx context.Context, fromCurrency, toCurrency domain.Currency) (bool, error) {
	fromLiquidity, err := r.currencyPoolRepo.GetAvailableLiquidity(ctx, fromCurrency)
	if err != nil {
		return false, errors.Wrapf(err, "error getting available liquidity for %s", fromCurrency)
	}

	toLiquidity, err := r.currencyPoolRepo.GetAvailableLiquidity(ctx, toCurrency)
	if err != nil {
		return false, errors.Wrapf(err, "error getting available liquidity for %s", toCurrency)
	}

	// @todo define threshold per currency via config
	thresholdPercent := decimal.NewFromInt(5)

	volume, err := r.volumeRepo.GetVolume(ctx, fromCurrency, toCurrency)
	if err != nil {
		return false, errors.Wrapf(err, "error getting transaction volume for %s to %s", fromCurrency, toCurrency)
	}

	imbalance := fromLiquidity.Sub(toLiquidity).Abs()
	threshold := volume.Mul(thresholdPercent)

	if imbalance.LessThan(threshold) {
		return false, nil
	}

	rebalanceAmount := imbalance.Div(decimal.NewFromInt(2))
	if rebalanceAmount.IsNegative() {
		return false, errors.Errorf("invalid rebalance amount: %s", rebalanceAmount.String())
	}

	rate, err := r.rateProvider.GetLatestRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return false, errors.Wrap(err, "error to get fx rate")
	}

	fromNewBalance, toNewBalance, err := r.currencyPoolRepo.Rebalance(ctx, fromCurrency, toCurrency, rebalanceAmount, rate)
	if err != nil {
		return false, errors.Wrap(err, "error rebalancing currency pool")
	}

	r.log.Info(
		"Rebalanced pair done",
		"from_currency", fromCurrency,
		"to_currency", toCurrency,
		"from_currency_new_balance", fromNewBalance.String(),
		"to_currency_new_balance", toNewBalance.String(),
	)

	return true, nil
}
