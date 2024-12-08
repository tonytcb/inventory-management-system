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
	const bufferSize = 100

	return &TransferNotifierChannel{
		log:   log,
		queue: make(chan *domain.TransferCreatedEvent, bufferSize),
	}
}

func (t *TransferNotifierChannel) Create(_ context.Context, event *domain.TransferCreatedEvent) {
	t.queue <- event
}

func (t *TransferNotifierChannel) Listen(ctx context.Context, handler TransferNotifierHandler) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-t.queue:
			if !ok {
				return
			}

			if err := handler.Settlement(ctx, event.Transfer, event.FxRate); err != nil {
				t.log.Error("error settling transfer", "transfer", event.Transfer.ID, "error", err.Error())
			}
		}
	}
}
