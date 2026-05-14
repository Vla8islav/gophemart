package repository

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) GetUserOrders(ctx context.Context, userID int64) ([]domain.UserOrder, error) {

	var orders []domain.UserOrder

	err := s.withRetry(ctx, func() error {
		rows, err := s.db.QueryContext(ctx,
			`
	SELECT number, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at DESC`,
			userID,
		)
		if err != nil {
			return fmt.Errorf("failed to get user orders for user id %d: %w", userID, err)
		}
		defer rows.Close()

		result := make([]domain.UserOrder, 0)

		for rows.Next() {

			var order domain.UserOrder
			err = rows.Scan(
				&order.Number,
				&order.Status,
				&order.Accrual,
				&order.UploadedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to scan user orders for user id %d: %w", userID, err)
			}
			result = append(result, order)
		}

		if err = rows.Err(); err != nil {
			return fmt.Errorf("failed to iterate user orders for user id %d: %w", userID, err)
		}

		orders = result

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}
