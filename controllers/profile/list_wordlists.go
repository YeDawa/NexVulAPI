package profile

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"nexvul/configs"
	"nexvul/controllers/users"
	"nexvul/models"
	"nexvul/utils"

	"github.com/labstack/echo/v4"
)

type Wordlists struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Public     bool   `json:"public"`
	TotalLines uint   `json:"total_lines"`
	CreatedAt  string `json:"created_at"`
	HTMLPage   string `json:"html_page"`
	APIPage    string `json:"api_page"`
	RawPage    string `json:"raw_page"`
	Username   string `json:"username"`
}

func ListPublicWordlists(c echo.Context) error {
	user := c.Param("user")
	username := users.GetUserIDByUsername(user)

	if username == "" {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "User not found",
		})
	}

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 50
	}

	offset := (page - 1) * limit
	search := strings.TrimSpace(c.QueryParam("search"))
	order := c.QueryParam("order")

	userID, err := strconv.Atoi(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	query := configs.DB.Where("user_id = ? AND public = 'true'", userID)
	if search != "" {
		query = query.Where("slug LIKE ? OR urls LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var orderBy string
	switch order {
	case "created_at asc", "created_at desc":
		orderBy = "created_at " + strings.ToUpper(strings.Split(order, " ")[1])
	default:
		orderBy = "created_at DESC"
	}

	query = query.Order(orderBy)

	var total int64
	query.Model(&models.CustomWordlists{}).Count(&total)

	var wordlistResponse []models.CustomWordlists
	result := query.Limit(limit).Offset(offset).Find(&wordlistResponse)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	var wordlistData []Wordlists
	for _, wordlist := range wordlistResponse {
		wordlistData = append(wordlistData, Wordlists{
			Slug:       wordlist.Slug,
			Username:   user,
			Name:       wordlist.Name,
			TotalLines: uint(wordlist.TotalLines),
			Public:     wordlist.Public,
			CreatedAt:  wordlist.CreatedAt.Format(time.RFC3339),
			HTMLPage:   utils.GetWordlistPage(wordlist.Slug),
			APIPage:    utils.GetWordlistApiPage(c, wordlist.Slug),
			RawPage:    utils.GetWordlistRawPage(c, wordlist.Slug),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    wordlistData,
		"page":    page,
		"limit":   limit,
		"total":   total,
	})
}
