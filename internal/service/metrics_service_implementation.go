package service

import (
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
