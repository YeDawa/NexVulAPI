package account

import (
	"encoding/json"
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

type Scans struct {
	Urls       []string `json:"urls"`
	Slug       string   `json:"slug"`
	Public     bool     `json:"public"`
	CreatedAt  string   `json:"created_at"`
	HTMLPage   string   `json:"html_page"`
	APIPage    string   `json:"api_page"`
	ReportPage string   `json:"report_page"`
	Username   string   `json:"username"`
}

func ListUserScans(c echo.Context) error {
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

	query := configs.DB.Where("user_id = ?", userIDUint)
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
	var publicScans int64
	var privateScans int64

	query.Model(&models.Scans{}).Count(&total)
	configs.DB.Model(&models.Scans{}).Where("user_id = ? AND public = 'true'", userID).Count(&publicScans)
	configs.DB.Model(&models.Scans{}).Where("user_id = ? AND public = 'false'", userID).Count(&privateScans)

	var scansResponse []models.Scans
	result := query.Limit(limit).Offset(offset).Find(&scansResponse)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	var scanData []Scans
	for _, scan := range scansResponse {
		user := users.GetUsernameByID(scan.UserId)

		var urls []string

		if err := json.Unmarshal([]byte(scan.Urls), &urls); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   "Failed to deserialize 'urls' field: " + err.Error(),
			})
		}

		scanData = append(scanData, Scans{
			Slug:       scan.Slug,
			Urls:       urls,
			Username:   user,
			Public:     scan.Public,
			CreatedAt:  scan.CreatedAt.Format(time.RFC3339),
			HTMLPage:   utils.GetScanPage(scan.Slug),
			APIPage:    utils.GetScanApiPage(c, scan.Slug),
			ReportPage: utils.GetScanApiReportPage(c, scan.Slug),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":       true,
		"data":          scanData,
		"page":          page,
		"limit":         limit,
		"total":         total,
		"public_scans":  publicScans,
		"private_scans": privateScans,
	})
}
