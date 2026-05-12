package service

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
)

func (m metricsService) CreateUser(ctx context.Context, userRegReq domain.UserRegisterRequest) (int64, error) {
	hash, err := helpers.HashPassword(userRegReq.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate the hash for the new user %s: %w", userRegReq.Login, err)
	}
	createUserParams := domain.CreateUserParams{
		Login:        userRegReq.Login,
		PasswordHash: hash,
	}
	return m.repository.CreateUser(ctx, createUserParams)
}
