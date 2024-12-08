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

func TestFXRatesRepository_SaveAndGetLatestRate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dependencies := integration.NewDependencies(ctx, t)

	dbConnection, err := NewPostgresConnection(ctx, dependencies.Cfg)
	require.NoError(t, err)

	truncateTables(ctx, t, dbConnection)

	repo := NewFXRatesRepository(dbConnection)

	var (
		rate1 = &domain.FXRate{FromCurrency: domain.USD, ToCurrency: domain.EUR, Rate: decimal.RequireFromString("0.85")}
		rate2 = &domain.FXRate{FromCurrency: domain.USD, ToCurrency: domain.EUR, Rate: decimal.RequireFromString("0.86")}
		rate3 = &domain.FXRate{FromCurrency: domain.USD, ToCurrency: domain.JPY, Rate: decimal.RequireFromString("110")}
		rate4 = &domain.FXRate{FromCurrency: domain.EUR, ToCurrency: domain.USD, Rate: decimal.RequireFromString("1.18")}
	)

	// Insert FX rates
	require.NoError(t, repo.Save(ctx, rate1))
	require.NoError(t, repo.Save(ctx, rate2))
	require.NoError(t, repo.Save(ctx, rate3))
	require.NoError(t, repo.Save(ctx, rate4))

	// Test GetLatestRate
	fxRate, err := repo.GetLatestRate(ctx, domain.USD, domain.EUR)
	require.NoError(t, err)
	assert.Equal(t, "0.86", fxRate.Rate.String())

	fxRate, err = repo.GetLatestRate(ctx, domain.USD, domain.JPY)
	require.NoError(t, err)
	assert.Equal(t, "110", fxRate.Rate.String())

	fxRate, err = repo.GetLatestRate(ctx, domain.EUR, domain.USD)
	require.NoError(t, err)
	assert.Equal(t, "1.18", fxRate.Rate.String())

	// Test rate not found
	_, err = repo.GetLatestRate(ctx, "BRL", domain.USD)
	assert.ErrorIs(t, err, domain.ErrRateNotFound)
}
