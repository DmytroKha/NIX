package http

import (
	"github.com/labstack/echo/v4"
	"nix_education/internal/infra/http/controllers"
)

func CommentsGroup(g *echo.Group, commentController controllers.CommentController) {
	g.GET("", commentController.FindAll)
	g.POST("", commentController.Save)
	g.GET("/:id", commentController.Find)
	g.PUT("/:id", commentController.Update)
	g.DELETE("/:id", commentController.Delete)
}
