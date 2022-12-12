package http

import (
	"github.com/labstack/echo/v4"
	"nix_education/internal/infra/http/controllers"
)

func UsersGroup(g *echo.Group, userController controllers.UserController) {
	g.PUT("/:id", userController.SetPassword)
}
