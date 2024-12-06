package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type CurrencyPoolRepository struct {
	db *pgxpool.Pool
}

func NewCurrencyPoolRepository(db *pgxpool.Pool) *CurrencyPoolRepository {
	return &CurrencyPoolRepository{db: db}
}

func (r *CurrencyPoolRepository) Debit(ctx context.Context, currency domain.Currency, amount decimal.Decimal) error {
	const query = `
        WITH deducted_balance AS (
            SELECT id, balance - $1 AS new_balance
            FROM currency_pools_ledger
            WHERE currency_code = $2
            ORDER BY updated_at DESC
            LIMIT 1
            FOR UPDATE
        )
        INSERT INTO currency_pools_ledger (currency_code, balance, updated_at)
        SELECT $2, new_balance, CURRENT_TIMESTAMP
        FROM deducted_balance
        WHERE new_balance >= 0
        RETURNING id, currency_code, balance, updated_at;
    `

	var db QueryRower = r.db
	if tx, ok := extractTxFromContext(ctx); ok {
		db = tx
	}

	var entry domain.CurrencyPool
	err := db.
		QueryRow(ctx, query, amount, string(currency)).
		Scan(&entry.ID, &entry.Currency, &entry.Balance, &entry.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrInsufficientBalance
		}
		return err
	}

	return nil
}

func (r *CurrencyPoolRepository) GetLatest(ctx context.Context, currency domain.Currency) (*domain.CurrencyPool, error) {
	const query = `
		SELECT id, currency_code, balance, updated_at
		FROM currency_pools_ledger
		WHERE currency_code = $1
		ORDER BY updated_at DESC
		LIMIT 1;
	`

	var pool domain.CurrencyPool
	err := r.db.QueryRow(ctx, query, string(currency)).Scan(
		&pool.ID,
		&pool.Currency,
		&pool.Balance,
		&pool.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrCurrencyNotFound
		}
		return nil, err
	}

	return &pool, nil
}
