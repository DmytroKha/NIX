package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"nix_education/config"
	"nix_education/internal/infra/http/controllers"
)

func SetJWTMiddlewares(g *echo.Group, cf config.Configuration) {
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(cf.JwtSecret),
		ContextKey:    controllers.UserKey,
	}))
}
