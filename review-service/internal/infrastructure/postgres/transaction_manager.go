package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(
	pool *pgxpool.Pool,
) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (m *TransactionManager) Do(
	ctx context.Context,
	callback func(pgx.Tx) error,
) error {
	if err := pgx.BeginFunc(ctx, m.pool, callback); err != nil {
		return fmt.Errorf("execute transaction: %w", err)
	}

	return nil
}
