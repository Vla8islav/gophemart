package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
)

func (m gophermartService) WithdrawFromUserBalance(ctx context.Context, userID int64,
	request domain.UserBalanceWithdrawRequest) error {

	if request.Sum <= 0 {
		return domain.ErrInvalidSum
	}
	if request.Order == "" {
		return domain.ErrOrderEmpty
	}
	if !helpers.ValidLuhn(request.Order) {
		return domain.ErrInvalidOrderNumber
	}

	return m.repository.WithdrawFromUserBalance(ctx, userID, request)
}
