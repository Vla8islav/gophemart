package domain

import (
	"context"
)

type GophemartRepository interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user CreateUserParams) (int64, error)
	CreateOrder(ctx context.Context, userID int64, orderNumber string) error
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	GetUserBalance(ctx context.Context, userID int64) (*UserBalance, error)
}
