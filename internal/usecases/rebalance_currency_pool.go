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
	pairs            []domain.CurrencyPair
}

func NewRebalanceCurrencyPool(
	rateProvider FXRateProvider,
	currencyPoolRepo CurrencyPoolRepository,
	volumeRepo TransactionsVolumeRepository,
	interval time.Duration,
	thresholdPercent float64,
	currenciesEnabled []domain.Currency,
) *RebalanceCurrencyPool {
	return &RebalanceCurrencyPool{
		log:              slog.Default(),
		rateProvider:     rateProvider,
		currencyPoolRepo: currencyPoolRepo,
		volumeRepo:       volumeRepo,
		interval:         interval,
		thresholdPercent: decimal.NewFromFloat(thresholdPercent),
		pairs:            domain.BuildPairs(currenciesEnabled),
	}
}

func (r *RebalanceCurrencyPool) Start(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			if err := r.CheckPairs(ctx); err != nil {
				r.log.Error("error checking pairs", "error", err.Error())
			}
		}
	}
}

func (r *RebalanceCurrencyPool) CheckPairs(ctx context.Context) error {
	for _, pair := range r.pairs {
		if _, err := r.RebalanceFromTo(ctx, pair.From, pair.To); err != nil {
			return errors.Wrapf(err, "error rebalancing from %s to %s", pair.From, pair.To)
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

	volume, err := r.volumeRepo.GetVolume(ctx, fromCurrency, toCurrency)
	if err != nil {
		return false, errors.Wrapf(err, "error getting transaction volume for %s to %s", fromCurrency, toCurrency)
	}

	if volume.IsZero() {
		return false, nil
	}

	imbalance := fromLiquidity.Sub(toLiquidity)

	threshold := volume.Mul(r.thresholdPercent.Div(decimal.NewFromInt(100)))
	if imbalance.Abs().LessThanOrEqual(threshold) {
		return false, nil
	}

	if imbalance.IsNegative() {
		r.log.Info("imbalance is negative, attempt to reverse pair", "from_currency", fromCurrency, "to_currency", toCurrency)
		return r.RebalanceFromTo(ctx, toCurrency, fromCurrency)
	}

	rebalanceAmount := imbalance.Div(decimal.NewFromInt(2))

	r.log.Info(
		"start rebalancing pair",
		"from_currency", fromCurrency,
		"from_currency_liquidity", fromLiquidity.String(),
		"to_currency", toCurrency,
		"to_currency_liquidity", toLiquidity.String(),
		"imbalance", imbalance.String(),
		"threshold", threshold.String(),
		"rebalance_amount", rebalanceAmount.String(),
	)

	rate, err := r.rateProvider.GetLatestRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return false, errors.Wrap(err, "error to get fx rate")
	}

	r.log.Info("rebalancing amounts", "rebalance_amount", rebalanceAmount.String(), "rate", rate.Rate.String())

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
