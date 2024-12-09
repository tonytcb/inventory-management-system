package app

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/tonytcb/inventory-management-system/internal/api/http"
	"github.com/tonytcb/inventory-management-system/internal/app/config"
	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/internal/infra/eventbroker"
	"github.com/tonytcb/inventory-management-system/internal/infra/storage"
	"github.com/tonytcb/inventory-management-system/internal/usecases"
)

type Application struct {
	cfg               *config.Config
	log               *slog.Logger
	httpServer        *http.Server
	databaseConn      *pgxpool.Pool
	transferNotifier  *eventbroker.TransferNotifierChannel
	settlementHandler *usecases.SettlementTransfer
	rebalanceHandler  *usecases.RebalanceCurrencyPool
}

func NewApplication(ctx context.Context, cfg *config.Config, log *slog.Logger) (*Application, error) {
	dbConn, err := storage.NewPostgresConnection(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database connection")
	}

	transferNotifierChan := make(chan *domain.TransferCreatedEvent, 100)

	httpServer, err := buildHTTPServer(cfg, dbConn, log, transferNotifierChan)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http server")
	}

	return &Application{
		cfg:               cfg,
		log:               log,
		httpServer:        httpServer,
		databaseConn:      dbConn,
		transferNotifier:  eventbroker.NewTransferNotifierChannel(log, transferNotifierChan),
		settlementHandler: buildSettlementUsecase(cfg, dbConn, log),
		rebalanceHandler:  buildRebalanceUsecase(cfg, dbConn, log),
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	errGroup, _ := errgroup.WithContext(ctx)

	a.log.Info("Running application")

	errGroup.Go(func() error {
		a.log.Info("Starting http server", "port", a.cfg.RestAPIPort)

		return a.httpServer.Start()
	})

	errGroup.Go(func() error {
		a.log.Info("Starting settlements event listener")

		return a.transferNotifier.Listen(ctx, a.settlementHandler)
	})

	errGroup.Go(func() error {
		a.log.Info("Starting rebalancing job")

		return a.rebalanceHandler.Start(ctx)
	})

	return nil
}

func (a *Application) Stop() error {
	a.log.Info("Stopping application")

	_ = a.httpServer.Stop()
	a.databaseConn.Close()

	return nil
}
