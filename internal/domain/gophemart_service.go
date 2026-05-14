package domain

import (
	"context"
)

type GophemartService interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, request UserRegisterRequest) (*AuthResult, error)
	LoginUser(ctx context.Context, request UserLoginRequest) (*AuthResult, error)
	GetUserBalance(ctx context.Context, userID int64) (*UserBalance, error)
}
