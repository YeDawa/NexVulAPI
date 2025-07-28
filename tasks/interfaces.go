package tasks

import "time"

type RequestPayload struct {
	URLs        []string `json:"urls"`
	WordlistURL string   `json:"wordlist_url"`
}

type AnalysisResult struct {
	Header string `json:"header"`
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
}

type SubdomainInfo struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	SSL       bool   `json:"ssl"`
}

type SiteAnalysis struct {
	URL             string           `json:"url"`
	Server          string           `json:"server,omitempty"`
	HttpMethod      string           `json:"method,omitempty"`
	ExecutionTime   time.Duration    `json:"execution_time"`
	StatusCode      int              `json:"status_code,omitempty"`
	ContentType     string           `json:"content_type,omitempty"`
	Results         []AnalysisResult `json:"results"`
	Subdomains      []string         `json:"subdomains"`
	SecurityScore   int              `json:"security_score"`
	Recommendations []string         `json:"recommendations"`
}

type MultiSiteResponse struct {
	Success       bool           `json:"success"`
	Sites         []SiteAnalysis `json:"data"`
	ExecutionTime time.Duration  `json:"execution_time"`
}

type SubdomainResult struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	SSL       bool   `json:"ssl"`
}

type RobotsDirective struct {
	UserAgent string   `json:"user_agent"`
	Allow     []string `json:"allow"`
	Disallow  []string `json:"disallow"`
}

type RobotsData struct {
	Target     string            `json:"target"`
	Sitemaps   []string          `json:"sitemaps"`
	Directives []RobotsDirective `json:"directives"`
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
