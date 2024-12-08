package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonytcb/inventory-management-system/internal/app/config"
	"github.com/tonytcb/inventory-management-system/internal/infra/storage"
	"github.com/tonytcb/inventory-management-system/internal/usecases"
)

func buildRebalanceUsecase(cfg *config.Config, conn *pgxpool.Pool, _ *slog.Logger) *usecases.RebalanceCurrencyPool {
	var (
		rateProvider     = storage.NewFXRatesRepository(conn)
		currencyPoolRepo = storage.NewCurrencyPoolRepository(conn)
		volumeRepo       = storage.NewTransactionsVolumeRepository(conn)
	)

	return usecases.NewRebalanceCurrencyPool(
		rateProvider,
		currencyPoolRepo,
		volumeRepo,
		cfg.RebalanceCheckInterval,
		cfg.RebalancePoolThresholdPercent,
		cfg.CurrenciesEnabled,
	)
}
