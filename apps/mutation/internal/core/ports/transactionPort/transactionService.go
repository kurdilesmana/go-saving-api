package transactionPort

import (
	"context"

	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/models/transactionModel"
)

type ITransactionService interface {
	Saving(ctx context.Context, request transactionModel.SavingSubs) (err error)
	CashWithdrawl(ctx context.Context, request transactionModel.CashWithdrawlSubs) (err error)
}
