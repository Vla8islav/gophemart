package service

import "context"

func (m gophermartService) Ping(ctx context.Context) error {
	return m.repository.Ping(ctx)
}
