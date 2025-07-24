package users

import (
	"encoding/base64"
	"net/http"

	"nexvul/configs"
	"nexvul/models"
	"nexvul/security"

	"nexvul/core"
	"nexvul/generator"

	"github.com/labstack/echo/v4"
)

type User struct {
	CookieName string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c echo.Context) error {
	req := new(LoginRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !core.ValidateEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid email format",
		})
	}

	var user models.Users
	result := configs.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Invalid email or password",
		})
	}

	salt, _ := base64.StdEncoding.DecodeString(user.Salt)

	if !security.VerifyPassword(req.Password, user.Password, salt) {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Invalid email or password",
		})
	}

	cookieValue, err := generator.CookieValue(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to generate cookie",
		})
	}

	core.SetCookie(c, configs.UserCookieName, cookieValue, 48)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"id":       user.Id,
		"name":     user.Name,
		"plan":     user.Plan,
		"status":   user.Status,
		"username": user.Username,
	})
}

func HasLogged(c echo.Context) error {
	cookieValue := core.GetCookie(c, configs.UserCookieName)
	isLoggedIn := cookieValue != ""

	return c.JSON(http.StatusOK, map[string]interface{}{
		"logged_in": isLoggedIn,
	})
}

func Logoff(c echo.Context) error {
	core.SetCookie(c, configs.UserCookieName, "", -1)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged off successfully",
	})
}
