package accountService

import (
	"context"
	"fmt"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/accountPort"
	"github.com/kurdilesmana/go-saving-api/pkg/convert"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type accountService struct {
	AccountRepo accountPort.IAccountRepository
	log         *logging.Logger
}

func NewAccountService(
	accountRepo accountPort.IAccountRepository,
	logger *logging.Logger,
) accountPort.IAccountService {
	return &accountService{
		AccountRepo: accountRepo,
		log:         logger,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, request accountModel.CreateAccountRequest) (response accountModel.CreateAccountResponse, err error) {
	// get request id
	requestID := mid.GetIDx(ctx)
	defer s.log.Info(logrus.Fields{"request_id": requestID}, request, "process service completed..")

	// format phoneno
	phoneNumber := convert.NormalizePhoneNumber(request.PhoneNumber)

	// validate account exist
	isExist, err := s.AccountRepo.CheckAccountExist(ctx, phoneNumber)
	if err != nil {
		err = fmt.Errorf("failed to check exist rekening")
		s.log.Warn(logrus.Fields{}, nil, err.Error())
		return
	}

	if isExist {
		err = fmt.Errorf("failed to create, rekening with phone number %s already exist", request.PhoneNumber)
		s.log.Warn(logrus.Fields{}, nil, err.Error())
		return
	}

	// Generate an account number
	AccountNumber := phoneNumber + "0001"

	// Hash PIN
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(request.PIN), bcrypt.MinCost)
	if err != nil {
		err = fmt.Errorf("failed to hashed PIN")
		s.log.Warn(logrus.Fields{}, nil, err.Error())
		return
	}

	// create account
	account := accountModel.Account{
		AccountNo:   AccountNumber,
		Name:        request.Name,
		PhoneNumber: phoneNumber,
		PIN:         string(hashedPIN),
		Balance:     0.0,
	}
	s.log.Info(logrus.Fields{"request_id": requestID}, account, "Create account to db")
	if err = s.AccountRepo.CreateAccount(ctx, account); err != nil {
		err = fmt.Errorf("failed to create rekening")
		s.log.Warn(logrus.Fields{}, nil, err.Error())
		return
	}

	// response
	response = accountModel.CreateAccountResponse{
		AccountNo: account.AccountNo,
		Balance:   account.Balance,
	}

	return
}
