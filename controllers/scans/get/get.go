package get_scan

import (
	"encoding/json"
	"net/http"
	"time"

	"httpshield/configs"
	"httpshield/controllers/users"
	"httpshield/models"
	"httpshield/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HeaderResult struct {
	Header string `json:"header"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type ScanData struct {
	URL             string         `json:"url"`
	Server          string         `json:"server"`
	Method          string         `json:"method"`
	ExecutionTime   int64          `json:"execution_time"`
	StatusCode      int            `json:"status_code"`
	ContentType     string         `json:"content_type"`
	Results         []HeaderResult `json:"results"`
	SecurityScore   int            `json:"security_score"`
	Recommendations []string       `json:"recommendations"`
}

type ScanResponse struct {
	Id         string     `json:"id"`
	Data       []ScanData `json:"data"`
	Urls       []string   `json:"urls"`
	HtmlPage   string     `json:"html_page"`
	ReportPage string     `json:"report_page"`
	ApiPage    string     `json:"api_page"`
	Public     bool       `json:"public"`
	Owner      ScanOwner  `json:"owner"`
	CreatedAt  string     `json:"created_at"`
}

type ScanOwner struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
}

func GetScanDetails(c echo.Context) error {
	id := c.Param("id")

	var scans models.Scans
	result := configs.DB.Where("slug = ?", id).First(&scans)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Item not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	var scanData []ScanData
	var urls []string

	if err := json.Unmarshal([]byte(scans.Data), &scanData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to deserialize 'data' field: " + err.Error(),
		})
	}

	if err := json.Unmarshal([]byte(scans.Urls), &urls); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to deserialize 'urls' field: " + err.Error(),
		})
	}

	var user models.Users
	if err := configs.DB.Where("id = ?", scans.UserId).First(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	owner := ScanOwner{
		Name:     user.Name,
		Username: user.Username,
		Avatar:   users.GetAvatarByID(user.Id),
	}

	response := ScanResponse{
		Id:         scans.Slug,
		Data:       scanData,
		Urls:       urls,
		Public:     scans.Public,
		Owner:      owner,
		HtmlPage:   utils.GetScanPage(scans.Slug),
		ApiPage:    utils.GetScanApiPage(c, scans.Slug),
		ReportPage: utils.GetScanApiReportPage(c, scans.Slug),
		CreatedAt:  scans.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}
