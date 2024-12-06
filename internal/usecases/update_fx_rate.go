package usecases

import (
	"context"
	"github.com/pkg/errors"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type FXRateRepository interface {
	Save(context.Context, *domain.FXRate) error
}

type FxRateUpdater struct {
	rateRepo FXRateRepository
}

func NewFxRateUpdater(rateRepo FXRateRepository) *FxRateUpdater {
	return &FxRateUpdater{rateRepo: rateRepo}
}

func (u *FxRateUpdater) Update(ctx context.Context, rate *domain.FXRate) error {
	if err := u.rateRepo.Save(ctx, rate); err != nil {
		return errors.Wrap(err, "error saving FX rate")
	}

	return nil
}
