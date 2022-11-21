package middlewares

import (
	"NIX/config"
	"NIX/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetJWTMiddlewares(g *echo.Group, cf config.Configuration) {
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(cf.JwtSecret),
		ContextKey:    controllers.UserKey,
	}))
}
