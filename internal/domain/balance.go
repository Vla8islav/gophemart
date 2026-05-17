package domain

type UserBalance struct {
	Current   int64
	Withdrawn int64
}

type UserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
