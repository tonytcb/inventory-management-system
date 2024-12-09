package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tonytcb/inventory-management-system/internal/app"
	"github.com/tonytcb/inventory-management-system/internal/app/config"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	cfg, err := config.Load("./config/default.env", "default.env")
	if err != nil {
		logger.Error("failed to load config", "error", err.Error())
		return
	}

	logger.Info("Application env vars", "data", cfg.LogFields())

	appCtx, cancel := context.WithCancel(context.Background())

	application, err := app.NewApplication(appCtx, cfg, logger)
	if err != nil {
		logger.Error("failed to create application", "error", err.Error())
		return
	}

	if err = application.Run(appCtx); err != nil {
		logger.Error("failed to run application:", "error", err.Error())
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	logger.Info("gracefully shutting down...")

	cancel()

	if err = application.Stop(); err != nil {
		logger.Error("failed to stop application", "error", err.Error())
		return
	}
}
