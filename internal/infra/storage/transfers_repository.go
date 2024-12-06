package storage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransferRepository struct {
	db *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) Save(ctx context.Context, transfer *domain.Transfer) error {
	const query = `
		INSERT INTO transfers (
			converted_amount,
			final_amount,
			original_amount,
			description,
			status,
			from_currency,
			to_currency
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	var db QueryRower = r.db
	if tx, ok := extractTxFromContext(ctx); ok {
		db = tx
	}

	err := db.QueryRow(
		ctx,
		query,
		transfer.ConvertedAmount,
		transfer.FinalAmount,
		transfer.OriginalAmount,
		transfer.Description,
		transfer.Status,
		transfer.From.Currency,
		transfer.To.Currency,
	).Scan(&transfer.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransferRepository) GetByID(ctx context.Context, id int) (*domain.Transfer, error) {
	const query = `
		SELECT 
		    id,
		    converted_amount,
		    final_amount,
		    original_amount,
		    description,
		    status,
		    from_currency,
		    to_currency,
		    created_at,
		    updated_at
		FROM transfers
		WHERE id = $1;
	`

	var transfer domain.Transfer
	var updatedAt = sql.NullTime{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&transfer.ID,
		&transfer.ConvertedAmount,
		&transfer.FinalAmount,
		&transfer.OriginalAmount,
		&transfer.Description,
		&transfer.Status,
		&transfer.From.Currency,
		&transfer.To.Currency,
		&transfer.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTransferNotFound
		}
		return nil, err
	}

	if updatedAt.Valid {
		transfer.UpdatedAt = updatedAt.Time
	}

	return &transfer, nil
}