package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransactionsVolumeRepository struct {
	db *pgxpool.Pool
}

func NewTransactionsVolumeRepository(db *pgxpool.Pool) *TransactionsVolumeRepository {
	return &TransactionsVolumeRepository{db: db}
}

func (r *TransactionsVolumeRepository) Save(
	ctx context.Context,
	fromCurrency,
	toCurrency domain.Currency,
	amount decimal.Decimal,
) error {
	const query = `
		INSERT INTO transaction_volumes (from_currency, to_currency, volume, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (from_currency, to_currency)
		DO UPDATE SET volume = transaction_volumes.volume + $3, updated_at = CURRENT_TIMESTAMP;
	`

	_, err := r.db.Exec(ctx, query, fromCurrency, toCurrency, amount)

	return err
}

func (r *TransactionsVolumeRepository) GetVolume(
	ctx context.Context,
	fromCurrency domain.Currency,
	toCurrency domain.Currency,
) (decimal.Decimal, error) {
	const query = `
		SELECT volume
		FROM transaction_volumes
		WHERE from_currency = $1 AND to_currency = $2;
	`

	var volume decimal.Decimal

	if err := r.db.QueryRow(ctx, query, fromCurrency, toCurrency).Scan(&volume); err != nil {
		return decimal.Zero, err
	}

	return volume, nil
}
