package http

import (
	"github.com/labstack/echo/v4"
	"nix_education/internal/infra/http/controllers"
)

func AuthGroup(g *echo.Group, authController controllers.AuthController) {
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
	g.POST("/loginGoogle", authController.LoginGoogle)
}
