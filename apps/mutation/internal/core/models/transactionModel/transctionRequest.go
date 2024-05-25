package transactionModel

type SavingSubs struct {
	AccountNumber   string  `json:"account_number" validate:"required"`
	TransactionTime string  `json:"transaction_time" validate:"required"`
	TransactionCode string  `json:"transaction_code" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,min=1"`
}

type CashWithdrawlSubs struct {
	AccountNumber   string  `json:"account_number" validate:"required"`
	TransactionTime string  `json:"transaction_time" validate:"required"`
	TransactionCode string  `json:"transaction_code" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,min=1"`
}
