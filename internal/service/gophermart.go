package service

import (
	"github.com/Vla8islav/gophemart/internal/domain"
)

type gophermartService struct {
	repository    domain.GophermartRepository
	authSecret    []byte
	accrualClient domain.GophermartAccrualClient
}

func NewMetricsService(repo domain.GophermartRepository,
	accrualClient domain.GophermartAccrualClient, authSecret string) domain.GophemartService {
	return gophermartService{
		repository:    repo,
		authSecret:    []byte(authSecret),
		accrualClient: accrualClient,
	}
}
