package accountModel

type AccountPayload struct {
	AccountNo   string  `json:"account_no"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Balance     float64 `json:"balance"`
}

type CreateAccountResponse struct {
	AccountNo string  `json:"account_no" validate:"required"`
	Balance   float64 `json:"balance" validate:"required"`
}
