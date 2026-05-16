package domain

import (
	"context"
)

type GophemartService interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, request UserRegisterRequest) (*AuthResult, error)
	LoginUser(ctx context.Context, request UserLoginRequest) (*AuthResult, error)

	GetUserBalance(ctx context.Context, userID int64) (*UserBalance, error)
	WithdrawFromUserBalance(ctx context.Context, userID int64, request UserBalanceWithdraw) error
	GetUserWithdrawals(ctx context.Context, userID int64) ([]Withdrawal, error)

	CreateOrder(ctx context.Context, userID int64, orderNumber string) error
	GetUserOrders(ctx context.Context, userID int64) ([]UserOrder, error)

	PollAccrual(ctx context.Context) error
	StartAccrualPolling(ctx context.Context) error
}
