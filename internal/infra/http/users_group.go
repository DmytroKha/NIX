package http

import (
	"NIX/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
)

func UsersGroup(g *echo.Group, userController controllers.UserController) {
	g.PUT("/:id", userController.SetPassword)
}
