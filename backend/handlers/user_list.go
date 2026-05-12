package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kopkapozla/rachav/database"
	"github.com/labstack/echo/v5"
)

func GetUserList(context *echo.Context) error {
	userList, err := database.Load()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, userList)
}

func PostUserList(context *echo.Context) error {
	var userList map[string]string

	err := json.NewDecoder(context.Request().Body).Decode(&userList)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"body must be map[string]string",
		)
	}

	err = database.Save(userList)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return context.NoContent(http.StatusOK)
}
