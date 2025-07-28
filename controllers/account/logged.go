package account

import (
	"encoding/base64"
	"net/http"
	"strings"

	"nexvul/configs"
	"nexvul/core"
	"nexvul/models"

	"github.com/drexedam/gravatar"
	"github.com/labstack/echo/v4"
)

func UserLogged(c echo.Context) error {
	cookieValue := core.GetCookie(c, configs.UserCookieName)

	decodedValue, _ := base64.StdEncoding.DecodeString(cookieValue)
	parts := strings.Split(string(decodedValue), ":")

	if len(parts) < 2 {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized user",
		})
	}

	userID := parts[1]
	var user models.Users
	result := configs.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	var profile models.Profile
	resultProfile := configs.DB.Where("id = ?", userID).First(&profile)
	if resultProfile.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	var totalScans int64
	var publicScans int64
	var privateScans int64

	configs.DB.Model(&models.Scans{}).Where("user_id = ?", userID).Count(&totalScans)
	configs.DB.Model(&models.Scans{}).Where("user_id = ? AND public = 'true'", userID).Count(&publicScans)
	configs.DB.Model(&models.Scans{}).Where("user_id = ? AND public = 'false'", userID).Count(&privateScans)

	var totalWordlists int64
	var publicWordlists int64
	var privateWordlists int64

	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ?", userID).Count(&totalWordlists)
	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ? AND public = 'true'", userID).Count(&publicWordlists)
	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ? AND public = 'false'", userID).Count(&privateWordlists)

	statsScans := map[string]int64{
		"total":   totalScans,
		"public":  publicScans,
		"private": privateScans,
	}

	statsWordlists := map[string]int64{
		"total":   totalWordlists,
		"public":  publicWordlists,
		"private": privateWordlists,
	}

	profileData := map[string]string{
		"public_name": profile.PublicName,
		"contact":     profile.Contact,
		"linkedin":    profile.Linkedin,
		"website":     profile.Website,
		"bio":         profile.Bio,
		"location":    profile.Location,
		"github":      profile.Github,
		"twitter":     profile.Twitter,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"name":       user.Name,
		"username":   user.Username,
		"plan":       user.Plan,
		"email":      user.Email,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"api_key":    user.ApiKey,
		"avatar":     gravatar.New(user.Email).Size(300).AvatarURL(),
		"profile":    profileData,
		"scans":      statsScans,
		"wordlists":  statsWordlists,
	})
}
