package tasks

import "time"

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

type ScanHeadersResponseReport struct {
	URL             string      `json:"url"`
	Server          string      `json:"server"`
	HttpMethod      string      `json:"http_method"`
	ExecutionTime   float64     `json:"execution_time"`
	StatusCode      int         `json:"status_code"`
	ContentType     string      `json:"content_type"`
	Results         interface{} `json:"results"`
	SecurityScore   float64     `json:"security_score"`
	Recommendations []string    `json:"recommendations"`
}

func ScanHeadersDataResponse(analysis SiteAnalysis) ScanHeadersResponseReport {
	return ScanHeadersResponseReport{
		URL:             analysis.URL,
		Server:          analysis.Server,
		HttpMethod:      analysis.HttpMethod,
		ExecutionTime:   analysis.ExecutionTime.Seconds(),
		StatusCode:      analysis.StatusCode,
		ContentType:     analysis.ContentType,
		Results:         analysis.Results,
		SecurityScore:   float64(analysis.SecurityScore),
		Recommendations: analysis.Recommendations,
	}
}