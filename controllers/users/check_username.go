package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckUsername(c echo.Context) error {
	username := c.Param("user")
	userID := GetUserIDByUsername(username)

	if userID == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"avaliable": true,
			"message":   "Username is available",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"avaliable": false,
		"message":   "Username is not available",
	})
}
