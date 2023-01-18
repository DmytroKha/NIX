package controllers

import (
	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	UserKey      = "user"
	UserIdKey    = "jti"
	UserEmailKey = "email"
)

func GetUserValueFromJWT(ctx echo.Context, key string) string {
	user := ctx.Get(UserKey)
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	val := claims[key].(string)
	return val
}

func FormatedResponse(ctx echo.Context, code int, i interface{}) error {
	ct := ctx.Request().Header.Get("Accept")
	if ct == "text/xml" {
		return ctx.XML(code, i)
	} else {
		return ctx.JSON(code, i)
	}
}
