package usecases

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/test/mocks"
)

func TestCreateTransfer_Create(t *testing.T) {
	wantRateOnSuccess := &domain.FXRate{
		ID:   1,
		Rate: decimal.RequireFromString("1.1"),
	}
	wantTransferOnSuccess := &domain.Transfer{
		ID:              1,
		OriginalAmount:  decimal.NewFromInt(100),
		ConvertedAmount: decimal.NewFromInt(110),
		FinalAmount:     decimal.RequireFromString("111.1"),
		Status:          domain.TransferStatusPending,
		From:            domain.Account{Currency: domain.USD},
		To:              domain.Account{Currency: domain.EUR},
	}

	type fields struct {
		rateProvider     func(*testing.T) FXRateProvider
		currencyPoolRepo func(*testing.T) CurrencyPoolRepository
		transferRepo     func(*testing.T) TransferRepository
		notifier         func(*testing.T) TransferNotifier
		txManager        func(*testing.T) TxManager
	}
	type args struct {
		ctx      context.Context
		transfer *domain.Transfer
		margin   decimal.Decimal
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Transfer
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "should create a transfer successfully",
			fields: fields{
				rateProvider: func(t *testing.T) FXRateProvider {
					m := mocks.NewFXRateProvider(t)
					m.EXPECT().
						GetLatestRate(mock.Anything, domain.USD, domain.EUR).
						Return(wantRateOnSuccess, nil).
						Once()
					return m
				},
				currencyPoolRepo: func(t *testing.T) CurrencyPoolRepository {
					m := mocks.NewCurrencyPoolRepository(t)
					m.EXPECT().
						Debit(mock.Anything, domain.USD, decimal.RequireFromString("100")).
						Return(nil).
						Once()
					return m
				},
				transferRepo: func(t *testing.T) TransferRepository {
					m := mocks.NewTransferRepository(t)
					m.EXPECT().
						Save(mock.Anything, mock.Anything).
						RunAndReturn(func(_ context.Context, got *domain.Transfer) error {
							got.ID = 1

							assert.Equal(t, wantTransferOnSuccess.OriginalAmount.String(), got.OriginalAmount.String())
							assert.Equal(t, wantTransferOnSuccess.ConvertedAmount.String(), got.ConvertedAmount.String())
							assert.Equal(t, wantTransferOnSuccess.FinalAmount.String(), got.FinalAmount.String())
							assert.Equal(t, wantTransferOnSuccess.Status, got.Status)

							return nil
						}).
						Once()
					return m
				},
				notifier: func(t *testing.T) TransferNotifier {
					m := mocks.NewTransferNotifier(t)
					m.EXPECT().
						Created(mock.Anything, mock.Anything).
						RunAndReturn(func(_ context.Context, got *domain.TransferCreatedEvent) error {
							assert.Equal(t, wantTransferOnSuccess.ID, got.Transfer.ID)
							assert.Equal(t, wantRateOnSuccess.Rate.String(), got.FxRate.Rate.String())
							return nil
						}).
						Once()
					return m
				},
				txManager: func(t *testing.T) TxManager {
					return &txManagerMock{}
				},
			},
			args: args{
				ctx: context.Background(),
				transfer: &domain.Transfer{
					OriginalAmount: decimal.NewFromInt(100),
					From:           domain.Account{Currency: domain.USD},
					To:             domain.Account{Currency: domain.EUR},
				},
				margin: decimal.NewFromFloat(0.01), // 1%
			},
			want: &domain.Transfer{
				ID:              1,
				OriginalAmount:  decimal.NewFromInt(100),
				ConvertedAmount: decimal.NewFromInt(110),
				FinalAmount:     decimal.RequireFromString("111.1"),
				Status:          domain.TransferStatusPending,
				From:            domain.Account{Currency: domain.USD},
				To:              domain.Account{Currency: domain.EUR},
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCreateTransfer(
				tt.fields.rateProvider(t),
				tt.fields.currencyPoolRepo(t),
				tt.fields.transferRepo(t),
				tt.fields.notifier(t),
				tt.fields.txManager(t),
			)

			got, err := c.Create(tt.args.ctx, tt.args.transfer, tt.args.margin)

			tt.wantErr(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.OriginalAmount.String(), got.OriginalAmount.String())
			assert.Equal(t, tt.want.ConvertedAmount.String(), got.ConvertedAmount.String())
			assert.Equal(t, tt.want.FinalAmount.String(), got.FinalAmount.String())
			assert.Equal(t, tt.want.Status, got.Status)
			assert.Equal(t, tt.want.From.Currency, got.From.Currency)
			assert.Equal(t, tt.want.To.Currency, got.To.Currency)
		})
	}
}

type txManagerMock struct {
}

func (t txManagerMock) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
