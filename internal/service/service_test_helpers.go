package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

type fakeCreateUserRepository struct {
	createUserFunc     func(ctx context.Context, params domain.CreateUserParams) (int64, error)
	getUserByLoginFunc func(ctx context.Context, login string) (*domain.User, error)
}

func (r fakeCreateUserRepository) Ping(ctx context.Context) error {
	return nil
}

func (r fakeCreateUserRepository) CreateUser(ctx context.Context, params domain.CreateUserParams) (int64, error) {
	return r.createUserFunc(ctx, params)
}

func (r fakeCreateUserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	return r.getUserByLoginFunc(ctx, login)
}
