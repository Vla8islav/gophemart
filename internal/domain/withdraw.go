package domain

type UserBalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type UserBalanceWithdraw struct {
	Order string
	Sum   int64
}
