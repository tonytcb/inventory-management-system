package usecases

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tonytcb/inventory-management-system/internal/domain"
	"github.com/tonytcb/inventory-management-system/test/mocks"
)

func TestSettlementTransfer_Settlement(t *testing.T) {
	type fields struct {
		transferRepo     func(*testing.T) TransferRepository
		currencyPoolRepo func(*testing.T) CurrencyPoolRepository
		transactionRepo  func(*testing.T) TransactionRepository
		volumeRepo       func(*testing.T) TransactionsVolumeRepository
		txManager        func(*testing.T) TxManager
	}
	type args struct {
		ctx      context.Context
		transfer *domain.Transfer
		fxRate   *domain.FXRate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "should settle a transfer successfully",
			fields: fields{
				transferRepo: func(t *testing.T) TransferRepository {
					m := mocks.NewTransferRepository(t)
					m.EXPECT().GetByID(mock.Anything, 1).
						Return(&domain.Transfer{ID: 1, Status: domain.TransferStatusPending}, nil).
						Once()
					m.EXPECT().
						UpdateStatus(mock.Anything, 1, domain.TransferStatusCompleted).
						Return(nil).
						Once()
					return m
				},
				currencyPoolRepo: func(t *testing.T) CurrencyPoolRepository {
					m := mocks.NewCurrencyPoolRepository(t)
					m.EXPECT().
						Credit(mock.Anything, domain.EUR, decimal.RequireFromString("110")).
						Return(nil).
						Once()
					return m
				},
				transactionRepo: func(t *testing.T) TransactionRepository {
					wantTransaction := &domain.Transaction{
						Type:        domain.TransactionTypeTransfer,
						ReferenceID: 1,
						Amount:      decimal.RequireFromString("110"),
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
				volumeRepo: func(t *testing.T) TransactionsVolumeRepository {
					m := mocks.NewTransactionsVolumeRepository(t)
					m.EXPECT().
						Save(mock.Anything, domain.USD, domain.EUR, decimal.RequireFromString("110")).
						Return(nil).
						Once()
					return m
				},
				txManager: func(*testing.T) TxManager {
					return &txManagerMock{}
				},
			},
			args: args{
				ctx: context.Background(),
				transfer: &domain.Transfer{
					ID:              1,
					OriginalAmount:  decimal.NewFromInt(100),
					ConvertedAmount: decimal.NewFromInt(110),
					FinalAmount:     decimal.RequireFromString("111.1"),
					From:            domain.Account{Currency: domain.USD},
					To:              domain.Account{Currency: domain.EUR},
				},
				fxRate: &domain.FXRate{ID: 1, Rate: decimal.RequireFromString("1.1")},
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSettlementTransfer(
				tt.fields.transferRepo(t),
				tt.fields.currencyPoolRepo(t),
				tt.fields.transactionRepo(t),
				tt.fields.volumeRepo(t),
				tt.fields.txManager(t),
			)

			err := s.Settlement(tt.args.ctx, tt.args.transfer, tt.args.fxRate)

			tt.wantErr(t, err)
		})
	}
}
