package controllers

import (
	"NIX/internal/app"
	"NIX/internal/domain"
	"NIX/internal/infra/http/requests"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CommentController struct {
	commentService app.CommentService
}

func NewCommentController(s app.CommentService) CommentController {
	return CommentController{
		commentService: s,
	}
}

func (c CommentController) Save(ctx echo.Context) error {
	var comment requests.CommentRequest
	err := ctx.Bind(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	p, err := comment.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	createdComment, err := c.commentService.Save(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdComment)
}

func (c CommentController) Find(ctx echo.Context) error {
	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	comment, err := c.commentService.Find(commentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, comment)
}

func (c CommentController) FindAll(ctx echo.Context) error {
	userId := int64(7)
	pagination, err := requests.DecodePaginationQuery(ctx.Request())

	comments, err := c.commentService.FindAll(userId, pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, comments)
}

func (c CommentController) Update(ctx echo.Context) error {
	var comment requests.CommentRequest
	err := ctx.Bind(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	p, err := comment.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	createdComment, err := c.commentService.Update(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdComment)
}

func (c CommentController) Delete(ctx echo.Context) error {
	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = c.commentService.Delete(commentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, domain.OK)
}
