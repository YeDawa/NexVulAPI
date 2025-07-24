package users

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"nexvul/configs"
	"nexvul/core"
	"nexvul/models"

	"github.com/drexedam/gravatar"
	"github.com/labstack/echo/v4"
)

func GetUserID(c echo.Context) string {
	cookieValue := core.GetCookie(c, configs.UserCookieName)
	if cookieValue == "" {
		return ""
	}

	decodedValue, err := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		return ""
	}

	parts := strings.Split(string(decodedValue), ":")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func GetAvatarByID(id uint) string {
	var email string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("email").
		Where("id = ?", id).
		Scan(&email).Error; err != nil {
		return ""
	}

	return gravatar.New(email).Size(300).AvatarURL()
}

func GetUsernameByID(id uint) string {
	var username string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("username").
		Where("id = ?", id).
		Scan(&username).Error; err != nil {
		return ""
	}

	return username
}

func GetHistory(id uint) bool {
	var history string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("history").
		Where("id = ?", id).
		Scan(&history).Error; err != nil {
		return false
	}

	return history == "true"
}

func GetUserIDByUsername(username string) string {
	var id uint

	if err := configs.DB.
		Model(&models.Users{}).
		Select("id").
		Where("username = ?", username).
		Scan(&id).Error; err != nil {
		return ""
	}

	return strconv.FormatUint(uint64(id), 10)
}

func GetNameByID(id uint) string {
	var name string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("name").
		Where("id = ?", id).
		Scan(&name).Error; err != nil {
		return ""
	}

	return name
}

func GetEmailByID(id uint) string {
	var email string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("email").
		Where("id = ?", id).
		Scan(&email).Error; err != nil {
		return ""
	}

	return email
}

func GetPlanByID(id uint) string {
	var plan string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("plan").
		Where("id = ?", id).
		Scan(&plan).Error; err != nil {
		return ""
	}

	return plan
}

func GetStatusByID(id uint) string {
	var status string

	if err := configs.DB.
		Model(&models.Users{}).
		Select("status").
		Where("id = ?", id).
		Scan(&status).Error; err != nil {
		return ""
	}

	return status
}

func GetCreatedAtByID(id uint) string {
	var createdAt time.Time

	if err := configs.DB.
		Model(&models.Users{}).
		Select("created_at").
		Where("id = ?", id).
		Scan(&createdAt).Error; err != nil {
		return ""
	}

	return createdAt.Format(time.RFC3339)
}
