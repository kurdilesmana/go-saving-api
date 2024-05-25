package mid

import (
	"context"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)

type MiddlewareTimeout struct {
	// Maximum processing time limit (milliseconds).
	Timeout int

	// Function to handle response when timeout has been exceeded, default JSON format.
	// example value: {"code": 504, "message": "processing time limit has been exceeded"}
	ErrorHandlerFunc func(code int, response interface{}) error

	// HTTP status code for response, default 504 - gateway timeout.
	StatusCode int

	// Response body when timeout has been exceeded
	Response interface{}
}

// HandlerFunc method for manage max limit timeout request.
func (m *MiddlewareTimeout) HandlerFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Instance variable to prepare middleware function
		var (
			timeout         = time.Duration(m.Timeout)
			errorHandleFunc = m.ErrorHandlerFunc
			statusCode      = m.StatusCode
			response        = m.Response
		)

		// Take maximum request limit from header, header-timeout milliseconds.
		if t, err := time.ParseDuration(c.Request().Header.Get("header-timeout")); err == nil && t != 0 {
			timeout = t
		}

		// Set default value when empty, 0 and nill
		if timeout == 0 {
			timeout = 30000 // Milliseconds 30 s
		}
		if errorHandleFunc == nil {
			errorHandleFunc = func(code int, response interface{}) error {
				return c.JSON(code, response)
			} // Response JSON format
		}
		if statusCode == 0 {
			statusCode = http.StatusGatewayTimeout // 504 - Gateway Timeout
		}
		if response == nil {
			response = echo.Map{"code": http.StatusGatewayTimeout, "message": "processing time limit has been exceeded"}
		}

		// Create context to control request time limit. Parent context from original request.
		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Millisecond*timeout)
		defer cancel()

		// Instance the channel and create a goroutine to continue the request and
		// wait for the processing results (sent through the channel).
		done := make(chan error, 1)
		go func() {
			// Recover when next(c) error
			defer func() {
				if r := recover(); r != nil {
					done <- c.JSON(http.StatusInternalServerError, echo.Map{"message": r})
				}
			}()

			// Overwrite request context
			c.SetRequest(c.Request().Clone(ctx))

			// Listen result for next(c)
			done <- next(c)
		}()

		// Multiplexer to get which comes first between processes or timeout
		select {
		case <-ctx.Done(): // Request Timeout
			return errorHandleFunc(statusCode, response)

		case err := <-done: // Normal Process
			return err
		}
	}
}
