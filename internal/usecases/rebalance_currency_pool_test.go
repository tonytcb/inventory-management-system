package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/test/mocks"
)

func TestRebalanceCurrencyPool_RebalanceFromTo(t *testing.T) {
	type fields struct {
		currencyPoolRepo func(*testing.T) CurrencyPoolRepository
		volumeRepo       func(*testing.T) TransactionsVolumeRepository
		rateProvider     func(*testing.T) FXRateProvider
		thresholdPercent float64
	}
	type args struct {
		ctx          context.Context
		fromCurrency domain.Currency
		toCurrency   domain.Currency
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "should not rebalance if imbalance is less than threshold",
			fields: fields{
				currencyPoolRepo: func(t *testing.T) CurrencyPoolRepository {
					m := mocks.NewCurrencyPoolRepository(t)
					m.EXPECT().GetAvailableLiquidity(mock.Anything, domain.USD).Return(decimal.NewFromInt(1000), nil).Once()
					m.EXPECT().GetAvailableLiquidity(mock.Anything, domain.EUR).Return(decimal.NewFromInt(950), nil).Once()
					return m
				},
				volumeRepo: func(t *testing.T) TransactionsVolumeRepository {
					m := mocks.NewTransactionsVolumeRepository(t)
					m.EXPECT().GetVolume(mock.Anything, domain.USD, domain.EUR).Return(decimal.NewFromInt(1000), nil).Once()
					return m
				},
				rateProvider: func(t *testing.T) FXRateProvider {
					return mocks.NewFXRateProvider(t)
				},
				thresholdPercent: 5,
			},
			args: args{
				ctx:          context.Background(),
				fromCurrency: domain.USD,
				toCurrency:   domain.EUR,
			},
			want:    false,
			wantErr: require.NoError,
		},
		{
			name: "should rebalance if imbalance is greater than threshold",
			fields: fields{
				currencyPoolRepo: func(t *testing.T) CurrencyPoolRepository {
					m := mocks.NewCurrencyPoolRepository(t)
					m.EXPECT().GetAvailableLiquidity(mock.Anything, domain.USD).Return(decimal.NewFromInt(1000), nil).Once()
					m.EXPECT().GetAvailableLiquidity(mock.Anything, domain.EUR).Return(decimal.NewFromInt(800), nil).Once()
					m.EXPECT().
						Rebalance(mock.Anything, domain.USD, domain.EUR, mock.Anything, mock.Anything).
						Return(decimal.NewFromInt(900), decimal.NewFromInt(900), nil).
						Once()
					return m
				},
				volumeRepo: func(t *testing.T) TransactionsVolumeRepository {
					m := mocks.NewTransactionsVolumeRepository(t)
					m.EXPECT().GetVolume(mock.Anything, domain.USD, domain.EUR).Return(decimal.NewFromInt(1000), nil).Once()
					return m
				},
				rateProvider: func(t *testing.T) FXRateProvider {
					m := mocks.NewFXRateProvider(t)
					m.EXPECT().
						GetLatestRate(mock.Anything, domain.USD, domain.EUR).
						Return(&domain.FXRate{Rate: decimal.NewFromFloat(1.1)}, nil).
						Once()
					return m
				},
				thresholdPercent: 5,
			},
			args: args{
				ctx:          context.Background(),
				fromCurrency: domain.USD,
				toCurrency:   domain.EUR,
			},
			want:    true,
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRebalanceCurrencyPool(
				tt.fields.rateProvider(t),
				tt.fields.currencyPoolRepo(t),
				tt.fields.volumeRepo(t),
				time.Second,
				tt.fields.thresholdPercent,
				[]domain.Currency{},
			)

			got, err := r.RebalanceFromTo(tt.args.ctx, tt.args.fromCurrency, tt.args.toCurrency)
			tt.wantErr(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
