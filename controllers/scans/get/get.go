package get_scan

import (
	"encoding/json"
	"net/http"
	"time"

	"httpshield/configs"
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
	Slug      string     `json:"slug"`
	Data      []ScanData `json:"data"`
	Urls      []string   `json:"urls"`
	HtmlPage  string     `json:"html_page"`
	ApiPage   string     `json:"api_page"`
	Public    bool       `json:"public"`
	CreatedAt string     `json:"created_at"`
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

	response := ScanResponse{
		Slug:      scans.Slug,
		Data:      scanData,
		Urls:      urls,
		Public:    scans.Public,
		HtmlPage:  utils.GetScanPage(scans.Slug),
		ApiPage:   utils.GetScanApiPage(c, scans.Slug),
		CreatedAt: scans.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}
