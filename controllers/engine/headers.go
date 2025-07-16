package engine

import (
	"net/http"
	"strings"
	"time"

	"httpshield/configs"
	"httpshield/utils"

	"github.com/labstack/echo/v4"
)

type RequestPayload struct {
	URLs []string `json:"urls"`
}

type AnalysisResult struct {
	Header string `json:"header"`
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
}

type SiteAnalysis struct {
	URL             string           `json:"url"`
	Server          string           `json:"server,omitempty"`
	HttpMethod      string           `json:"method,omitempty"`
	ExecutionTime   time.Duration    `json:"execution_time"`
	StatusCode      int              `json:"status_code,omitempty"`
	ContentType     string           `json:"content_type,omitempty"`
	Results         []AnalysisResult `json:"results"`
	SecurityScore   int              `json:"security_score"`
	Recommendations []string         `json:"recommendations"`
}

type MultiSiteResponse struct {
	Success       bool           `json:"success"`
	Sites         []SiteAnalysis `json:"data"`
	ExecutionTime time.Duration  `json:"execution_time"`
}

func analyzeSingleURL(client *http.Client, targetURL string) SiteAnalysis {
	executionStart := time.Now()
	analysis := SiteAnalysis{URL: targetURL}

	resp, err := client.Head(targetURL)
	if err != nil || resp == nil {
		analysis.Results = []AnalysisResult{
			{Header: "N/A", Status: "Error", Note: "Failed to fetch URL headers"},
		}
		analysis.SecurityScore = 0
		return analysis
	}

	defer resp.Body.Close()
	normalizedHeaders := make(map[string]string)

	for k, v := range resp.Header {
		if len(v) > 0 {
			normalizedHeaders[strings.ToLower(k)] = v[0]
		}
	}

	var results []AnalysisResult
	var recommendations []string
	missing := 0

	for _, header := range configs.RequiredHeaders {
		lh := strings.ToLower(header)
		if val, ok := normalizedHeaders[lh]; ok && val != "" {
			results = append(results, AnalysisResult{
				Header: header,
				Status: "Present",
				Note:   utils.HeaderDescription(header),
			})
		} else {
			results = append(results, AnalysisResult{
				Header: header,
				Status: "Missing",
				Note:   utils.HeaderDescription(header),
			})

			missing++
			recommendations = append(recommendations, utils.GenerateRecommendation(header))
		}
	}

	total := len(configs.RequiredHeaders)
	score := int(float64(total-missing) / float64(total) * 100)
	if score < 0 {
		score = 0
	}

	analysis.Results = results
	analysis.StatusCode = resp.StatusCode
	analysis.Server = resp.Header.Get("Server")
	analysis.HttpMethod = resp.Request.Method
	analysis.ContentType = resp.Header.Get("Content-Type")
	analysis.Recommendations = recommendations
	analysis.SecurityScore = score
	analysis.ExecutionTime = time.Since(executionStart)
	return analysis
}

func AnalyzeHeaders(c echo.Context) error {
	var req RequestPayload
	if err := c.Bind(&req); err != nil || len(req.URLs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing URLs"})
	}

	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}
	var siteAnalyses []SiteAnalysis

	for _, targetURL := range req.URLs {
		siteAnalyses = append(siteAnalyses, analyzeSingleURL(client, targetURL))
	}

	executionTime := time.Since(startTime)
	return c.JSON(http.StatusOK, MultiSiteResponse{
		Success:       true,
		Sites:         siteAnalyses,
		ExecutionTime: executionTime,
	})
}
