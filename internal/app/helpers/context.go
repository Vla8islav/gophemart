package helpers

import (
	"context"
	"time"
)

func GetDefaultContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	return ctx, cancel
}
