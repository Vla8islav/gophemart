package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/helpers"
)

func (m gophermartService) PollAccrual(ctx context.Context) error {
	orders, err := m.repository.GetActiveOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active orders: %w", err)
	}

	infos, pollErr := m.pollOrdersInfo(ctx, orders)

	if err := m.updateOrdersFromAccrualInfo(ctx, infos); err != nil {
		return errors.Join(pollErr, err)
	}

	return pollErr
}

func (m gophermartService) pollOrdersInfo(
	ctx context.Context,
	orders []domain.UserOrder,
) ([]domain.AccrualOrderInfoResponse, error) {
	var pollingErrors []error
	var accrualResponses []domain.AccrualOrderInfoResponse

	for _, order := range orders {
		if err := ctx.Err(); err != nil {
			pollingErrors = append(pollingErrors, err)
			return accrualResponses, errors.Join(pollingErrors...)
		}
		info, err := m.accrualClient.GetOrderInfo(ctx, order.Number)
		if err != nil {
			pollingErrors = append(pollingErrors, fmt.Errorf("order %s: %w", order.Number, err))
			continue
		}

		if info == nil {
			pollingErrors = append(
				pollingErrors,
				fmt.Errorf("order %s: accrual client returned nil info without error", order.Number),
			)
			continue
		}

		accrualResponses = append(accrualResponses, *info)
	}

	return accrualResponses, errors.Join(pollingErrors...)
}

func (m gophermartService) updateOrdersFromAccrualInfo(
	ctx context.Context,
	infos []domain.AccrualOrderInfoResponse,
) error {
	updateParams := make([]domain.UpdateOrderParams, 0, len(infos))

	for _, info := range infos {
		param := domain.UpdateOrderParams{
			Number: info.Order,
			Status: info.Status,
		}

		if info.Accrual != nil {
			accrual := helpers.FloatToCents(*info.Accrual)
			param.Accrual = &accrual
		}

		updateParams = append(updateParams, param)
	}

	if len(updateParams) == 0 {
		return nil
	}

	if err := m.repository.UpdateOrders(ctx, updateParams); err != nil {
		return fmt.Errorf("failed to update orders from accrual info: %w", err)
	}

	return nil
}
