package middlewares

import (
	"fmt"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/labstack/echo/v4"
)

func GetDataFromSessionContext(ctx echo.Context) (*accountModel.AccountPayload, error) {
	payloadClaims, ok := ctx.Get("session").(*accountModel.AccountPayload)
	if !ok || payloadClaims == nil {
		return nil, fmt.Errorf("invalid session")
	}

	return payloadClaims, nil
}
