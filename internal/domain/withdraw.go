package domain

type UserBalanceWithdrawRequest struct {
	Order string `json:"order"`
	Sum   int64  `json:"sum"`
}
