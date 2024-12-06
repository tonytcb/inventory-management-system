package storage

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/test/integration"
)

func TestTransactionRepository_SaveAndGetByReferenceID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dependencies := integration.NewDependencies(ctx, t)

	dbConnection, err := NewPostgresConnection(ctx, dependencies.Cfg)
	require.NoError(t, err)

	truncateTables(ctx, t, dbConnection)

	repo := NewTransactionRepository(dbConnection)
	rateRepo := NewFXRatesRepository(dbConnection)

	var rate1 = &domain.FXRate{FromCurrency: domain.USD, ToCurrency: domain.EUR, Rate: decimal.RequireFromString("0.85")}

	// Save initial rate
	require.NoError(t, rateRepo.Save(ctx, rate1))

	var (
		transaction1 = &domain.Transaction{
			ReferenceID: 1,
			Type:        domain.TransactionTypeTransfer,
			Amount:      decimal.RequireFromString("100"),
			FXRate:      &domain.FXRate{ID: rate1.ID},
			Revenue:     decimal.RequireFromString("1"),
		}
		transaction2 = &domain.Transaction{
			ReferenceID: 2,
			Type:        domain.TransactionTypeTransfer,
			Amount:      decimal.RequireFromString("200"),
			FXRate:      &domain.FXRate{ID: rate1.ID},
			Revenue:     decimal.RequireFromString("2"),
		}
	)

	// Save transactions
	require.NoError(t, repo.Save(ctx, transaction1))
	require.NoError(t, repo.Save(ctx, transaction2))

	// Test GetByReferenceID
	got1, err := repo.GetByReferenceID(ctx, transaction1.ReferenceID)
	require.NoError(t, err)
	assert.Equal(t, transaction1.ReferenceID, got1.ReferenceID)
	assert.Equal(t, transaction1.Type, got1.Type)
	assert.Equal(t, transaction1.Amount.String(), got1.Amount.String())
	assert.Equal(t, transaction1.FXRate.ID, got1.FXRate.ID)
	assert.Equal(t, transaction1.Revenue.String(), got1.Revenue.String())

	got2, err := repo.GetByReferenceID(ctx, transaction2.ReferenceID)
	require.NoError(t, err)
	assert.Equal(t, transaction2.ReferenceID, got2.ReferenceID)
	assert.Equal(t, transaction2.Type, got2.Type)
	assert.Equal(t, transaction2.Amount.String(), got2.Amount.String())
	assert.Equal(t, transaction2.FXRate.ID, got2.FXRate.ID)
	assert.Equal(t, transaction2.Revenue.String(), got2.Revenue.String())

	// Test transaction not found
	got3, err := repo.GetByReferenceID(ctx, 123)
	assert.ErrorIs(t, err, domain.ErrTransactionNotFound)
	assert.Nil(t, got3)
}
