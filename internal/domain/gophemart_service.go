package domain

import "context"

type GophemartService interface {
	Ping(ctx context.Context) error
}
