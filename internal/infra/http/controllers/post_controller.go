package controllers

import (
	"NIX/internal/app"
	"NIX/internal/domain"
	"NIX/internal/infra/http/requests"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type PostController struct {
	postService app.PostService
}

func NewPostController(s app.PostService) PostController {
	return PostController{
		postService: s,
	}
}

func (c PostController) Save(ctx echo.Context) error {
	var post requests.PostRequest
	err := ctx.Bind(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	p, err := post.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	createdPost, err := c.postService.Save(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdPost)
}

func (c PostController) Find(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	post, err := c.postService.Find(postId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, post)
}

func (c PostController) FindAll(ctx echo.Context) error {
	userId := int64(7)
	pagination, err := requests.DecodePaginationQuery(ctx.Request())

	posts, err := c.postService.FindAll(userId, pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, posts)
}

func (c PostController) Update(ctx echo.Context) error {
	var post requests.PostRequest
	err := ctx.Bind(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&post)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	p, err := post.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	createdPost, err := c.postService.Update(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdPost)
}

func (c PostController) Delete(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = c.postService.Delete(postId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, domain.OK)
}
