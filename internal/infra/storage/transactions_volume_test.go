//go:build integration
// +build integration

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

func TestTransactionsVolumeRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dependencies := integration.NewDependencies(ctx, t)

	dbConnection, err := NewPostgresConnection(ctx, dependencies.Cfg)
	require.NoError(t, err)

	truncateTables(ctx, t, dbConnection)

	repo := NewTransactionsVolumeRepository(dbConnection)

	// Test Save method
	err = repo.Save(ctx, domain.USD, domain.EUR, decimal.RequireFromString("100.50"))
	require.NoError(t, err)

	err = repo.Save(ctx, domain.USD, domain.EUR, decimal.RequireFromString("50.25"))
	require.NoError(t, err)

	// Test GetVolume method
	volume, err := repo.GetVolume(ctx, domain.USD, domain.EUR)
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("150.75").String(), volume.String())

	// Test Save method with different currency pair
	err = repo.Save(ctx, domain.EUR, domain.USD, decimal.RequireFromString("200.00"))
	require.NoError(t, err)

	volume, err = repo.GetVolume(ctx, domain.EUR, domain.USD)
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("200.00").String(), volume.String())
}
