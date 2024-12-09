package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonytcb/inventory-management-system/internal/api/http"
	"github.com/tonytcb/inventory-management-system/internal/app/config"
	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/internal/infra/eventbroker"
	"github.com/tonytcb/inventory-management-system/internal/infra/storage"
	"github.com/tonytcb/inventory-management-system/internal/usecases"
)

func buildHTTPServer(
	cfg *config.Config,
	conn *pgxpool.Pool,
	log *slog.Logger,
	transferNotifierChan chan *domain.TransferCreatedEvent,
) (*http.Server, error) {
	// Initialize database repositories
	var (
		rateProvider     = storage.NewFXRatesRepository(conn)
		currencyPoolRepo = storage.NewCurrencyPoolRepository(conn)
		transfersRepo    = storage.NewTransferRepository(conn)
	)

	var transferNotifier = eventbroker.NewTransferNotifierChannel(log, transferNotifierChan)
	var txManager = storage.NewTxManager(conn)

	// Initialize use cases
	var (
		createTransferUsecase = usecases.NewCreateTransfer(rateProvider, currencyPoolRepo, transfersRepo, transferNotifier, txManager)
		updateRateUsecase     = usecases.NewFxRateUpdater(rateProvider)
	)

	// Initialize http handlers
	var (
		createTransferHandler = http.NewCreateTransferHandler(cfg, log, createTransferUsecase)
		updateRateHandler     = http.NewUpdateRateHandler(log, updateRateUsecase)
	)

	return http.NewServer(
		cfg,
		createTransferHandler,
		updateRateHandler,
	), nil
}
