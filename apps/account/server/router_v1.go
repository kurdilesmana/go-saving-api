package server

import (
	middle "github.com/kurdilesmana/go-saving-api/apps/account/server/middlewares"
	"github.com/labstack/echo/v4"
)

func routerGroupV1(handler Handler, e *echo.Echo, middle *middle.AuthMiddleware) {
	e.GET("/favicon.ico", func(c echo.Context) error { return nil })

	api := e.Group("/api") // /api

	v1 := api.Group("/v1", func(next echo.HandlerFunc) echo.HandlerFunc { // middleware for /api/v1
		return func(c echo.Context) error {
			c.Set("Version", "v1")
			return next(c)
		}
	})

	v1.POST("/daftar", handler.accountHandler.CreateHandler)
	v1.GET("/inquiry", handler.accountHandler.InquiryHandler, middle.AuthorizeToken())
	v1.POST("/tabung", handler.transactionHandler.SavingHandler, middle.AuthorizeToken())
	v1.POST("/tarik", handler.transactionHandler.CashWithdrawlHandler, middle.AuthorizeToken())
}
