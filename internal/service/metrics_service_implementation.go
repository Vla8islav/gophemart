package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/models"
)

type metricsService struct {
	repository domain.GophemartRepository
}

func NewMetricsService(repo domain.GophemartRepository) domain.GophemartRepository {
	return metricsService{repository: repo}
}

func (m metricsService) Ping(ctx context.Context) error {
	return m.repository.Ping(ctx)
}

func (m metricsService) CreateUser(ctx context.Context, user models.User) error {
	return m.repository.CreateUser(ctx, user)
}
