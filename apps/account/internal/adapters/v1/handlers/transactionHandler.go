package handlers

import (
	"net/http"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/transactionPort"
	"github.com/kurdilesmana/go-saving-api/apps/account/server/middlewares"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"
	"github.com/kurdilesmana/go-saving-api/pkg/web"
	"github.com/labstack/echo/v4"
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

// SavingHandler godoc
// @Summary API Saving Account
// @Description Saving For Account
// @Tags		Account
// @Accept	json
// @Produce	json
// @Param Authorization header string true "Authorization"  example(Basic NjI4MjE2NzU1MzYzNDQwMDAxOjEyMzQ1Ng==)
// @Param request body transactionModel.SavingRequest true "Request Parameters"
// @Success 200 {object} transactionModel.TransactionResponse "Response Success"
// @Router		/tabung	[post]
func (h *TransactionHandler) SavingHandler(ctx echo.Context) error {
	requestID := mid.GetID(ctx)
	userCtx := mid.SetIDx(ctx.Request().Context(), requestID)

	accountSession, err := middlewares.GetDataFromSessionContext(ctx)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, nil, "error retrieving user from Session")
		return web.ResponseFormatter(ctx, http.StatusUnauthorized, "Unauthorized", nil, err)
	}

	var payload transactionModel.SavingRequest
	if err := ctx.Bind(&payload); err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "saving payload")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, "Bad Request", nil, err)
	}

	// Validate slice payload
	err = h.Validator.Validate(payload)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "saving payload validation")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, "bad request", nil, err)
	}

	respData, err := h.TransactionService.Saving(userCtx, accountSession, payload)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "error saving")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, err.Error(), nil, err)
	}

	return web.ResponseFormatter(ctx, http.StatusOK, "Success", respData, nil)
}

// CashWithdrawlHandler godoc
// @Summary API CashWithdrawl Account
// @Description CashWithdrawl For Account
// @Tags		Account
// @Accept	json
// @Produce	json
// @Param Authorization header string true "Authorization"  example(Basic NjI4MjE2NzU1MzYzNDQwMDAxOjEyMzQ1Ng==)
// @Param request body transactionModel.CashWithdrawlRequest true "Request Parameters"
// @Success 200 {object} transactionModel.TransactionResponse "Response Success"
// @Router		/tarik	[post]
func (h *TransactionHandler) CashWithdrawlHandler(ctx echo.Context) error {
	requestID := mid.GetID(ctx)
	userCtx := mid.SetIDx(ctx.Request().Context(), requestID)

	accountSession, err := middlewares.GetDataFromSessionContext(ctx)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, nil, "error retrieving user from Session")
		return web.ResponseFormatter(ctx, http.StatusUnauthorized, "Unauthorized", nil, err)
	}

	var payload transactionModel.CashWithdrawlRequest
	if err := ctx.Bind(&payload); err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "saving payload")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, "Bad Request", nil, err)
	}

	// Validate slice payload
	err = h.Validator.Validate(payload)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "saving payload validation")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, "bad request", nil, err)
	}

	respData, err := h.TransactionService.CashWithdrawl(userCtx, accountSession, payload)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "error saving")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, err.Error(), nil, err)
	}

	return web.ResponseFormatter(ctx, http.StatusOK, "Success", respData, nil)
}
