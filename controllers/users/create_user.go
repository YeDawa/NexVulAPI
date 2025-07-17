package users

import (
	"net/http"
	"encoding/base64"
	
	"httpshield/core"
	"httpshield/models"
	"httpshield/configs"
	"httpshield/security"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(c echo.Context) error {
	req := new(CreateUserRequest)

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

	salt, _ := security.GenerateRandomSalt(16)
	newUser := models.Users{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: security.HashPassword(req.Password, salt),
		Salt:     base64.StdEncoding.EncodeToString(salt),
	}

	result := configs.DB.Create(&newUser)
	if result.Error != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User created successfully",
	})
}
