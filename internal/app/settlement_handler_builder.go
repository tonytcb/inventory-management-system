package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonytcb/inventory-management-system/internal/app/config"
	"github.com/tonytcb/inventory-management-system/internal/infra/storage"
	"github.com/tonytcb/inventory-management-system/internal/usecases"
)

func buildSettlementUsecase(_ *config.Config, conn *pgxpool.Pool, _ *slog.Logger) *usecases.SettlementTransfer {
	var (
		currencyPoolRepo = storage.NewCurrencyPoolRepository(conn)
		transfersRepo    = storage.NewTransferRepository(conn)
		transactionsRepo = storage.NewTransactionRepository(conn)
		volumeRepo       = storage.NewTransactionsVolumeRepository(conn)
		txManager        = storage.NewTxManager(conn)
	)

	return usecases.NewSettlementTransfer(transfersRepo, currencyPoolRepo, transactionsRepo, volumeRepo, txManager)
}
