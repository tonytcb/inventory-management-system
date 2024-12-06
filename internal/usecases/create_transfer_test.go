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

func TestCreateTransfer_Create(t *testing.T) {
	type fields struct {
		rateProvider     func(*testing.T) FXRateProvider
		currencyPoolRepo func(*testing.T) CurrencyPoolRepository
		transferRepo     func(*testing.T) TransferRepository
		transactionRepo  func(*testing.T) TransactionRepository
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
						Return(&domain.FXRate{ID: 1, Rate: decimal.RequireFromString("1.1")}, nil).
						Once()
					return m
				},
				currencyPoolRepo: func(t *testing.T) CurrencyPoolRepository {
					m := mocks.NewCurrencyPoolRepository(t)
					m.EXPECT().
						Debit(mock.Anything, domain.EUR, decimal.RequireFromString("111.1")).
						Return(nil).
						Once()
					return m
				},
				transferRepo: func(t *testing.T) TransferRepository {
					wantTransfer := &domain.Transfer{
						OriginalAmount:  decimal.NewFromFloat(100),
						ConvertedAmount: decimal.NewFromFloat(110),
						FinalAmount:     decimal.RequireFromString("111.1"),
						Status:          domain.TransferStatusPending,
						From:            domain.Account{Currency: domain.USD},
						To:              domain.Account{Currency: domain.EUR},
					}

					m := mocks.NewTransferRepository(t)
					m.EXPECT().
						Save(mock.Anything, wantTransfer).
						RunAndReturn(func(_ context.Context, t *domain.Transfer) error {
							t.ID = 1
							t.CreatedAt = time.Now().UTC()
							return nil
						}).
						Once()
					return m
				},
				transactionRepo: func(t *testing.T) TransactionRepository {
					wantTransaction := &domain.Transaction{
						ReferenceID: 1,
						Type:        domain.TransactionTypeTransfer,
						Amount:      decimal.NewFromFloat(110),
						FXRate:      &domain.FXRate{ID: 1, Rate: decimal.RequireFromString("1.1")},
						Revenue:     decimal.RequireFromString("1.1"),
					}

					m := mocks.NewTransactionRepository(t)
					m.EXPECT().
						Save(mock.Anything, wantTransaction).
						Return(nil).
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
					OriginalAmount: decimal.NewFromFloat(100),
					From:           domain.Account{Currency: domain.USD},
					To:             domain.Account{Currency: domain.EUR},
				},
				margin: decimal.NewFromFloat(0.01), // 1%
			},
			want: &domain.Transfer{
				ID:              1,
				OriginalAmount:  decimal.NewFromFloat(100),
				ConvertedAmount: decimal.NewFromFloat(110),
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
				tt.fields.transactionRepo(t),
				tt.fields.txManager(t),
			)

			got, err := c.Create(tt.args.ctx, tt.args.transfer, tt.args.margin)

			tt.wantErr(t, err)
			if err != nil {
				return
			}

			tt.want.CreatedAt = got.CreatedAt

			assert.EqualValues(t, tt.want, got)
		})
	}
}

type txManagerMock struct {
}

func (t txManagerMock) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
