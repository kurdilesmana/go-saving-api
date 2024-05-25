package handlers

import (
	"net/http"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/accountPort"
	"github.com/kurdilesmana/go-saving-api/apps/account/server/middlewares"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"
	"github.com/kurdilesmana/go-saving-api/pkg/web"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	AccountService accountPort.IAccountService
	Logger         *logging.Logger
	Validator      *validator.RequestValidator
}

func NewAccountHandler(
	accountService accountPort.IAccountService,
	logger *logging.Logger,
	validator *validator.RequestValidator,
) AccountHandler {
	return AccountHandler{
		AccountService: accountService,
		Logger:         logger,
		Validator:      validator,
	}
}

// CreateHandler godoc
// @Summary API Create Account
// @Description Create For Account
// @Tags		Account
// @Accept	json
// @Produce	json
// @Param CreateAccountRequest body accountModel.CreateAccountRequest true "Request Parameters"
// @Success 200 {object} accountModel.CreateAccountResponse "Response Success"
// @Router		/daftar	[post]
func (h *AccountHandler) CreateHandler(ctx echo.Context) error {
	requestID := mid.GetID(ctx)
	userCtx := mid.SetIDx(ctx.Request().Context(), requestID)

	var payload accountModel.CreateAccountRequest
	if err := ctx.Bind(&payload); err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "create account payload")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, "Bad Request", nil, err)
	}

	// Validate slice payload
	err := h.Validator.Validate(payload)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "create account payload validation")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, "bad request", nil, err)
	}

	respData, err := h.AccountService.CreateAccount(userCtx, payload)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, payload, "error create account")
		return web.ResponseFormatter(ctx, http.StatusBadRequest, err.Error(), nil, err)
	}

	return web.ResponseFormatter(ctx, http.StatusOK, "Success", respData, nil)
}

// InquiryHandler godoc
// @Summary API Inquiry Account
// @Description Inquiry For Account
// @Tags		Account
// @Accept	json
// @Produce	json
// @Param Authorization header string true "Authorization"  example(Basic NjI4MjE2NzU1MzYzNDQwMDAxOjEyMzQ1Ng==)
// @Success 200 {object} accountModel.Account "Response Success"
// @Router		/inquiry	[GET]
func (h *AccountHandler) InquiryHandler(ctx echo.Context) error {
	requestID := mid.GetID(ctx)

	userSession, err := middlewares.GetDataFromSessionContext(ctx)
	if err != nil {
		h.Logger.Error(logrus.Fields{"request_id": requestID, "error": err.Error()}, nil, "error retrieving user from Session")
		return web.ResponseFormatter(ctx, http.StatusUnauthorized, "Unauthorized", nil, err)
	}

	return web.ResponseFormatter(ctx, http.StatusOK, "Success", userSession, nil)
}
