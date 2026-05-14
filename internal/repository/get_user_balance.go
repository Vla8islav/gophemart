package repository

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (s *PostgresStorage) GetUserBalance(ctx context.Context, userID int64) (*domain.UserBalance, error) {

	var balance domain.UserBalance

	err := s.withRetry(ctx, func() error {
		err := s.db.QueryRowContext(ctx,
			`select
    coalesce(o.total_accrual, 0) - coalesce(w.total_withdrawals, 0) as remaining_accrual,
    coalesce(w.total_withdrawals, 0) as withdrawals_amount
from (
         select coalesce(sum(accrual), 0) as total_accrual
         from orders
         where user_id = $1
           and status = 'PROCESSED'
     ) o
         cross join (
    select coalesce(sum(amount), 0) as total_withdrawals
    from withdrawals
    where user_id = $1
) w`,
			userID,
		).Scan(&balance.Current, &balance.Withdrawn)
		if err != nil {
			return fmt.Errorf("failed to get accrual sum by user id %d: %w", userID, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &balance, nil
}
