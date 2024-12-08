package usecases

import (
	"context"

	"github.com/pkg/errors"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransactionRepository interface {
	Save(context.Context, *domain.Transaction) error
}

type SettlementTransfer struct {
	transferRepo     TransferRepository
	currencyPoolRepo CurrencyPoolRepository
	transactionRepo  TransactionRepository
	volumeRepo       TransactionsVolumeRepository
	txManager        TxManager
}

func NewSettlementTransfer(
	transferRepo TransferRepository,
	currencyPoolRepo CurrencyPoolRepository,
	transactionRepo TransactionRepository,
	volumeRepo TransactionsVolumeRepository,
	txManager TxManager,
) *SettlementTransfer {
	return &SettlementTransfer{
		transferRepo:     transferRepo,
		currencyPoolRepo: currencyPoolRepo,
		transactionRepo:  transactionRepo,
		volumeRepo:       volumeRepo,
		txManager:        txManager,
	}
}

func (s *SettlementTransfer) Settlement(ctx context.Context, transfer *domain.Transfer, fxRate *domain.FXRate) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		if err := s.isTransferPending(ctx, transfer.ID); err != nil {
			return errors.Wrap(err, "error checking if transfer is pending")
		}

		if err := s.transferRepo.UpdateStatus(ctx, transfer.ID, domain.TransferStatusCompleted); err != nil {
			return errors.Wrap(err, "error updating transfer status to completed")
		}

		if err := s.currencyPoolRepo.Credit(ctx, transfer.To.Currency, transfer.ConvertedAmount); err != nil {
			return errors.Wrap(err, "error crediting amount to currency pool")
		}

		transaction := s.buildTransaction(transfer, fxRate)
		if err := s.transactionRepo.Save(ctx, transaction); err != nil {
			return errors.Wrap(err, "error saving transaction")
		}

		if err := s.volumeRepo.Save(ctx, transfer.From.Currency, transfer.To.Currency, transfer.ConvertedAmount); err != nil {
			return errors.Wrap(err, "error saving transaction")
		}

		return nil
	})
}

func (s *SettlementTransfer) isTransferPending(ctx context.Context, id int) error {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "error getting transfer")
	}

	if !transfer.IsStatusPending() {
		return errors.New("transfer is not pending")
	}

	return nil
}

func (s *SettlementTransfer) buildTransaction(transfer *domain.Transfer, fxRate *domain.FXRate) *domain.Transaction {
	return &domain.Transaction{
		Type:        domain.TransactionTypeTransfer,
		ReferenceID: transfer.ID,
		Amount:      transfer.ConvertedAmount,
		FXRate:      fxRate,
		Revenue:     transfer.Margin(),
	}
}
