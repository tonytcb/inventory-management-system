package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

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
	var req updateRateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to parse request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	fromCurrency, toCurrency, err := req.ParseCurrencyPair()
	if err != nil {
		h.log.Error("invalid currency pair parameter", "pair", req.Pair, "error", err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid currency pair"})
		return
	}

	rate, err := decimal.NewFromString(req.Rate)
	if err != nil {
		h.log.Error("invalid rate", "rate", req.Rate, "error", err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid rate"})
		return
	}

	fxRate := &domain.FXRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
		UpdatedAt:    req.Timestamp,
	}

	if err := h.fxRateUpdater.Update(c.Request.Context(), fxRate); err != nil {
		h.log.Error("failed to update rate", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rate"})
		return
	}

	h.log.Info("fx rate updated successfully", "rate", rate)

	c.JSON(http.StatusOK, gin.H{"message": "rate updated successfully"})
}

type updateRateReq struct {
	Pair      string    `json:"pair"`
	Rate      string    `json:"rate"`
	Timestamp time.Time `json:"timestamp"`
}

func (r updateRateReq) ParseCurrencyPair() (domain.Currency, domain.Currency, error) {
	parts := strings.Split(r.Pair, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid currency pair format")
	}

	fromCurrency := domain.Currency(parts[0])
	toCurrency := domain.Currency(parts[1])

	if !fromCurrency.IsValid() || !toCurrency.IsValid() {
		return "", "", errors.New("invalid currency in pair")
	}

	return fromCurrency, toCurrency, nil
}
