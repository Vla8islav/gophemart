package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (m gophermartService) GetUserBalance(ctx context.Context, userID int64) (*domain.UserBalance, error) {
	return m.repository.GetUserBalance(ctx, userID)
}
