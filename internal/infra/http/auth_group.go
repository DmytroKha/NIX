package http

import (
	"NIX/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
)

func AuthGroup(g *echo.Group, authController controllers.AuthController) {
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
	g.POST("/loginGoogle", authController.LoginGoogle)
}
