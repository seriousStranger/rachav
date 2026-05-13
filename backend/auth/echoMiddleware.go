package auth

import (
	"log/slog"

	"github.com/kopkapozla/rachav/config"
	"github.com/labstack/echo/v5"
)

func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		user, pass := config.Config.GetAuthPair()
		isAuth := СheckAuth(c.Request().Header.Get("Api-Authorization"),
			"Basic ",
			map[string]string{user: pass},
		)
		if !isAuth {
			slog.Warn("api failed login attempt")
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
