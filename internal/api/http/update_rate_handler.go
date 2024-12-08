package http

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type FxRateUpdater interface {
	Update(ctx context.Context, rate *domain.FXRate) error
}

type updateRateHandler struct {
	log           *slog.Logger
	fxRateUpdater FxRateUpdater
}

func NewUpdateRateHandler(log *slog.Logger, fxRateUpdater FxRateUpdater) UpdateRateHandler {
	return &updateRateHandler{log: log, fxRateUpdater: fxRateUpdater}
}

func (h updateRateHandler) Handle(c *gin.Context) {
	// todo: implement
}
