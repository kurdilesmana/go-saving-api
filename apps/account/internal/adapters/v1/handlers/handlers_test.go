package handlers_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/adapters/v1/handlers"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/transactionModel"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BaseResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) CreateAccount(ctx context.Context, u accountModel.CreateAccountRequest) (accountModel.CreateAccountResponse, error) {
	args := m.Called(ctx, u)
	return args[0].(accountModel.CreateAccountResponse), args.Error(1)
}

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) Saving(ctx context.Context, accountSession *accountModel.AccountPayload, request transactionModel.SavingRequest) (
	transactionModel.TransactionResponse, error) {
	args := m.Called(ctx, accountSession, request)
	return args[0].(transactionModel.TransactionResponse), args.Error(1)
}

func (m *MockTransactionService) CashWithdrawl(ctx context.Context, accountSession *accountModel.AccountPayload, request transactionModel.CashWithdrawlRequest) (
	transactionModel.TransactionResponse, error) {
	args := m.Called(ctx, accountSession, request)
	return args[0].(transactionModel.TransactionResponse), args.Error(1)
}

func basicAuthEncode(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestCreateAccount(t *testing.T) {
	// Setup
	mockAccountService := new(MockAccountService)
	user := accountModel.CreateAccountRequest{
		Name:        "John Doe",
		PhoneNumber: "082123567890",
		PIN:         "12356",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	userResp := accountModel.CreateAccountResponse{
		AccountNo: "082123567890001",
		Balance:   0,
	}
	userResJSON, err := json.Marshal(userResp)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/daftar", bytes.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("request-id", "test-create-account")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	requestID := mid.GetID(c)
	ctx := mid.SetIDx(c.Request().Context(), requestID)
	mockAccountService.On("CreateAccount", ctx, user).Return(userResp, nil)

	log := logging.NewLogger("TEST_APP")
	validator := validator.NewRequestValidator()
	h := &handlers.AccountHandler{AccountService: mockAccountService, Logger: log, Validator: validator}

	// Run the test
	err = h.CreateHandler(c)

	// Assertions
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, http.StatusCreated, rec.Code)
	var responseUser accountModel.CreateAccountResponse
	err = json.Unmarshal(userResJSON, &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, userResp, responseUser)

	// Verify mock expectations
	mockAccountService.AssertExpectations(t)
}

func TestSaving(t *testing.T) {
	// Setup
	mockTransactionService := new(MockTransactionService)
	expectedAccount := accountModel.AccountPayload{
		AccountNo:   "082123567890001",
		Name:        "John Doe",
		PhoneNumber: "082123567890",
		Balance:     0,
	}

	username := "082123567890001"
	password := "12356"
	expectedAuthHeader := "Basic " + basicAuthEncode(username, password)

	tabung := transactionModel.SavingRequest{Amount: 100}
	tabungJSON, err := json.Marshal(tabung)
	if err != nil {
		t.Fatal(err)
	}

	tabungResp := transactionModel.TransactionResponse{
		AccountNo: "082123567890001",
		Balance:   100,
	}

	tabungRespJSON, err := json.Marshal(tabungResp)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Request tabung nomor rekening dikenali", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tabung", bytes.NewReader(tabungJSON))
		req.Header.Set(echo.HeaderAuthorization, expectedAuthHeader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("request-id", "test-saving-transaction-ok")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set session valid
		c.Set("session", &expectedAccount)

		requestID := mid.GetID(c)
		ctx := mid.SetIDx(c.Request().Context(), requestID)
		mockTransactionService.On("Saving", ctx, &expectedAccount, tabung).Return(tabungResp, nil)

		log := logging.NewLogger("TEST_APP")
		validator := validator.NewRequestValidator()
		h := &handlers.TransactionHandler{TransactionService: mockTransactionService, Logger: log, Validator: validator}

		// Run the test
		err = h.SavingHandler(c)

		// Assertions
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		var responseTabung transactionModel.TransactionResponse
		err = json.Unmarshal(tabungRespJSON, &responseTabung)
		assert.NoError(t, err)
		assert.Equal(t, tabungResp, responseTabung)

		// Verify mock expectations
		mockTransactionService.AssertExpectations(t)
	})

	t.Run("Request tabung nomor rekening tidak dikenali", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tabung", bytes.NewReader(tabungJSON))
		req.Header.Set(echo.HeaderAuthorization, expectedAuthHeader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("request-id", "test-saving-transaction-ok")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		requestID := mid.GetID(c)
		ctx := mid.SetIDx(c.Request().Context(), requestID)
		mockTransactionService.On("Saving", ctx, &expectedAccount, tabung).Return(tabungResp, nil)

		log := logging.NewLogger("TEST_APP")
		validator := validator.NewRequestValidator()
		h := &handlers.TransactionHandler{TransactionService: mockTransactionService, Logger: log, Validator: validator}

		// Run the test
		h.SavingHandler(c)

		// Expected unauthorized status code
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestCashWithdrawlOK(t *testing.T) {
	// Setup
	mockTransactionService := new(MockTransactionService)
	expectedAccount := accountModel.AccountPayload{
		AccountNo:   "082123567890001",
		Name:        "John Doe",
		PhoneNumber: "082123567890",
		Balance:     0,
	}

	username := "082123567890001"
	password := "12356"
	expectedAuthHeader := "Basic " + basicAuthEncode(username, password)

	tarik := transactionModel.CashWithdrawlRequest{Amount: 100}
	tarikJSON, err := json.Marshal(tarik)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Request tarik nomor rekening dikenali dan saldo cukup", func(t *testing.T) {
		tarikResp := transactionModel.TransactionResponse{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tarik", bytes.NewReader(tarikJSON))
		req.Header.Set(echo.HeaderAuthorization, expectedAuthHeader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("request-id", "test-cashwithdrawl-transaction-ok")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// set session valid
		c.Set("session", &expectedAccount)

		requestID := mid.GetID(c)
		ctx := mid.SetIDx(c.Request().Context(), requestID)
		mockTransactionService.On("CashWithdrawl", ctx, &expectedAccount, tarik).Return(tarikResp, nil)

		log := logging.NewLogger("TEST_APP")
		validator := validator.NewRequestValidator()
		h := &handlers.TransactionHandler{TransactionService: mockTransactionService, Logger: log, Validator: validator}

		// Run the test
		err = h.CashWithdrawlHandler(c)

		// Assertions
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, http.StatusOK, rec.Code)
		var responsetarik transactionModel.TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &responsetarik)
		assert.NoError(t, err)
		assert.Equal(t, tarikResp, responsetarik)

		// Verify mock expectations
		mockTransactionService.AssertExpectations(t)
	})

	t.Run("Request tarik nomor rekening tidak dikenali", func(t *testing.T) {
		tarikResp := transactionModel.TransactionResponse{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tarik", bytes.NewReader(tarikJSON))
		req.Header.Set(echo.HeaderAuthorization, expectedAuthHeader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("request-id", "test-cashwithdrawl-transaction-ok")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		requestID := mid.GetID(c)
		ctx := mid.SetIDx(c.Request().Context(), requestID)
		mockTransactionService.On("CashWithdrawl", ctx, &expectedAccount, tarik).Return(tarikResp, nil)

		log := logging.NewLogger("TEST_APP")
		validator := validator.NewRequestValidator()
		h := &handlers.TransactionHandler{TransactionService: mockTransactionService, Logger: log, Validator: validator}

		// Run the test
		h.CashWithdrawlHandler(c)

		// Expected unauthorized status code
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestCashWithdrawlErrBalance(t *testing.T) {
	// Setup
	mockTransactionService := new(MockTransactionService)
	expectedAccount := accountModel.AccountPayload{
		AccountNo:   "082123567890001",
		Name:        "John Doe",
		PhoneNumber: "082123567890",
		Balance:     0,
	}

	username := "082123567890001"
	password := "12356"
	expectedAuthHeader := "Basic " + basicAuthEncode(username, password)

	tarik := transactionModel.CashWithdrawlRequest{Amount: 100}
	tarikJSON, err := json.Marshal(tarik)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Request tarik nomor rekening dikenali dan saldo tidak cukup", func(t *testing.T) {
		var resptarik transactionModel.TransactionResponse
		errBalance := errors.New("failed to cashwithdrawl, insufficient balance")
		tarikRespErr := BaseResponse{
			Message: errBalance.Error(),
			Error:   errBalance.Error(),
			Data:    nil,
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tarik", bytes.NewReader(tarikJSON))
		req.Header.Set(echo.HeaderAuthorization, expectedAuthHeader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("request-id", "test-cashwithdrawl-transaction-failed")
		recErr := httptest.NewRecorder()
		c := e.NewContext(req, recErr)

		// set session valid
		c.Set("session", &expectedAccount)

		requestID := mid.GetID(c)
		ctx := mid.SetIDx(c.Request().Context(), requestID)
		mockTransactionService.On("CashWithdrawl", ctx, &expectedAccount, tarik).Return(resptarik, errBalance)

		log := logging.NewLogger("TEST_APP")
		validator := validator.NewRequestValidator()
		h := &handlers.TransactionHandler{TransactionService: mockTransactionService, Logger: log, Validator: validator}

		// Run the test
		h.CashWithdrawlHandler(c)

		assert.Equal(t, http.StatusBadRequest, recErr.Code)
		var respErrTarik BaseResponse
		err = json.Unmarshal(recErr.Body.Bytes(), &respErrTarik)
		assert.Nil(t, err)
		assert.Equal(t, tarikRespErr, respErrTarik)
	})
}
