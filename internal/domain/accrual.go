package domain

import "context"

type AccrualOrderInfoResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type GophermartAccrualClient interface {
	GetOrderInfo(ctx context.Context, orderNumber string) (*AccrualOrderInfoResponse, error)
}
