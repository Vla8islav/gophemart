package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
)

func (m gophermartService) CreateOrder(ctx context.Context, userID int64, orderNumber string) error {
	orderNumber = strings.TrimSpace(orderNumber)
	if orderNumber == "" {
		return domain.ErrOrderEmpty
	}

	if !helpers.ValidLuhn(orderNumber) {
		return domain.ErrInvalidOrderNumber
	}

	err := m.repository.CreateOrder(ctx, userID, orderNumber)
	if err != nil {
		return fmt.Errorf("failed to CreateOrder: %w", err)
	}

	return nil
}
