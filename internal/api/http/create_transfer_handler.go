package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/app/config"
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
	cfg             *config.Config
	log             *slog.Logger
	transferCreator TransferCreator
}

func NewCreateTransferHandler(cfg *config.Config, log *slog.Logger, transferCreator TransferCreator) CreateTransferHandler {
	return &createTransferHandler{cfg: cfg, log: log, transferCreator: transferCreator}
}

func (h createTransferHandler) Handle(c *gin.Context) {
	var req createTransferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to parse request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := req.IsValid(); err != nil {
		h.log.Error("Invalid Create Transfer request", "body", req, "error", err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	transfer := req.ToTransferDomain()
	margin := decimal.NewFromFloat(h.cfg.RevenueMarginPercent).Div(decimal.NewFromInt(100))

	createdTransfer, err := h.transferCreator.Create(c.Request.Context(), transfer, margin)
	if err != nil {
		if errors.Is(err, domain.ErrCurrencyNotFound) {
			h.log.Info("currency not found", "transfer", transfer, "error", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "currency not found"})
			return
		}

		if errors.Is(err, domain.ErrRateNotFound) {
			h.log.Info("currency rate not found", "body", req, "error", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "currency rate not found"})
			return
		}

		if errors.Is(err, domain.ErrInsufficientBalance) {
			h.log.Info("insufficient balance on pool", "body", req, "error", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "insufficient balance on pool"})
			return
		}

		h.log.Error("failed to create transfer", "body", req, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transfer"})
		return
	}

	h.log.Info("transfer created successfully", "transfer", transfer)

	c.JSON(http.StatusCreated, gin.H{"transfer": createdTransfer})
}

type createTransferReq struct {
	FromAccount struct {
		Currency string `json:"currency"`
	} `json:"from_account"`
	ToAccount struct {
		Currency string `json:"currency"`
	} `json:"to_account"`
	Amount string `json:"amount"`
}

func (req createTransferReq) IsValid() error {
	fromCurrency := domain.Currency(req.FromAccount.Currency)
	if !fromCurrency.IsValid() {
		return errors.New("invalid from_account.currency")
	}

	toCurrency := domain.Currency(req.ToAccount.Currency)
	if !toCurrency.IsValid() {
		return errors.New("invalid to_account.currency")
	}

	if _, err := decimal.NewFromString(req.Amount); err != nil {
		return errors.New("invalid amount")
	}

	return nil
}

func (req createTransferReq) ToTransferDomain() *domain.Transfer {
	fromCurrency := domain.Currency(req.FromAccount.Currency)
	toCurrency := domain.Currency(req.ToAccount.Currency)
	amount := decimal.RequireFromString(req.Amount)

	return &domain.Transfer{
		From:           domain.Account{Currency: fromCurrency},
		To:             domain.Account{Currency: toCurrency},
		OriginalAmount: amount,
	}
}
