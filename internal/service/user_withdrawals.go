package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (m gophermartService) GetUserWithdrawals(ctx context.Context, userID int64) ([]domain.Withdrawal, error) {
	return m.repository.GetUserWithdrawals(ctx, userID)
}
