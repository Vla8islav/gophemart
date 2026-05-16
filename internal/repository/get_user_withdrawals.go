package repository

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) GetUserWithdrawals(ctx context.Context, userID int64) ([]domain.Withdrawal, error) {
	var withdrawals []domain.Withdrawal

	err := s.withRetry(ctx, func() error {
		rows, err := s.db.QueryContext(ctx,
			`
	SELECT id, user_id, order_number, amount, processed_at
	FROM withdrawals
	WHERE user_id = $1
	ORDER BY processed_at DESC`,
			userID,
		)
		if err != nil {
			return fmt.Errorf("failed to get user withdrawals for user id %d: %w", userID, err)
		}
		defer rows.Close()

		result := make([]domain.Withdrawal, 0)

		for rows.Next() {
			var withdrawal domain.Withdrawal
			err = rows.Scan(
				&withdrawal.ID,
				&withdrawal.UserID,
				&withdrawal.OrderNumber,
				&withdrawal.Amount,
				&withdrawal.ProcessedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to scan user withdrawals for user id %d: %w", userID, err)
			}
			result = append(result, withdrawal)
		}

		if err = rows.Err(); err != nil {
			return fmt.Errorf("failed to iterate user withdrawals for user id %d: %w", userID, err)
		}

		withdrawals = result

		return nil
	})

	if err != nil {
		return nil, err
	}

	return withdrawals, nil

}
