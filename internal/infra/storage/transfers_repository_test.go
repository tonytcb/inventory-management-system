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

func TestTransferRepository_SaveAndGetByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dependencies := integration.NewDependencies(ctx, t)

	dbConnection, err := NewPostgresConnection(ctx, dependencies.Cfg)
	require.NoError(t, err)

	truncateTables(ctx, t, dbConnection)

	repo := NewTransferRepository(dbConnection)

	// Save transfers
	transfer1 := &domain.Transfer{
		ConvertedAmount: decimal.RequireFromString("100.00"),
		FinalAmount:     decimal.RequireFromString("98.00"),
		OriginalAmount:  decimal.RequireFromString("100.00"),
		Description:     "Transfer 1",
		Status:          domain.TransferStatusPending,
		From:            domain.Account{Currency: domain.USD},
		To:              domain.Account{Currency: domain.EUR},
	}
	transfer2 := &domain.Transfer{
		ConvertedAmount: decimal.RequireFromString("200.00"),
		FinalAmount:     decimal.RequireFromString("196.00"),
		OriginalAmount:  decimal.RequireFromString("200.00"),
		Description:     "Transfer 2",
		Status:          domain.TransferStatusPending,
		From:            domain.Account{Currency: domain.USD},
		To:              domain.Account{Currency: domain.JPY},
	}

	require.NoError(t, repo.Save(ctx, transfer1))
	require.NoError(t, repo.Save(ctx, transfer2))

	// Test GetByID
	got1, err := repo.GetByID(ctx, transfer1.ID)
	require.NoError(t, err)
	assert.Equal(t, transfer1.ID, got1.ID)
	assert.Equal(t, transfer1.ConvertedAmount.String(), got1.ConvertedAmount.String())
	assert.Equal(t, transfer1.FinalAmount.String(), got1.FinalAmount.String())
	assert.Equal(t, transfer1.OriginalAmount.String(), got1.OriginalAmount.String())
	assert.Equal(t, transfer1.Description, got1.Description)
	assert.Equal(t, transfer1.Status, got1.Status)
	assert.Equal(t, transfer1.From.Currency, got1.From.Currency)
	assert.Equal(t, transfer1.To.Currency, got1.To.Currency)

	got2, err := repo.GetByID(ctx, transfer2.ID)
	require.NoError(t, err)
	assert.Equal(t, transfer2.ID, got2.ID)
	assert.Equal(t, transfer2.ConvertedAmount.String(), got2.ConvertedAmount.String())
	assert.Equal(t, transfer2.FinalAmount.String(), got2.FinalAmount.String())
	assert.Equal(t, transfer2.OriginalAmount.String(), got2.OriginalAmount.String())
	assert.Equal(t, transfer2.Description, got2.Description)
	assert.Equal(t, transfer2.Status, got2.Status)
	assert.Equal(t, transfer2.From.Currency, got2.From.Currency)
	assert.Equal(t, transfer2.To.Currency, got2.To.Currency)

	// Test transfer not found
	got3, err := repo.GetByID(ctx, 123)
	require.ErrorIs(t, err, domain.ErrTransferNotFound)
	require.Nil(t, got3)
}
