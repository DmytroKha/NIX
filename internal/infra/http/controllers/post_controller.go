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

	jtl := GetUserValueFromJWT(ctx, UserIdKey)
	userId, err := strconv.Atoi(jtl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p.UserId = int64(userId)

	createdPost, err := c.postService.Save(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var postDto resources.PostDto

	return ctx.JSON(http.StatusCreated, postDto.DomainToDto(createdPost))
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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	post, err := c.postService.Find(postId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	var postDto resources.PostDto

	return ctx.JSON(http.StatusOK, postDto.DomainToDto(post))
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
	rqst := ctx.Request()
	pagination, err := requests.DecodePaginationQuery(rqst)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	posts, err := c.postService.FindAll(pagination)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	var postDto resources.PostDto
	//hdr := ctx.Response().Header()
	//ct := hdr.Get("Content-Type")
	ct := rqst.Header.Get("Accept")
	if ct == "text/xml" {
		return ctx.XML(http.StatusOK, postDto.DomainToDtoCollection(posts))
	} else {
		return ctx.JSON(http.StatusOK, postDto.DomainToDtoCollection(posts))
	}
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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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

	jtl := GetUserValueFromJWT(ctx, UserIdKey)
	userId, err := strconv.Atoi(jtl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p.UserId = int64(userId)
	p.Id = postId

	updatedPost, err := c.postService.Update(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var postDto resources.PostDto

	return ctx.JSON(http.StatusOK, postDto.DomainToDto(updatedPost))
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
// @Success      200  {object}  domain.Post
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /posts/{id} [delete]
func (c PostController) Delete(ctx echo.Context) error {
	postId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	jtl := GetUserValueFromJWT(ctx, UserIdKey)
	userId, err := strconv.Atoi(jtl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = c.postService.Delete(postId, int64(userId))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, domain.OK)
}
