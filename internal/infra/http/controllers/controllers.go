package controllers

import (
	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

/* should not use built-in type string as key for value;
define your own type to avoid collisions */
// type CtxStrKey string

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
