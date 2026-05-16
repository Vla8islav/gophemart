package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

func (m gophermartService) GetUserOrders(ctx context.Context, userID int64) ([]domain.UserOrder, error) {
	return m.repository.GetUserOrders(ctx, userID)
}
