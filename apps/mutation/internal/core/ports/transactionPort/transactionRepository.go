package transactionPort

import (
	"context"

	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/models/transactionModel"
)

type ITransactionRepository interface {
	Create(ctx context.Context, transaction *transactionModel.Transaction) (ID int, err error)
	CreateDetail(ctx context.Context, transactionDetail *transactionModel.TransactionDetail) (err error)
}
