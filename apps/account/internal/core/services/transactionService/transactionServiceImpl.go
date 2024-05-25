package transactionService

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/accountPort"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/cachePort"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/sirupsen/logrus"
)

type transactionService struct {
	AccountRepo accountPort.IAccountRepository
	CacheRepo   cachePort.ICacheRepository
	log         *logging.Logger
}

func NewTransactionService(
	accountRepo accountPort.IAccountRepository,
	cacheRepo cachePort.ICacheRepository,
	logger *logging.Logger,
) transactionPort.ITransactionService {
	return &transactionService{
		AccountRepo: accountRepo,
		CacheRepo:   cacheRepo,
		log:         logger,
	}
}

func (s *transactionService) Saving(ctx context.Context, accountSession *accountModel.AccountPayload, request transactionModel.SavingRequest) (response transactionModel.TransactionResponse, err error) {
	// get request id
	requestID := mid.GetIDx(ctx)
	defer s.log.Info(logrus.Fields{"request_id": requestID}, request, "process service completed..")

	// get account info
	s.log.Info(logrus.Fields{}, accountSession, "get update info account")
	account, err := s.AccountRepo.GetByAccountNo(ctx, accountSession.AccountNo)
	if err != nil {
		s.log.Warn(logrus.Fields{}, accountSession, err.Error())
		err = fmt.Errorf("failed to saving, error get update info account")
		return
	}

	if account.AccountNo == "" {
		s.log.Warn(logrus.Fields{}, accountSession, "accout not exist")
		err = fmt.Errorf("failed to saving, account not exist")
		return
	}

	// publish to redis
	pubData := map[string]interface{}{
		"transaction_time": time.Now().Format("2006-01-02 15:04:05"),
		"transaction_code": "SAV",
		"account_number":   account.AccountNo,
		"amount":           request.Amount,
	}
	err = s.CacheRepo.PublishStreamEvent(ctx, "savings", pubData)
	if err != nil {
		s.log.Warn(logrus.Fields{}, pubData, err.Error())
		err = fmt.Errorf("failed to saving, error publish data to redis")
		return
	}

	// Waiting to cache done
	time.Sleep(1 * time.Second)
	mutationAmount, _ := s.CacheRepo.GetValue(ctx, requestID)
	s.log.Info(logrus.Fields{"mutation_amount": mutationAmount}, nil, "cache amount")

	// update balance account
	mutAmount, _ := strconv.ParseFloat(strings.TrimSpace(mutationAmount), 64)
	balance, err := s.AccountRepo.UpdateBalanceAccount(ctx, account.AccountNo, mutAmount)
	if err != nil {
		s.log.Warn(logrus.Fields{}, nil, err.Error())
		err = fmt.Errorf("failed to update balance account")
		return
	}

	// response
	response.AccountNo = account.AccountNo
	response.Balance = balance

	return
}

func (s *transactionService) CashWithdrawl(ctx context.Context, accountSession *accountModel.AccountPayload, request transactionModel.CashWithdrawlRequest) (response transactionModel.TransactionResponse, err error) {
	// get request id
	requestID := mid.GetIDx(ctx)
	defer s.log.Info(logrus.Fields{"request_id": requestID}, request, "process service completed..")

	// get account info
	s.log.Info(logrus.Fields{}, accountSession, "get update info account")
	account, err := s.AccountRepo.GetByAccountNo(ctx, accountSession.AccountNo)
	if err != nil {
		s.log.Warn(logrus.Fields{}, accountSession, err.Error())
		err = fmt.Errorf("failed to saving, error get update info account")
		return
	}

	if account.AccountNo == "" {
		s.log.Warn(logrus.Fields{}, accountSession, "account not exist")
		err = fmt.Errorf("failed to saving, account not exist")
		return
	}

	if (account.Balance - request.Amount) < 0.01 {
		s.log.Warn(logrus.Fields{}, accountSession, "insufficient balance")
		err = fmt.Errorf("failed to cashwithdrawl, insufficient balance")
		return
	}

	// publish to redis
	pubData := map[string]interface{}{
		"transaction_time": time.Now().Format("2006-01-02 15:04:05"),
		"transaction_code": "CWD",
		"account_number":   account.AccountNo,
		"amount":           request.Amount,
	}
	err = s.CacheRepo.PublishStreamEvent(ctx, "cashwithdrawl", pubData)
	if err != nil {
		s.log.Warn(logrus.Fields{}, pubData, err.Error())
		err = fmt.Errorf("failed to saving, error publish data to redis")
		return
	}

	// Waiting to cache done
	time.Sleep(1 * time.Second)
	mutationAmount, _ := s.CacheRepo.GetValue(ctx, requestID)
	s.log.Info(logrus.Fields{"mutation_amount": mutationAmount}, nil, "cache amount")

	// update balance account
	mutAmount, _ := strconv.ParseFloat(strings.TrimSpace(mutationAmount), 64)
	balance, err := s.AccountRepo.UpdateBalanceAccount(ctx, account.AccountNo, mutAmount)
	if err != nil {
		s.log.Warn(logrus.Fields{}, nil, err.Error())
		err = fmt.Errorf("failed to update balance account")
		return
	}

	// response
	response.AccountNo = account.AccountNo
	response.Balance = balance

	return
}
