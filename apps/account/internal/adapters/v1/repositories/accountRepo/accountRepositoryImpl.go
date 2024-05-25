package accountRepo

import (
	"context"
	"fmt"
	"time"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/accountPort"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accountRepository struct {
	DB             *gorm.DB
	KeyTransaction string
	timeout        time.Duration
	log            *logging.Logger
}

func NewaccountRepo(db *gorm.DB, keyTransaction string, timeout int, logger *logging.Logger) accountPort.IAccountRepository {
	return &accountRepository{
		DB:             db,
		KeyTransaction: keyTransaction,
		timeout:        time.Duration(timeout) * time.Second,
		log:            logger,
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, input accountModel.Account) (err error) {
	r.log.Info(logrus.Fields{}, input, "start create account...")
	trx, ok := ctx.Value(r.KeyTransaction).(*gorm.DB)
	if !ok {
		trx = r.DB
	}

	ctxWT, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res := trx.WithContext(ctxWT).Create(&input)
	if res.Error != nil {
		remark := "failed to create rekening"
		r.log.Error(logrus.Fields{
			"error": res.Error.Error(),
		}, input, remark)
		err = fmt.Errorf(remark)
		return
	}

	return
}

func (r *accountRepository) CheckAccountExist(ctx context.Context, PhoneNumber string) (exist bool, err error) {
	r.log.Info(logrus.Fields{"PhoneNumber": PhoneNumber}, nil, "check account exist...")
	trx, ok := ctx.Value(r.KeyTransaction).(*gorm.DB)
	if !ok {
		trx = r.DB
	}

	ctxWT, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var account accountModel.Account
	res := trx.WithContext(ctxWT).Where(accountModel.Account{PhoneNumber: PhoneNumber}).Find(&account)
	if res.Error != nil {
		remark := "failed to get account by phonenumber"
		r.log.Error(logrus.Fields{
			"error": res.Error.Error(),
		}, nil, remark)
		err = fmt.Errorf(remark)
		return
	}

	if res.RowsAffected > 0 {
		exist = true
	}

	return
}

func (r *accountRepository) GetByAccountNo(ctx context.Context, accountNo string) (*accountModel.Account, error) {
	r.log.Info(logrus.Fields{"accountNo": accountNo}, nil, "get account by account no...")
	trx, ok := ctx.Value(r.KeyTransaction).(*gorm.DB)
	if !ok {
		trx = r.DB
	}

	ctxWT, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var account accountModel.Account
	res := trx.WithContext(ctxWT).Where("account_no = ?", accountNo).Find(&account)
	if res.Error != nil {
		remark := "failed to get account by account no"
		r.log.Error(logrus.Fields{
			"error": res.Error.Error(),
		}, nil, remark)
		err := fmt.Errorf(remark)
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) UpdateBalanceAccount(ctx context.Context, accountNo string, amount float64) (balance float64, err error) {
	r.log.Info(logrus.Fields{"accountNo": accountNo}, nil, "get account by account no...")
	trx, ok := ctx.Value(r.KeyTransaction).(*gorm.DB)
	if !ok {
		trx = r.DB
	}

	ctxWT, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var account accountModel.Account
	res := trx.WithContext(ctxWT).Model(&account).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "balance"}}}).
		Where("account_no = ?", accountNo).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount))

	if res.Error != nil {
		remark := "failed to update balance account"
		r.log.Error(logrus.Fields{
			"error": res.Error.Error(),
		}, nil, remark)
		err = fmt.Errorf(remark)
		return
	}
	balance = account.Balance

	return
}
