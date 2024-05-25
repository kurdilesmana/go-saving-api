package accountPort

import (
	"context"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, input accountModel.Account) (err error)
	CheckAccountExist(ctx context.Context, phoneNumber string) (exist bool, err error)
	GetByAccountNo(ctx context.Context, accountNo string) (account *accountModel.Account, err error)
	UpdateBalanceAccount(ctx context.Context, accountNo string, amount float64) (balance float64, err error)
}
