package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"nix_education/internal/app"
	"nix_education/internal/domain"
	"nix_education/internal/infra/http/requests"
	"nix_education/internal/infra/http/resources"
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
// @Produce      xml
// @Param        postId   path      string  true  "Post ID"
// @Param        input   body      requests.CommentRequest  true  "Comment body"
// @Success      201  {object}  resources.CommentDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments [post]
func (c CommentController) Save(ctx echo.Context) error {
	var comment requests.CommentRequest
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Bind(&comment)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&comment)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	email := GetUserValueFromJWT(ctx, UserEmailKey)
	createdComment, err := c.commentService.Save(comment, postId, email)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var commentDto resources.CommentDto
	return FormatedResponse(ctx, http.StatusCreated, commentDto.DatabaseToDto(createdComment))
}

// FindComment godoc
// @Summary      Show a comment
// @Security     ApiKeyAuth
// @Description  get comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        postId   path      string  true  "Post ID"
// @Param        id   path      string  true  "Comment ID"
// @Success      200  {object}  resources.CommentDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments/{id} [get]
func (c CommentController) Find(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	comment, err := c.commentService.Find(postId, commentId)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	var commentDto resources.CommentDto
	return FormatedResponse(ctx, http.StatusOK, commentDto.DatabaseToDto(comment))
}

// ListComments godoc
// @Summary      Show all comments
// @Security     ApiKeyAuth
// @Description  get all comments
// @Tags         comments
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        postId   path      string  true  "Post ID"
// @Success      200  {object}  resources.CommentDto
// @Failure      400  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments [get]
func (c CommentController) FindAll(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	pagination, err := requests.DecodePaginationQuery(ctx.Request())
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	comments, err := c.commentService.FindAll(postId, pagination)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	var commentDto resources.CommentDto
	return FormatedResponse(ctx, http.StatusOK, commentDto.DatabaseToDtoCollection(comments))
}

// UpdateComment godoc
// @Summary      Update comment
// @Security     ApiKeyAuth
// @Description  update comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        postId   path      string  true  "Post ID"
// @Param        id   path      string  true  "Comment ID"
// @Param        input   body      requests.CommentRequest  true  "Comment body"
// @Success      200  {object}  resources.CommentDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments/{id} [put]
func (c CommentController) Update(ctx echo.Context) error {
	var comment requests.CommentRequest
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Bind(&comment)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&comment)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	email := GetUserValueFromJWT(ctx, UserEmailKey)

	updatedComment, err := c.commentService.Update(comment, postId, id, email)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var commentDto resources.CommentDto
	return FormatedResponse(ctx, http.StatusOK, commentDto.DatabaseToDto(updatedComment))
}

// DeleteComment godoc
// @Summary      Delete comment
// @Security     ApiKeyAuth
// @Description  delete comment by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        postId   path      string  true  "Post ID"
// @Param        id   path      string  true  "Comment ID"
// @Success      200  {object}  resources.CommentDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{postId}/comments/{id} [delete]
func (c CommentController) Delete(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	commentId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	email := GetUserValueFromJWT(ctx, UserEmailKey)
	err = c.commentService.Delete(postId, commentId, email)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	return FormatedResponse(ctx, http.StatusOK, domain.OK)
}
