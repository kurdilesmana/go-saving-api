package accountModel

type Account struct {
	AccountNo   string  `json:"account_no" gorm:"column:account_no;type:varchar(20);primary_key;"`
	Name        string  `json:"name" gorm:"column:name;type:varchar(50);null;" `
	PhoneNumber string  `json:"phone_number" gorm:"column:phone_number;type:varchar(20);null;" `
	PIN         string  `json:"pin" gorm:"column:pin;type:text;null;" `
	Balance     float64 `json:"balance" gorm:"column:balance;type:float;null;" `
}

func (b *Account) TableName() string {
	return "public.account"
}
