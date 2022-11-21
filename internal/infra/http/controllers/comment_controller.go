package controllers

import (
	"NIX/internal/app"
	"NIX/internal/domain"
	"NIX/internal/infra/http/requests"
	"errors"
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

// NewComment godoc
// @Summary      Create a new comment
// @Security     ApiKeyAuth
// @Description  save a comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Param        input   body      requests.CommentRequest  true  "Comment body"
// @Success      201  {object}  domain.Comment
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments [post]
func (c CommentController) Save(ctx echo.Context) error {
	var comment requests.CommentRequest
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Bind(&comment)
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

	p.PostId = postId

	createdComment, err := c.commentService.Save(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdComment)
}

// FindComment godoc
// @Summary      Show a comment
// @Security     ApiKeyAuth
// @Description  get comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Param        id   path      string  true  "Comment ID"
// @Success      200  {object}  domain.Comment
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments/{id} [get]
func (c CommentController) Find(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	comment, err := c.commentService.Find(postId, commentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, comment)
}

// ListComments godoc
// @Summary      Show all comments
// @Security     ApiKeyAuth
// @Description  get all comments
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Success      200  {object}  domain.Comment
// @Failure      400  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments [get]
func (c CommentController) FindAll(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	pagination, err := requests.DecodePaginationQuery(ctx.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	comments, err := c.commentService.FindAll(postId, pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary      Update comment
// @Security     ApiKeyAuth
// @Description  update comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Param        id   path      string  true  "Comment ID"
// @Param        input   body      requests.CommentRequest  true  "Comment body"
// @Success      201  {object}  domain.Comment
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments/{id} [put]
func (c CommentController) Update(ctx echo.Context) error {
	var comment requests.CommentRequest

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Bind(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	p, err := comment.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	email := GetUserValueFromJWT(ctx, UserEmailKey)
	if email != p.Email {
		err = errors.New("user email mismatch")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	p.PostId = postId
	p.Id = id

	createdComment, err := c.commentService.Update(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdComment)
}

// DeleteComment godoc
// @Summary      Delete comment
// @Security     ApiKeyAuth
// @Description  delete comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Param        id   path      string  true  "Comment ID"
// @Success      200  {object}  domain.Comment
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments/{id} [delete]
func (c CommentController) Delete(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	email := GetUserValueFromJWT(ctx, UserEmailKey)

	err = c.commentService.Delete(postId, commentId, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, domain.OK)
}
