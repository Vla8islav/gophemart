package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
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
