package wordlists

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"nexvul/configs"
	"nexvul/controllers/users"
	"nexvul/generator"
	"nexvul/models"
	"nexvul/utils"

	"github.com/labstack/echo/v4"
)

type CreateWordlistRequest struct {
	URL string `json:"url"`
}

type ImportWordlistResponse struct {
	Slug     string `json:"slug"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	HTMLPage string `json:"html_page"`
}

func ImportWordlist(c echo.Context) error {
	slug := generator.String(8, 16)
	req := new(CreateWordlistRequest)
	userID := users.GetUserID(c)
	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	lastPart := utils.GetLastPartOfURL(req.URL)
	totalLines, err := utils.CountRemoteFileLines(req.URL)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("failed to count lines in wordlist: %v", err),
		})
	}

	customWordlist := models.CustomWordlists{
		Slug:       slug,
		Name:       lastPart,
		Url:        req.URL,
		FileName:   lastPart,
		CreatedAt:  time.Now(),
		TotalLines: totalLines,
		UserId:     uint(userIDUint),
	}

	if err := configs.DB.Create(&customWordlist).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("failed to create wordlist: %v", err),
		})
	}

	response := ImportWordlistResponse{
		Slug:     slug,
		Success:  true,
		HTMLPage: utils.GetWordlistPage(slug),
		Message:  "Wordlist created successfully",
	}

	return c.JSON(http.StatusOK, response)
}
