package scans

import (
	"net/http"
	"time"

	"httpshield/configs"
	"httpshield/generator"
	"httpshield/models"
	"httpshield/tasks"
	"httpshield/utils"

	// "httpshield/controllers/users"

	"github.com/labstack/echo/v4"
)

func AnalyzeHeaders(c echo.Context) error {
	// userID := users.GetUserID(c)

	// userIDUint, err := strconv.ParseUint(userID, 10, 64)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 		"success": false,
	// 		"error":   "Unathorized: Invalid user ID",
	// 	})
	// }

	var req tasks.RequestPayload
	if err := c.Bind(&req); err != nil || len(req.URLs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing URLs"})
	}

	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	var siteAnalyses []tasks.SiteAnalysis

	for _, targetURL := range req.URLs {
		siteAnalyses = append(siteAnalyses, tasks.AnalyzeSingleURL(client, targetURL))
	}

	executionTime := time.Since(startTime)

	jsonResultsData, err := utils.ToJSONString(siteAnalyses)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize analysis data",
		})
	}

	jsonUrlsData, err := utils.ToJSONString(req.URLs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize urls",
		})
	}

	// Replace generator.String with a local implementation to avoid import cycle
	slug := generator.String(8, 16)
	newScan := models.Scans{
		Slug:      slug,
		Urls:      jsonUrlsData,
		CreatedAt: time.Now(),
		Data:      jsonResultsData,
		// UserId:        uint(userIDUint),
	}

	result := configs.DB.Create(&newScan)
	if result.Error != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success":        true,
		"id":             slug,
		"data":           siteAnalyses,
		"execution_time": executionTime,
		"html_page":      utils.GetScanPage(slug),
		"message":        "Scan created successfully",
	})
}
