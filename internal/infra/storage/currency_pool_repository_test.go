package storage

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/test/integration"
)

func TestCurrencyPoolRepository_Debit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dependencies := integration.NewDependencies(ctx, t)

	dbConnection, err := NewPostgresConnection(ctx, dependencies.Cfg)
	require.NoError(t, err)

	truncateTables(ctx, t, dbConnection)

	repo := NewCurrencyPoolRepository(dbConnection)

	// Insert initial balances
	insertCurrencyPool(ctx, t, dbConnection, domain.USD, decimal.RequireFromString("0"))
	insertCurrencyPool(ctx, t, dbConnection, domain.USD, decimal.RequireFromString("1"))
	insertCurrencyPool(ctx, t, dbConnection, domain.USD, decimal.RequireFromString("2"))
	insertCurrencyPool(ctx, t, dbConnection, domain.USD, decimal.RequireFromString("1000")) // current
	insertCurrencyPool(ctx, t, dbConnection, domain.EUR, decimal.RequireFromString("658"))
	insertCurrencyPool(ctx, t, dbConnection, domain.EUR, decimal.RequireFromString("659.12345678")) // current
	insertCurrencyPool(ctx, t, dbConnection, domain.AUD, decimal.RequireFromString("1000"))

	// Test debit operations
	assert.NoError(t, repo.Debit(ctx, domain.USD, decimal.RequireFromString("100")))
	assert.NoError(t, repo.Debit(ctx, domain.USD, decimal.RequireFromString("0.001")))
	assert.NoError(t, repo.Debit(ctx, domain.USD, decimal.RequireFromString("1.0001")))
	assert.NoError(t, repo.Debit(ctx, domain.EUR, decimal.RequireFromString("0.12345677")))

	// Test insufficient balance
	err = repo.Debit(ctx, domain.USD, decimal.RequireFromString("1000.0000001"))
	assert.Error(t, domain.ErrInsufficientBalance)

	// Verify USD balance
	usd, err := repo.GetLatest(ctx, domain.USD)
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("898.9989").String(), usd.Balance.String())

	// Verify EUR balance
	eur, err := repo.GetLatest(ctx, domain.EUR)
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("659.00000001").String(), eur.Balance.String())

	// Verify AUD balance
	aud, err := repo.GetLatest(ctx, domain.AUD)
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("1000").String(), aud.Balance.String())

	// Verify Not found
	_, err = repo.GetLatest(ctx, domain.JPY)
	require.ErrorIs(t, err, domain.ErrCurrencyNotFound)
}

func truncateTables(ctx context.Context, _ *testing.T, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "TRUNCATE currency_pools_ledger")
	_, _ = db.Exec(ctx, "TRUNCATE fx_rates")
	_, _ = db.Exec(ctx, "TRUNCATE transactions")
	_, _ = db.Exec(ctx, "TRUNCATE transfer")
}

func insertCurrencyPool(ctx context.Context, t *testing.T, db *pgxpool.Pool, currency domain.Currency, balance decimal.Decimal) {
	const query = `
		INSERT INTO currency_pools_ledger (currency_code, balance, updated_at)
		VALUES ($1, $2, $3);
	`
	_, err := db.Exec(ctx, query, currency, balance, time.Now())
	require.NoError(t, err)
}
