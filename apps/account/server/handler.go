package server

import (
	"github.com/kurdilesmana/go-saving-api/apps/account/deps"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"

	handler "github.com/kurdilesmana/go-saving-api/apps/account/internal/adapters/v1/handlers"
)

type Handler struct {
	accountHandler     handler.AccountHandler
	transactionHandler handler.TransactionHandler
}

func SetupHandler(dep deps.Dependency) Handler {
	//init validator
	validator := validator.NewRequestValidator()

	return Handler{
		accountHandler:     handler.NewAccountHandler(dep.AccountService, dep.Logger, validator),
		transactionHandler: handler.NewTransactionHandler(dep.TransactionService, dep.Logger, validator),
	}
}
