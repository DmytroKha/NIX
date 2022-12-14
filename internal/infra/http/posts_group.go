package http

import (
	"github.com/labstack/echo/v4"
	"nix_education/internal/infra/http/controllers"
)

func PostsGroup(g *echo.Group, postController controllers.PostController) {
	g.GET("", postController.FindAll)
	g.POST("", postController.Save)
	g.GET("/:id", postController.Find)
	g.PUT("/:id", postController.Update) //розібратися з контекстом
	g.DELETE("/:id", postController.Delete)
}
