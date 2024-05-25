package transactionRepo

import (
	"context"
	"fmt"
	"time"

	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type mutationRepository struct {
	DB             *gorm.DB
	KeyTransaction string
	timeout        time.Duration
	log            *logging.Logger
}

func NewMutationRepo(db *gorm.DB, keyTransaction string, timeout int, logger *logging.Logger) transactionPort.ITransactionRepository {
	return &mutationRepository{
		DB:             db,
		KeyTransaction: keyTransaction,
		timeout:        time.Duration(timeout) * time.Second,
		log:            logger,
	}
}

func (r *mutationRepository) Create(ctx context.Context, transaction *transactionModel.Transaction) (ID int, err error) {
	r.log.Info(logrus.Fields{}, transaction, "start create transaction...")
	trx, ok := ctx.Value(r.KeyTransaction).(*gorm.DB)
	if !ok {
		trx = r.DB
	}

	ctxWT, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res := trx.WithContext(ctxWT).Create(&transaction)
	if res.Error != nil {
		remark := "failed to create transaction"
		r.log.Error(logrus.Fields{
			"error": res.Error.Error(),
		}, transaction, remark)
		err = fmt.Errorf(remark)
		return
	}

	ID = int(transaction.ID)

	return
}

func (r *mutationRepository) CreateDetail(ctx context.Context, transactionDetail *transactionModel.TransactionDetail) (err error) {
	r.log.Info(logrus.Fields{}, transactionDetail, "start create transaction detail...")
	trx, ok := ctx.Value(r.KeyTransaction).(*gorm.DB)
	if !ok {
		trx = r.DB
	}

	ctxWT, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res := trx.WithContext(ctxWT).Create(&transactionDetail)
	if res.Error != nil {
		remark := "failed to create transaction detail"
		r.log.Error(logrus.Fields{
			"error": res.Error.Error(),
		}, transactionDetail, remark)
		err = fmt.Errorf(remark)
		return
	}
	return
}
