package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/tonytcb/inventory-management-system/internal/app/config"
)

type CreateTransferHandler interface {
	Handle(c *gin.Context)
}

type UpdateRateHandler interface {
	Handle(c *gin.Context)
}

type Server struct {
	log *slog.Logger
	cfg *config.Config
	srv *http.Server

	// handlers
	createTransferHandler CreateTransferHandler
	updateRateHandler     UpdateRateHandler
}

func NewServer(
	cfg *config.Config,
	createTransferHandler CreateTransferHandler,
	updateRateHandler UpdateRateHandler,
) *Server {
	return &Server{cfg: cfg,
		createTransferHandler: createTransferHandler,
		updateRateHandler:     updateRateHandler,
	}
}

func (m *Server) Start() error {
	router := gin.Default()

	router.POST("/health", m.createTransferHandler.Handle)
	router.PUT("/fx-rate", m.updateRateHandler.Handle)

	m.srv = &http.Server{
		Addr:    m.cfg.RestAPIPort,
		Handler: router.Handler(),
	}

	if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "error to start http server")
	}

	return nil
}

func (m *Server) Stop() error {
	return m.srv.Shutdown(context.Background())
}
