package get_scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nexvul/configs"
	"nexvul/generator"
	"nexvul/models"
	"nexvul/tasks"
	"nexvul/utils"

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
	var subdomainInfo []tasks.SubdomainInfo

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

	if scans.Subdomains != "" {
		if err := json.Unmarshal([]byte(scans.Subdomains), &subdomainInfo); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   "Failed to deserialize 'subdomains' field: " + err.Error(),
			})
		}

		domainMap := make(map[string][]tasks.SubdomainInfo)
		for _, info := range subdomainInfo {
			domainMap[info.Domain] = append(domainMap[info.Domain], info)
		}

		subdomainInfo = make([]tasks.SubdomainInfo, 0)
		for _, infos := range domainMap {
			subdomainInfo = append(subdomainInfo, infos...)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	var analyses []tasks.SiteAnalysis
	for _, targetURL := range urls {
		analysis := tasks.AnalyzeSingleURL(client, targetURL)
		analyses = append(analyses, analysis)
	}

	var wordlistData ScanWordlist
	if scans.Wordlist != 0 {
		var wordlist models.CustomWordlists
		if err := configs.DB.Where("id = ?", scans.Wordlist).First(&wordlist).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Failed to retrieve wordlist: " + err.Error(),
			})
		}

		wordlistData = ScanWordlist{
			HtmlPage:   utils.GetWordlistPage(wordlist.Slug),
		}
	}

	pdfBytes, err := generator.GeneratePDF(analyses, wordlistData.HtmlPage, subdomainInfo)
	if err != nil {
		fmt.Println("GenerateMultiSitePDF error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	filename := fmt.Sprintf("report_%d.pdf", time.Now().UnixNano())
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
