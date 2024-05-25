package transactionService

import (
	"context"
	"fmt"
	"time"

	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/ports/cachePort"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/pkg/convert"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/sirupsen/logrus"
)

type transactionService struct {
	TransactionRepo transactionPort.ITransactionRepository
	CacheRepo       cachePort.ICacheRepository
	log             *logging.Logger
}

func NewTransactionService(
	transactionRepo transactionPort.ITransactionRepository,
	cacheRepo cachePort.ICacheRepository,
	logger *logging.Logger,
) transactionPort.ITransactionService {
	return &transactionService{
		TransactionRepo: transactionRepo,
		CacheRepo:       cacheRepo,
		log:             logger,
	}
}

func (s *transactionService) Saving(ctx context.Context, request transactionModel.SavingSubs) (err error) {
	// get request id
	eventID := mid.GetIDx(ctx)
	defer s.log.Info(logrus.Fields{"event_id": eventID}, request, "process service completed..")

	// create transaction
	transactionTime, _ := convert.StringToTime(request.TransactionTime, "1")
	transaction := transactionModel.Transaction{
		TransactionTime: transactionTime,
		TransactionCode: request.TransactionCode,
		Amount:          request.Amount,
	}
	IDTransaction, err := s.TransactionRepo.Create(ctx, &transaction)
	if err != nil {
		err = fmt.Errorf("failed to create transaction")
		s.log.Warn(logrus.Fields{}, transaction, err.Error())
		return
	}

	// create transaction detail
	transactionDetail := transactionModel.TransactionDetail{
		TransactionID: IDTransaction,
		Mutation:      "C",
		AccountNumber: request.AccountNumber,
		Amount:        request.Amount,
	}
	if err = s.TransactionRepo.CreateDetail(ctx, &transactionDetail); err != nil {
		s.log.Warn(logrus.Fields{}, transaction, err.Error())
		err = fmt.Errorf("failed to create transaction")
		return
	}

	// set update amount for saldo
	err = s.CacheRepo.SetValue(ctx, eventID, request.Amount, 5*time.Second)
	if err != nil {
		s.log.Warn(logrus.Fields{}, transaction, err.Error())
		err = fmt.Errorf("failed to set mutation to cache")
		return
	}

	return
}

func (s *transactionService) CashWithdrawl(ctx context.Context, request transactionModel.CashWithdrawlSubs) (err error) {
	// get request id
	eventID := mid.GetIDx(ctx)
	defer s.log.Info(logrus.Fields{"event_id": eventID}, request, "process service completed..")

	// create transaction
	transactionTime, _ := convert.StringToTime(request.TransactionTime, "1")
	transaction := transactionModel.Transaction{
		TransactionTime: transactionTime,
		TransactionCode: request.TransactionCode,
		Amount:          request.Amount,
	}
	IDTransaction, err := s.TransactionRepo.Create(ctx, &transaction)
	if err != nil {
		err = fmt.Errorf("failed to create transaction")
		s.log.Warn(logrus.Fields{}, transaction, err.Error())
		return
	}

	// create transaction detail
	transactionDetail := transactionModel.TransactionDetail{
		TransactionID: IDTransaction,
		Mutation:      "D",
		AccountNumber: request.AccountNumber,
		Amount:        request.Amount,
	}
	if err = s.TransactionRepo.CreateDetail(ctx, &transactionDetail); err != nil {
		s.log.Warn(logrus.Fields{}, transaction, err.Error())
		err = fmt.Errorf("failed to create transaction")
		return
	}

	// set update amount for saldo
	err = s.CacheRepo.SetValue(ctx, eventID, (request.Amount * -1), 5*time.Second)
	if err != nil {
		s.log.Warn(logrus.Fields{}, transaction, err.Error())
		err = fmt.Errorf("failed to set mutation to cache")
		return
	}

	return
}
