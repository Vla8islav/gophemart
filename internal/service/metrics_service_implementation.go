package service

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/domain"
)

type metricsService struct {
	repository domain.GophemartRepository
	authSecret []byte
}

func NewMetricsService(repo domain.GophemartRepository, authSecret string) domain.GophemartService {
	return metricsService{
		repository: repo,
		authSecret: []byte(authSecret),
	}
}

func (m metricsService) Ping(ctx context.Context) error {
	return m.repository.Ping(ctx)
}
