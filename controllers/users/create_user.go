package users

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"httpshield/configs"
	"httpshield/core"
	"httpshield/generator"
	"httpshield/models"
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

	apiKey := "hs_" + generator.String(32, 36)
	salt, _ := security.GenerateRandomSalt(16)

	newUser := models.Users{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		ApiKey:   apiKey,
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

	userId := GetUserIDByUsername(req.Username)
	userIdUint, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to convert userId to uint",
		})
	}
	newProfile := models.Profile{
		UserId:     uint(userIdUint),
		PublicName: req.Name,
		Contact:    req.Email,
	}

	configs.DB.Create(&newProfile)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User created successfully",
	})
}
