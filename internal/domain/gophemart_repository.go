package domain

import (
	"context"
)

type GophemartRepository interface {
	Ping(ctx context.Context) error
}
