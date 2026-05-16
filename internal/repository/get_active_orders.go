package repository

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) GetActiveOrders(ctx context.Context) ([]domain.UserOrder, error) {

	var orders []domain.UserOrder

	err := s.withRetry(ctx, func() error {
		rows, err := s.db.QueryContext(ctx,
			`
	SELECT number, status, accrual, uploaded_at
	FROM orders
	WHERE status NOT IN ('INVALID', 'PROCESSED')
	ORDER BY uploaded_at`,
		)
		if err != nil {
			return fmt.Errorf("failed to get active orders: %w", err)
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
				return fmt.Errorf("failed to scan active orders: %w", err)
			}
			result = append(result, order)
		}

		if err = rows.Err(); err != nil {
			return fmt.Errorf("failed to iterate active orders: %w", err)
		}

		orders = result

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}
