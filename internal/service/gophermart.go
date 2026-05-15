package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vla8islav/gophemart/internal/domain"
)

type gophermartService struct {
	repository      domain.GophermartRepository
	authSecret      []byte
	accrualClient   domain.GophermartAccrualClient
	pollingInterval int
}

func NewMetricsService(
	repo domain.GophermartRepository,
	accrualClient domain.GophermartAccrualClient,
	authSecret string,
	pollingInterval int,
) domain.GophemartService {
	return gophermartService{
		repository:      repo,
		authSecret:      []byte(authSecret),
		accrualClient:   accrualClient,
		pollingInterval: pollingInterval,
	}
}

func (g gophermartService) StartAccrualPolling(ctx context.Context) error {
	if g.pollingInterval <= 0 {
		return fmt.Errorf("polling interval must be positive")
	}

	pollingInterval := time.Duration(g.pollingInterval) * time.Second

	go func() {
		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				pollCtx, cancel := context.WithTimeout(ctx, pollingInterval)
				if err := g.PollAccrual(pollCtx); err != nil {
					log.Printf("failed to update orders from accrual system: %v", err)
				}
				cancel()
			}
		}
	}()

	return nil
}
