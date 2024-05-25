package transactionPort

import (
	"context"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/transactionModel"
)

type ITransactionService interface {
	Saving(ctx context.Context, accountSession *accountModel.AccountPayload, request transactionModel.SavingRequest) (response transactionModel.TransactionResponse, err error)
	CashWithdrawl(ctx context.Context, accountSession *accountModel.AccountPayload, request transactionModel.CashWithdrawlRequest) (response transactionModel.TransactionResponse, err error)
}
