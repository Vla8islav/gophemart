package domain

import (
	"context"
)

type GophemartService interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, request UserRegisterRequest) (*AuthResult, error)
}
