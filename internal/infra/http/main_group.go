package http

import (
	"NIX/internal/infra/http/controllers"
	"NIX/internal/infra/http/requests"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func MainGroup(e *echo.Echo, authController controllers.AuthController) {
	e.Validator = requests.NewValidator()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/callback", authController.Callback)
}
