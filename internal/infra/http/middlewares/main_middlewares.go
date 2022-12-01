package middlewares

import (
	"github.com/labstack/echo/v4"
)

func SetMainMiddlewares(e *echo.Echo) {
	e.Use(serverHeader)
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "NIX/v1.0")
		return next(c)
	}
}
