package mid

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const requestIDKey = "request-id"
const requestIDHeader = "X-Request-Id"

// RequestID read header with key X-Request-Id, if exist that value used to traceID
// if not, generate uuid for traceID
func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Request().Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set(requestIDKey, requestID) // set to context echo
		return next(c)
	}
}

func GetID(ctx echo.Context) string {
	requestID, ok := ctx.Get(requestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

// Because echo request context is not included value from echo.Context
// we need to build this method, hiks.
type keyCtx string

func SetIDx(ctx context.Context, requsetID string) context.Context {
	return context.WithValue(ctx, keyCtx(requestIDKey), requsetID)
}

func GetIDx(ctx context.Context) string {
	requestID, ok := ctx.Value(keyCtx(requestIDKey)).(string)
	if !ok {
		return ""
	}
	return requestID
}
