package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (m metricsService) GetUserBalance(ctx context.Context, userID int64) (*domain.UserBalance, error) {
	return m.repository.GetUserBalance(ctx, userID)
}
