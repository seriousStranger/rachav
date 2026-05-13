package panel

import (
	"net/http"

	"github.com/kopkapozla/rachav/config"
	"github.com/labstack/echo/v5"
)

func GetHost(context *echo.Context) error {
	return context.JSON(http.StatusOK, map[string]string{"host": config.Config.GetHost()})
}
