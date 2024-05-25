package server

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func routerDocs(e *echo.Echo) {
	e.GET("/docs/*", echoSwagger.WrapHandler)
}
