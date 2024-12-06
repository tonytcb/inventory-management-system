package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type QueryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{db: db}
}

func (m TxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "error to start transaction")
	}

	if fnErr := fn(putTxToContext(ctx, tx)); fnErr != nil {
		if err := tx.Rollback(ctx); err != nil {
			return errors.Wrap(err, "error to rollback transaction")
		}

		return fnErr
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "error to commit transaction")
	}

	return nil
}

type txKey string

var ctxWithTx = txKey("tx")

func extractTxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx := ctx.Value(ctxWithTx)
	if t, ok := tx.(pgx.Tx); ok {
		return t, true
	}

	return nil, false
}

func putTxToContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxWithTx, tx)
}
