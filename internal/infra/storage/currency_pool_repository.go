package storage

import (
	"context"
	"github.com/pkg/errors"

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

func (r *CurrencyPoolRepository) Credit(ctx context.Context, currency domain.Currency, amount decimal.Decimal) error {
	const query = `
        WITH added_balance AS (
            SELECT id, balance + $1 AS new_balance
            FROM currency_pools_ledger
            WHERE currency_code = $2
            ORDER BY updated_at DESC
            LIMIT 1
            FOR UPDATE
        )
        INSERT INTO currency_pools_ledger (currency_code, balance, updated_at)
        SELECT $2, new_balance, CURRENT_TIMESTAMP
        FROM added_balance
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

func (r *CurrencyPoolRepository) GetAvailableLiquidity(ctx context.Context, currency domain.Currency) (decimal.Decimal, error) {
	latest, err := r.GetLatest(ctx, currency)
	if err != nil {
		return decimal.Zero, err
	}

	return latest.Balance, nil
}

func (r *CurrencyPoolRepository) Rebalance(
	ctx context.Context,
	fromCurrency domain.Currency,
	toCurrency domain.Currency,
	amount decimal.Decimal,
	rate *domain.FXRate,
) (decimal.Decimal, decimal.Decimal, error) {
	convertedAmount := amount.Mul(rate.Rate)

	const debitQuery = `
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

	const creditQuery = `
        WITH added_balance AS (
            SELECT id, balance + $1 AS new_balance
            FROM currency_pools_ledger
            WHERE currency_code = $2
            ORDER BY updated_at DESC
            LIMIT 1
            FOR UPDATE
        )
        INSERT INTO currency_pools_ledger (currency_code, balance, updated_at)
        SELECT $2, new_balance, CURRENT_TIMESTAMP
        FROM added_balance
        RETURNING id, currency_code, balance, updated_at;
    `

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	defer tx.Rollback(ctx)

	// Debit fromCurrency
	var fromEntry domain.CurrencyPool
	err = tx.QueryRow(ctx, debitQuery, amount, string(fromCurrency)).
		Scan(&fromEntry.ID, &fromEntry.Currency, &fromEntry.Balance, &fromEntry.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return decimal.Zero, decimal.Zero, errors.Wrap(domain.ErrInsufficientBalance, "error to debit toCurrency")
		}
		return decimal.Zero, decimal.Zero, err
	}

	// Credit toCurrency
	var toEntry domain.CurrencyPool
	err = tx.QueryRow(ctx, creditQuery, convertedAmount, string(toCurrency)).
		Scan(&toEntry.ID, &toEntry.Currency, &toEntry.Balance, &toEntry.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return decimal.Zero, decimal.Zero, errors.Wrap(domain.ErrInsufficientBalance, "error to credit toCurrency")
		}
		return decimal.Zero, decimal.Zero, err
	}

	return fromEntry.Balance, toEntry.Balance, tx.Commit(ctx)
}
