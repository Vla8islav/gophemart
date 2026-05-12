package domain

import (
	"context"

	"github.com/Vla8islav/gophemart/internal/models"
)

type GophemartRepository interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user models.User) (int64, error)
}
