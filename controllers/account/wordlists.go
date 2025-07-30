package account

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"nexvul/configs"
	"nexvul/models"
	"nexvul/utils"

	"nexvul/controllers/users"

	"github.com/labstack/echo/v4"
)

type Wordlists struct {
	Name       string `json:"name"`
	TotalLines uint   `json:"total_lines"`
	Slug       string `json:"slug"`
	CreatedAt  string `json:"created_at"`
	HTMLPage   string `json:"html_page"`
	APIPage    string `json:"api_page"`
}

func ListUserWordlists(c echo.Context) error {
	userID := users.GetUserID(c)
	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized user",
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

	query := configs.DB.Where("user_id = ?", uint(userIDUint))

	if search != "" {
		query = query.Where("name LIKE ? OR url LIKE ?", "%"+search+"%", "%"+search+"%")
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
	var publicWordlists int64
	var privateWordlists int64

	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ?", userID).Count(&total)
	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ? AND public = 'true'", userID).Count(&publicWordlists)
	configs.DB.Model(&models.CustomWordlists{}).Where("user_id = ? AND public = 'false'", userID).Count(&privateWordlists)

	var wordlistsResponse []models.CustomWordlists
	result := query.Limit(limit).Offset(offset).Find(&wordlistsResponse)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	var wordlistData []Wordlists
	for _, wordlist := range wordlistsResponse {
		wordlistData = append(wordlistData, Wordlists{
			Name:       wordlist.Name,
			Slug:       wordlist.Slug,
			TotalLines: uint(wordlist.TotalLines),
			CreatedAt:  wordlist.CreatedAt.Format(time.RFC3339),
			HTMLPage:   utils.GetWordlistPage(wordlist.Slug),
			APIPage:    utils.GetWordlistApiPage(c, wordlist.Slug),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":           true,
		"data":              wordlistData,
		"page":              page,
		"limit":             limit,
		"total":             total,
		"public_wordlists":  publicWordlists,
		"private_wordlists": privateWordlists,
	})
}
