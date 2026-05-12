package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

type fakeCreateUserRepository struct {
	createUserFunc func(ctx context.Context, params domain.CreateUserParams) (int64, error)
}

func (r fakeCreateUserRepository) Ping(ctx context.Context) error {
	return nil
}

func (r fakeCreateUserRepository) CreateUser(ctx context.Context, params domain.CreateUserParams) (int64, error) {
	return r.createUserFunc(ctx, params)
}
