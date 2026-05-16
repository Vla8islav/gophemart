package domain

import "time"

type Withdrawal struct {
	ID          int64
	UserID      int64
	Amount      int64
	OrderNumber string
	ProcessedAt time.Time
}

type WithdrawResponseItem struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"` // should marshal to RFC3339 by default
}

type WithdrawResponse []WithdrawResponseItem

type UserBalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type UserBalanceWithdraw struct {
	Order string
	Sum   int64
}
