package scans

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"nexvul/configs"
	"nexvul/generator"
	"nexvul/models"
	"nexvul/tasks"
	"nexvul/utils"

	"nexvul/controllers/users"

	"github.com/labstack/echo/v4"
)

type ScanRequest struct {
	URLs         []string `json:"urls"`
	Public       bool     `json:"public"`
	Domains      []string `json:"domains"`
	CORS         []string `json:"cors"`
	WordlistData string   `json:"wordlist_url"`
}

type Wordlist struct {
	Id  uint   `json:"id"`
	URL string `json:"url"`
}

func ScanHandler(c echo.Context) error {
	userID := users.GetUserID(c)
	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	var req ScanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if len(req.URLs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLs are required"})
	}

	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}

	var siteAnalyses []tasks.SiteAnalysis
	for _, url := range req.URLs {
		siteAnalyses = append(siteAnalyses, tasks.AnalyzeSingleURL(client, url))
	}

	domains := req.Domains
	if len(domains) == 0 {
		domainSet := make(map[string]struct{})

		for _, u := range req.URLs {
			d, err := tasks.ExtractDomain(u)
			if err == nil && d != "" {
				domainSet[d] = struct{}{}
			}
		}
		for d := range domainSet {
			domains = append(domains, d)
		}
	}

	var subdomainResults []tasks.SubdomainResult
	var WordlistData Wordlist

	if req.WordlistData != "" && len(domains) > 0 {
		var wordlist models.CustomWordlists
		if err := configs.DB.Where("slug = ?", req.WordlistData).First(&wordlist).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Failed to retrieve wordlist: " + err.Error(),
			})
		}

		WordlistData = Wordlist{
			Id: wordlist.Id,
		}

		wordlistSlice, err := tasks.FetchRemoteWordlist(wordlist.Url)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch wordlist: %v\n", err)
		} else {
			var wg sync.WaitGroup
			found := make(chan tasks.SubdomainResult, 1000)
			var mutex sync.Mutex

			var collectorWg sync.WaitGroup
			collectorWg.Add(1)
			go func() {
				defer collectorWg.Done()
				for res := range found {
					mutex.Lock()
					subdomainResults = append(subdomainResults, res)
					mutex.Unlock()
				}
			}()

			for _, domain := range domains {
				domain = strings.TrimSpace(domain)
				if domain == "" {
					continue
				}
				for _, sub := range wordlistSlice {
					sub = strings.TrimSpace(sub)
					if sub == "" {
						continue
					}
					wg.Add(1)
					go tasks.ScanSubdomain(&wg, sub, domain, found)
					time.Sleep(2 * time.Millisecond)
				}
			}

			wg.Wait()
			close(found)
			collectorWg.Wait()
		}
	}

	var robotsResults []tasks.RobotsData
	var corsResults []tasks.CORSScanResult

	for _, url := range req.URLs {
		report, err := tasks.ParseRobotsTxt(url)

		if err == nil && len(report.Directives) > 0 {
			robotsResults = append(robotsResults, report)
		}

		corsScan, err := tasks.ScanCORS(url)
		if err == nil {
			corsResults = append(corsResults, *corsScan)
		}
	}

	jsonCorsResults, err := json.Marshal(corsResults)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize CORS data",
		})
	}

	jsonRobotsData, err := json.Marshal(robotsResults)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize robots.txt data",
		})
	}

	executionTime := time.Since(startTime)
	jsonResultsData, err := json.Marshal(siteAnalyses)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize analysis data",
		})
	}

	jsonUrlsData, err := json.Marshal(req.URLs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize URLs",
		})
	}

	jsonSubdomainsData, err := json.Marshal(subdomainResults)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to serialize subdomains",
		})
	}

	if !req.Public {
		req.Public = true
	}

	slug := generator.String(8, 16)
	newScan := models.Scans{
		Slug:      slug,
		Public:    req.Public,
		Urls:      string(jsonUrlsData),
		CORS:      string(jsonCorsResults),
		CreatedAt: time.Now(),
		Data:      string(jsonResultsData),
		Robots:    string(jsonRobotsData),
		UserId:    uint(userIDUint),
	}

	if req.WordlistData != "" && len(subdomainResults) > 0 {
		newScan.Subdomains = string(jsonSubdomainsData)
		newScan.Wordlist = WordlistData.Id
	}

	result := configs.DB.Create(&newScan)
	if result.Error != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	response := map[string]interface{}{
		"success":        true,
		"id":             slug,
		"execution_time": executionTime,
		"html_page":      utils.GetScanPage(slug),
		"message":        "Scan created successfully",
	}

	return c.JSON(http.StatusCreated, response)
}
