package service

import (
	"context"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
)

type metricsService struct {
	repository domain.GophemartRepository
}

func NewMetricsService(repo domain.GophemartRepository) domain.GophemartService {
	return metricsService{repository: repo}
}

func (m metricsService) Ping(ctx context.Context) error {
	return m.repository.Ping(ctx)
}

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
