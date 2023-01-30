package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"nix_education/internal/app"
	"nix_education/internal/infra/http/requests"
	"nix_education/internal/infra/http/resources"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(us app.UserService) UserController {
	return UserController{
		userService: us,
	}
}

// SetPass godoc
// @Summary      Set password for user google acc
// @Security     ApiKeyAuth
// @Description  set password for user google acc
// @Tags         users
// @Accept       json
// @Produce      json
// @Produce      xml
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
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&user)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	updatedUser, err := c.userService.SetPassword(user)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var userDto resources.UserDto
	return FormatedResponse(ctx, http.StatusOK, userDto.DatabaseToDto(updatedUser))
}
