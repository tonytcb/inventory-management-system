package eventbroker

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/test/mocks"
)

func TestTransferNotifierChannel_CreatedAndListen(t *testing.T) {
	logger := slog.Default()

	notifier := NewTransferNotifierChannel(logger, nil)

	handler := mocks.NewTransferNotifierHandler(t)
	handler.EXPECT().
		Settlement(mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := notifier.Listen(ctx, handler)
		assert.NoError(t, err)
	}()

	transfer := &domain.Transfer{
		ID: 1,
		From: domain.Account{
			Currency: domain.Currency("USD"),
		},
		To: domain.Account{
			Currency: domain.Currency("EUR"),
		},
		OriginalAmount: decimal.NewFromInt(100),
	}

	fxRate := &domain.FXRate{
		ID:           100,
		FromCurrency: domain.Currency("USD"),
		ToCurrency:   domain.Currency("EUR"),
		Rate:         decimal.NewFromFloat(0.85),
	}

	event := domain.NewTransferCreatedEvent(transfer, fxRate)
	err := notifier.Created(context.Background(), event)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	handler.AssertCalled(t, "Settlement", mock.Anything, transfer, fxRate)
}
