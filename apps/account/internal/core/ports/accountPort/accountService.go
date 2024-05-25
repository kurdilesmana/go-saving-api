package accountPort

import (
	"context"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, request accountModel.CreateAccountRequest) (response accountModel.CreateAccountResponse, err error)
}
