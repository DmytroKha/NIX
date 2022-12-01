package controllers

import (
	"NIX/internal/app"
	"NIX/internal/infra/http/requests"
	"NIX/internal/infra/http/resources"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(us app.UserService) UserController {
	return UserController{
		userService: us,
	}
}

func (c UserController) Save(ctx echo.Context) error {
	var user requests.UserRequest
	err := ctx.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	u, err := user.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	createdUser, err := c.userService.Save(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, createdUser)
}

// SetPass godoc
// @Summary      Set password for user google acc
// @Security     ApiKeyAuth
// @Description  set password for user google acc
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        input   body      requests.UserRequest  true  "User body"
// @Success      200  {object}  resources.UserDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /users/{id} [put]
func (c UserController) SetPassword(ctx echo.Context) error {
	var user requests.UserRequest
	err := ctx.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = ctx.Validate(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	u, err := user.ToDomainModel()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updatedUser, err := c.userService.SetPassword(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var userDto resources.UserDto

	return ctx.JSON(http.StatusOK, userDto.DomainToDto(updatedUser))
}