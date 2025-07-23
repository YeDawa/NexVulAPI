package scans

import (
	"fmt"
	"net/http"
	// "strconv"
	"time"

	// "httpshield/configs"
	// "httpshield/generator"
	// "httpshield/models"
	"httpshield/tasks"
	// "httpshield/utils"

	// "httpshield/controllers/users"

	"github.com/labstack/echo/v4"
)

func AnalyzeHeaders(c echo.Context) error {
	// userID := users.GetUserID(c)
	// userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	var req tasks.RequestPayload
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	var siteAnalyses []tasks.SiteAnalysis

	// Caso URLs diretas sejam passadas
	if len(req.URLs) > 0 {
		for _, targetURL := range req.URLs {
			// Analisa a URL principal
			siteAnalyses = append(siteAnalyses, tasks.AnalyzeSingleURL(client, targetURL))

			// Se WordlistURL estiver presente, analisa subdomínios também
			if req.WordlistURL != "" {
				subResults, err := tasks.AnalyzeSubdomainsFromURL(targetURL, req.WordlistURL)
				fmt.Printf("Analyzing subdomains for %s with wordlist %s\n", targetURL, req.WordlistURL)
				if err == nil {
					siteAnalyses = append(siteAnalyses, subResults...)
					fmt.Printf("Subdomains found for %s: %v\n", targetURL, subResults)
				}
			}
		}
	} else if len(req.URLs) > 0 && req.WordlistURL != "" {
		subResults, err := tasks.AnalyzeSubdomainsFromURL(req.URLs[0], req.WordlistURL)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to fetch or process subdomains: " + err.Error(),
			})
		}

		siteAnalyses = subResults
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Either URLs or (domain and wordlist_url) must be provided",
		})
	}

	executionTime := time.Since(startTime)

	// jsonResultsData, err := utils.ToJSONString(siteAnalyses)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 		"success": false,
	// 		"error":   "Failed to serialize analysis data",
	// 	})
	// }

	// jsonUrlsData, err := utils.ToJSONString(req.URLs)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 		"success": false,
	// 		"error":   "Failed to serialize URLs",
	// 	})
	// }

	// slug := generator.String(8, 16)
	// newScan := models.Scans{
	// 	Slug:      slug,
	// 	Urls:      jsonUrlsData,
	// 	CreatedAt: time.Now(),
	// 	Data:      jsonResultsData,
	// 	// UserId:    uint(userIDUint),
	// }

	// result := configs.DB.Create(&newScan)
	// if result.Error != nil {
	// 	return c.JSON(http.StatusConflict, map[string]interface{}{
	// 		"success": false,
	// 		"error":   result.Error.Error(),
	// 	})
	// }

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		// "id":             slug,
		"data":           siteAnalyses,
		"execution_time": executionTime,
		"wordlist_url":   req.WordlistURL,
		// "html_page":      utils.GetScanPage(slug),
		"message": "Scan created successfully",
	})
}
