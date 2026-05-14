package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) CreateOrder(ctx context.Context, userID int64, orderNumber string) error {

	err := s.withRetry(ctx, func() error {
		var existingUserID int64

		err := s.db.QueryRowContext(ctx,
			"SELECT user_id FROM orders WHERE number = $1",
			orderNumber).Scan(&existingUserID)

		if err == nil {
			// Такой заказ уже нашелся
			if existingUserID == userID {
				return domain.ErrOrderAlreadyUploadedByUser
			}

			return domain.ErrOrderAlreadyUploadedByAnotherUser
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to check the existing order UserID %d Number %s: %w",
				userID, orderNumber, err)
		}

		_, err = s.db.ExecContext(ctx, `
			INSERT INTO orders (number, user_id, status)
			VALUES ($1, $2, 'NEW')
		`, orderNumber, userID)
		if err != nil {
			return fmt.Errorf("create order failed UserID %d Number %q: %w", userID, orderNumber, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("create order failed UserID %d Number %q: %w", userID, orderNumber, err)
	}

	return nil
}
