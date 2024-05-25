package transactionModel

type SavingRequest struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
}

type CashWithdrawlRequest struct {
	Amount float64 `json:"amount" validate:"required,min=1"`
}
