package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type FXRatesRepository struct {
	db *pgxpool.Pool
}

func NewFXRatesRepository(db *pgxpool.Pool) *FXRatesRepository {
	return &FXRatesRepository{db: db}
}

func (r *FXRatesRepository) Save(ctx context.Context, fxRate *domain.FXRate) error {
	const query = `
		INSERT INTO fx_rates (from_currency, to_currency, rate, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	if fxRate.UpdatedAt.IsZero() {
		fxRate.UpdatedAt = time.Now().UTC()
	}

	var db QueryRower = r.db
	if tx, ok := extractTxFromContext(ctx); ok {
		db = tx
	}

	return db.QueryRow(ctx, query, fxRate.FromCurrency, fxRate.ToCurrency, fxRate.Rate, fxRate.UpdatedAt).Scan(&fxRate.ID)
}

func (r *FXRatesRepository) GetLatestRate(ctx context.Context, fromCurrency, toCurrency domain.Currency) (*domain.FXRate, error) {
	const query = `
		SELECT id, from_currency, to_currency, rate, updated_at
		FROM fx_rates
		WHERE from_currency = $1 AND to_currency = $2
		ORDER BY id DESC
		LIMIT 1;
	`

	var fxRate domain.FXRate
	err := r.db.QueryRow(ctx, query, string(fromCurrency), string(toCurrency)).Scan(
		&fxRate.ID,
		&fxRate.FromCurrency,
		&fxRate.ToCurrency,
		&fxRate.Rate,
		&fxRate.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRateNotFound
		}
		return nil, err
	}

	return &fxRate, nil
}
