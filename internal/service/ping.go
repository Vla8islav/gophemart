package service

import "context"

func (m metricsService) Ping(ctx context.Context) error {
	return m.repository.Ping(ctx)
}
