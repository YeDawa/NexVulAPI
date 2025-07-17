package scans

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func ExportPDF(c echo.Context) error {
	var req RequestPayload
	if err := c.Bind(&req); err != nil || len(req.URLs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing URLs"})
	}

	client := &http.Client{Timeout: 10 * time.Second}
	var siteAnalyses []SiteAnalysis
	for _, targetURL := range req.URLs {
		siteAnalyses = append(siteAnalyses, analyzeSingleURL(client, targetURL))
	}

	pdfBytes, err := GeneratePDF(siteAnalyses)
	if err != nil {
		fmt.Println("GenerateMultiSitePDF error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	filename := fmt.Sprintf("report_%d.pdf", time.Now().UnixNano())
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
