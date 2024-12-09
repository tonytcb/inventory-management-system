package eventbroker

import (
	"context"
	"log/slog"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransferNotifierHandler interface {
	Settlement(context.Context, *domain.Transfer, *domain.FXRate) error
}

type TransferNotifierChannel struct {
	log   *slog.Logger
	queue chan *domain.TransferCreatedEvent
}

func NewTransferNotifierChannel(log *slog.Logger) *TransferNotifierChannel {
	const bufferSize = 100_000

	return &TransferNotifierChannel{
		log:   log,
		queue: make(chan *domain.TransferCreatedEvent, bufferSize),
	}
}

func (t *TransferNotifierChannel) Created(_ context.Context, event *domain.TransferCreatedEvent) error {
	t.queue <- event
	return nil
}

func (t *TransferNotifierChannel) Listen(ctx context.Context, handler TransferNotifierHandler) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-t.queue:
			if !ok {
				return nil
			}

			if err := handler.Settlement(context.Background(), event.Transfer, event.FxRate); err != nil {
				t.log.Error("error settling transfer", "transfer", event.Transfer.ID, "error", err.Error())
				continue
			}

			t.log.Info("transfer settled successfully", "transfer", event.Transfer.ID)
		}
	}
}
