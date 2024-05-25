package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/apps/mutation/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	TransactionService transactionPort.ITransactionService
	Logger             *logging.Logger
	Validator          *validator.RequestValidator
}

func NewTransactionHandler(
	transactionService transactionPort.ITransactionService,
	logger *logging.Logger,
	validator *validator.RequestValidator,
) TransactionHandler {
	return TransactionHandler{
		TransactionService: transactionService,
		Logger:             logger,
		Validator:          validator,
	}
}

func (h *TransactionHandler) SavingHandler(ctx context.Context, eventID, eventData string) (err error) {
	h.Logger.Info(logrus.Fields{"event_id": eventID, "event_data": eventData}, nil, "start handler saving...")
	trxCtx := mid.SetIDx(ctx, eventID)

	// bind request
	var reqSub transactionModel.SavingSubs
	err = json.Unmarshal([]byte(eventData), &reqSub)
	if err != nil {
		h.Logger.Error(logrus.Fields{}, nil, err.Error())
		err = fmt.Errorf("failed unmarshal publish data")
		return
	}

	// validate request
	h.Logger.Info(logrus.Fields{"event_id": eventID}, reqSub, "publish data...")
	err = h.Validator.Validate(reqSub)
	if err != nil {
		h.Logger.Error(logrus.Fields{}, reqSub, err.Error())
		err = fmt.Errorf("failed validate publish data")
		return
	}

	// process saving
	err = h.TransactionService.Saving(trxCtx, reqSub)
	if err != nil {
		h.Logger.Error(logrus.Fields{}, reqSub, err.Error())
		err = fmt.Errorf("failed process saving using publish data")
		return
	}

	return nil
}

func (h *TransactionHandler) CashWithdrawlHandler(ctx context.Context, eventID, eventData string) (err error) {
	h.Logger.Info(logrus.Fields{"event_id": eventID, "event_data": eventData}, nil, "start handler saving...")
	trxCtx := mid.SetIDx(ctx, eventID)

	// bind request
	var reqSub transactionModel.CashWithdrawlSubs
	err = json.Unmarshal([]byte(eventData), &reqSub)
	if err != nil {
		h.Logger.Error(logrus.Fields{}, nil, err.Error())
		err = fmt.Errorf("failed unmarshal publish data")
		return
	}

	// validate request
	h.Logger.Info(logrus.Fields{"event_id": eventID}, reqSub, "publish data...")
	err = h.Validator.Validate(reqSub)
	if err != nil {
		h.Logger.Error(logrus.Fields{}, reqSub, err.Error())
		err = fmt.Errorf("failed validate publish data")
		return
	}

	// process saving
	err = h.TransactionService.CashWithdrawl(trxCtx, reqSub)
	if err != nil {
		h.Logger.Error(logrus.Fields{}, reqSub, err.Error())
		err = fmt.Errorf("failed process saving using publish data")
		return
	}

	return nil
}
