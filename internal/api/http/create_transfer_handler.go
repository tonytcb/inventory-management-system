package http

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransferCreator interface {
	Create(
		ctx context.Context,
		transfer *domain.Transfer,
		margin decimal.Decimal,
	) (*domain.Transfer, error)
}

type createTransferHandler struct {
	log             *slog.Logger
	transferCreator TransferCreator
}

func NewCreateTransferHandler(log *slog.Logger, transferCreator TransferCreator) CreateTransferHandler {
	return &createTransferHandler{log: log, transferCreator: transferCreator}
}

func (h createTransferHandler) Handle(c *gin.Context) {
	// todo: implement
}
