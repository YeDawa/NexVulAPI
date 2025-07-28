package profile

import (
	"net/http"
	"strconv"
	"time"

	"nexvul/configs"
	"nexvul/models"

	"nexvul/controllers/users"

	"github.com/labstack/echo/v4"
)

type ProfileDetails struct {
	Avatar     string `json:"avatar"`
	PublicName string `json:"public_name"`
	Website    string `json:"website"`
	Bio        string `json:"bio"`
	Location   string `json:"location"`
	Contact    string `json:"contact"`
	CreatedAt  string `json:"created_at"`
	Twitter    string `json:"twitter"`
	Github     string `json:"github"`
	Status     string `json:"status"`
	Plan       string `json:"plan"`
	Linkedin   string `json:"linkedin"`
}

type ProfileStats struct {
	Scans     int64 `json:"scans"`
	Wordlists int64 `json:"wordlists"`
}

func ProfilePublic(c echo.Context) error {
	username := c.Param("user")
	userId := users.GetUserIDByUsername(username)

	userIDUint, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	var UserInfo models.Users
	var profile models.Profile

	resultProfile := configs.DB.Where("id = ?", userId).First(&profile)
	resultUser := configs.DB.Where("id = ?", userIDUint).First(&UserInfo)

	if resultProfile.Error != nil || resultUser.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	var totalScans int64
	var totalWordlists int64
	
	configs.DB.Model(&models.Scans{}).Where("user_id = ? AND public = 'true'", uint(userIDUint)).Count(&totalScans)
	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ? AND public = 'true'", uint(userIDUint)).Count(&totalWordlists)

	ProfileDetails := ProfileDetails{
		Avatar:     users.GetAvatarByID(UserInfo.Id),
		PublicName: profile.PublicName,
		Website:    profile.Website,
		Bio:        profile.Bio,
		Location:   profile.Location,
		Contact:    profile.Contact,
		Plan:       string(UserInfo.Plan),
		Status:     string(UserInfo.Status),
		CreatedAt:  UserInfo.CreatedAt.Local().Format(time.RFC3339),
		Twitter:    profile.Twitter,
		Github:     profile.Github,
		Linkedin:   profile.Linkedin,
	}

	ProfileStats := ProfileStats{
		Scans:     totalScans,
		Wordlists: totalWordlists,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"stats":   ProfileStats,
		"details": ProfileDetails,
	})
}
