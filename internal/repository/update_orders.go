package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) UpdateOrders(ctx context.Context, updates []domain.UpdateOrderParams) error {
	if len(updates) == 0 {
		return nil
	}

	positionalArguments := make([]string, len(updates))
	var values []interface{}

	for i, update := range updates {
		position1 := i*3 + 1
		position2 := i*3 + 2
		position3 := i*3 + 3

		positionalArguments[i] = fmt.Sprintf("($%d::text, $%d::text, $%d::bigint)", position1, position2, position3)
		values = append(values, update.Number, update.Status, update.Accrual)
	}

	return s.withRetryTx(ctx,
		func(tx *sql.Tx) error { return s.batchUpdatesOrdersTx(ctx, tx, positionalArguments, values) },
	)
}

func (s *PostgresStorage) batchUpdatesOrdersTx(ctx context.Context,
	tx *sql.Tx, positionalArguments []string, values []interface{}) error {

	query := fmt.Sprintf(`
		UPDATE orders AS o
		SET
			status = v.status,
			accrual = v.accrual,
			updated_at = now()
		FROM (VALUES %s) AS v(number, status, accrual)
		WHERE o.number = v.number
	`, strings.Join(positionalArguments, ","))
	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("batch update orders failed %v: %w", positionalArguments, err)
	}
	return nil
}
