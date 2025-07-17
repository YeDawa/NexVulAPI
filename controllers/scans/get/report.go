package get_scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"httpshield/configs"
	"httpshield/generator"
	"httpshield/models"
	"httpshield/tasks"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GenerateReport(c echo.Context) error {
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

	client := &http.Client{Timeout: 10 * time.Second}
	var analyses []tasks.SiteAnalysis
	for _, targetURL := range urls {
		analysis := tasks.AnalyzeSingleURL(client, targetURL)
		analyses = append(analyses, analysis)
	}

	pdfBytes, err := generator.GeneratePDF(analyses)
	if err != nil {
		fmt.Println("GenerateMultiSitePDF error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	filename := fmt.Sprintf("report_%d.pdf", time.Now().UnixNano())
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
