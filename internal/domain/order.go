package domain

import "time"

type UserOrderGetResponseItem struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserOrderGetResponse []UserOrderGetResponseItem

type UserOrder struct {
	Number     string
	Status     string
	Accrual    *int64
	UploadedAt time.Time
}

type UpdateOrderParams struct {
	Number  string
	Status  string
	Accrual *int64
}
