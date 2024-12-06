package storage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonytcb/inventory-management-system/internal/domain"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	const query = `
		INSERT INTO transactions (reference_id, type, amount, fx_rate_id, revenue)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`

	var db QueryRower = r.db
	if tx, ok := extractTxFromContext(ctx); ok {
		db = tx
	}

	err := db.QueryRow(
		ctx,
		query,
		transaction.ReferenceID,
		transaction.Type,
		transaction.Amount,
		transaction.FXRate.ID,
		transaction.Revenue,
	).Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) GetByReferenceID(ctx context.Context, referenceID int) (*domain.Transaction, error) {
	const query = `
		SELECT id, reference_id, type, amount, fx_rate_id, revenue, created_at, updated_at
		FROM transactions
		WHERE reference_id = $1
		LIMIT 1;
	`

	var transaction = domain.Transaction{FXRate: &domain.FXRate{}}
	var updatedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, referenceID).Scan(
		&transaction.ID,
		&transaction.ReferenceID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.FXRate.ID,
		&transaction.Revenue,
		&transaction.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, err
	}

	if updatedAt.Valid {
		transaction.UpdatedAt = updatedAt.Time
	}

	return &transaction, nil
}
