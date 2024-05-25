package transactionModel

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionTime time.Time `json:"transaction_time"`
	TransactionCode string    `json:"transaction_code"`
	Amount          float64   `json:"amount"`
}

type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey"`
	TransactionID int     `json:"transaction_id"`
	Mutation      string  `json:"mutation"`
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
}
