package domain

import (
	"context"
)

type GophemartRepository interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user CreateUserParams) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
}
