package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) WithdrawFromUserBalance(ctx context.Context,
	userID int64, request domain.UserBalanceWithdrawRequest) error {

	return s.withRetryTx(ctx,
		func(tx *sql.Tx) error { return s.withdrawMoneyTx(ctx, tx, userID, request) },
	)
}

func (s *PostgresStorage) withdrawMoneyTx(ctx context.Context,
	tx *sql.Tx, userID int64, request domain.UserBalanceWithdrawRequest) error {

	// Postgres conditional insert black magic
	query := `
		INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
		SELECT $1, $2, $3, now()
		WHERE (
			SELECT COALESCE(SUM(accrual), 0)
			FROM orders
			WHERE user_id = $1
			  AND status = 'PROCESSED'
		) - (
			SELECT COALESCE(SUM(amount), 0)
			FROM withdrawals
			WHERE user_id = $1
		) >= $3;
	`

	result, err := tx.ExecContext(ctx, query, userID, request.Order, request.Sum)
	if err != nil {
		return fmt.Errorf("withdraw from user balance failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected after withdrawal insert: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotEnoughMoney
	}

	return nil
}
