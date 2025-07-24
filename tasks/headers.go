package tasks

import (
	"net/http"
	"strings"
	"time"

	"nexvul/configs"
	"nexvul/utils"
)

func AnalyzeSingleURL(client *http.Client, targetURL string) SiteAnalysis {
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
