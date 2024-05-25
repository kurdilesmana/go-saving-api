package server

import (
	"github.com/kurdilesmana/go-saving-api/apps/mutation/deps"
	"github.com/kurdilesmana/go-saving-api/pkg/validator"

	handler "github.com/kurdilesmana/go-saving-api/apps/mutation/internal/adapters/v1/handlers"
)

type Handler struct {
	transactionHandler handler.TransactionHandler
}

func SetupHandler(dep deps.Dependency) Handler {
	//init validator
	validator := validator.NewRequestValidator()

	return Handler{
		transactionHandler: handler.NewTransactionHandler(dep.TransactionService, dep.Logger, validator),
	}
}
