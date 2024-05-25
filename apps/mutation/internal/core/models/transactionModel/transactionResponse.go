package transactionModel

type TransactionResponse struct {
	AccountNo string  `json:"account_no" validate:"required"`
	Balance   float64 `json:"balance" validate:"required"`
}
