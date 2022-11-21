package router

import (
	"NIX/config"
	"NIX/internal/infra/http"
	"NIX/internal/infra/http/controllers"
	"NIX/internal/infra/http/middlewares"
	"github.com/labstack/echo/v4"
)

func New(
	userController controllers.UserController,
	authController controllers.AuthController,
	postController controllers.PostController,
	commentController controllers.CommentController,
	cf config.Configuration) *echo.Echo {

	e := echo.New()

	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	users := api.Group("/users")
	posts := api.Group("/posts")
	comments := posts.Group("/:postId/comments")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(api)
	middlewares.SetJWTMiddlewares(posts, cf)

	http.MainGroup(e, authController)
	http.AuthGroup(auth, authController)
	http.UsersGroup(users, userController)
	http.PostsGroup(posts, postController)
	http.CommentsGroup(comments, commentController)

	return e
}
