package engine

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"httpshield/utils"
	"httpshield/configs"

	"github.com/go-pdf/fpdf"
	"github.com/labstack/echo/v4"
)

func GenerateMultiSitePDF(sites []SiteAnalysis) ([]byte, error) {
	if err := utils.DownloadFontIfNeeded(); err != nil {
		return nil, fmt.Errorf("font download error: %w", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("Roboto", "", configs.FontPath)

	pdf.SetFont("Roboto", "", 12)

	for _, site := range sites {
		pdf.AddPage()

		pdf.Cell(0, 10, utils.SanitizeText("Security Headers Report"))
		pdf.Ln(12)

		pdf.Cell(0, 8, utils.SanitizeText(fmt.Sprintf("Site: %s", site.URL)))
		pdf.Ln(8)

		pdf.Cell(0, 8, utils.SanitizeText(fmt.Sprintf("Security Score: %d%%", site.SecurityScore)))
		pdf.Ln(10)

		pdf.Cell(0, 8, utils.SanitizeText("Header Analysis"))
		pdf.Ln(8)

		for _, res := range site.Results {
			pdf.Cell(0, 6, utils.SanitizeText(fmt.Sprintf("%s: %s", res.Header, res.Status)))
			pdf.Ln(6)

			if res.Note != "" {
				pdf.SetFont("Roboto", "", 10)
				pdf.MultiCell(0, 5, utils.SanitizeText("→ "+res.Note), "", "", false)
				pdf.SetFont("Roboto", "", 12)
			}
		}

		if len(site.Recommendations) > 0 {
			pdf.Ln(8)
			pdf.Cell(0, 8, utils.SanitizeText("Recommendations"))
			pdf.Ln(8)

			for _, rec := range site.Recommendations {
				pdf.MultiCell(0, 5, utils.SanitizeText("• "+rec), "", "", false)
			}
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	return buf.Bytes(), nil
}

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

	pdfBytes, err := GenerateMultiSitePDF(siteAnalyses)
	if err != nil {
		fmt.Println("GenerateMultiSitePDF error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	filename := fmt.Sprintf("report_%d.pdf", time.Now().UnixNano())
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
