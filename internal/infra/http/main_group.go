package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"nix_education/internal/infra/http/controllers"
	"nix_education/internal/infra/http/requests"
)

func MainGroup(e *echo.Echo, authController controllers.AuthController) {
	e.Validator = requests.NewValidator()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/callback", authController.Callback)
}
