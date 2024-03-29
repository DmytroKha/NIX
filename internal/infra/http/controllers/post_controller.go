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

type PostController struct {
	postService app.PostService
}

func NewPostController(s app.PostService) PostController {
	return PostController{
		postService: s,
	}
}

// NewPost godoc
// @Summary      Create a new post
// @Security     ApiKeyAuth
// @Description  save a post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        input   body      requests.PostRequest  true  "Post body"
// @Success      201  {object}  resources.PostDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /posts [post]
func (c PostController) Save(ctx echo.Context) error {
	var post requests.PostRequest
	err := ctx.Bind(&post)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&post)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	jtl := GetUserValueFromJWT(ctx, UserIdKey)
	userId, err := strconv.Atoi(jtl)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	createdPost, err := c.postService.Save(post, int64(userId))
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var postDto resources.PostDto
	return FormatedResponse(ctx, http.StatusCreated, postDto.DatabaseToDto(createdPost))
}

// FindPost godoc
// @Summary      Show a post
// @Security     ApiKeyAuth
// @Description  get post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  resources.PostDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{id} [get]
func (c PostController) Find(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	post, err := c.postService.Find(postId)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	var postDto resources.PostDto
	return FormatedResponse(ctx, http.StatusOK, postDto.DatabaseToDto(post))
}

// ListPosts godoc
// @Summary      Show all posts
// @Security     ApiKeyAuth
// @Description  get all posts
// @Tags         posts
// @Accept       json
// @Produce      json
// @Produce      xml
// @Success      200  {object}  resources.PostDto
// @Failure      400  {string}  echo.HTTPError
// @Router       /posts [get]
func (c PostController) FindAll(ctx echo.Context) error {
	pagination, err := requests.DecodePaginationQuery(ctx.Request())
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	posts, err := c.postService.FindAll(pagination)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	var postDto resources.PostDto
	return FormatedResponse(ctx, http.StatusOK, postDto.DatabaseToDtoCollection(posts))
}

// UpdatePost godoc
// @Summary      Update post
// @Security     ApiKeyAuth
// @Description  update post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        id   path      string  true  "Post ID"
// @Param        input   body      requests.PostRequest  true  "Post body"
// @Success      200  {object}  resources.PostDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /posts/{id} [put]
func (c PostController) Update(ctx echo.Context) error {
	var post requests.PostRequest
	err := ctx.Bind(&post)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&post)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	jtl := GetUserValueFromJWT(ctx, UserIdKey)
	userId, err := strconv.Atoi(jtl)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	updatedPost, err := c.postService.Update(post, postId, int64(userId))
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var postDto resources.PostDto
	return FormatedResponse(ctx, http.StatusOK, postDto.DatabaseToDto(updatedPost))
}

// DeletePost godoc
// @Summary      Delete post
// @Security     ApiKeyAuth
// @Description  delete post by ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        id   path      string  true  "Post ID"
// @Success      200  {object}  resources.PostDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{id} [delete]
func (c PostController) Delete(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	jtl := GetUserValueFromJWT(ctx, UserIdKey)
	userId, err := strconv.Atoi(jtl)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = c.postService.Delete(postId, int64(userId))
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	return FormatedResponse(ctx, http.StatusOK, domain.OK)
}
