package scans

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"httpshield/configs"
	"httpshield/controllers/users"
	"httpshield/generator"
	"httpshield/models"
	"httpshield/tasks"
	"httpshield/utils"

	"github.com/labstack/echo/v4"
)

type ScanRequest struct {
	URLs        []string `json:"urls"`         // URLs para análise de headers e, se domains vazio, para scan de subdomínios
	Domains     []string `json:"domains"`      // Domínios para scan de subdomínios (opcional)
	WordlistURL string   `json:"wordlist_url"` // Wordlist remota para subdomínios
}

type SubdomainResult struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	SSL       bool   `json:"ssl"`
}

func fetchRemoteWordlist(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			list = append(list, line)
		}
	}

	return list, scanner.Err()
}

func tryRequest(url string, useTLS bool) (bool, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	if useTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true, nil
	}
	return false, nil
}

func scanSubdomain(wg *sync.WaitGroup, subdomain, domain string, found chan<- SubdomainResult) {
	defer wg.Done()

	full := fmt.Sprintf("%s.%s", subdomain, domain)

	_, err := net.LookupHost(full)
	if err != nil {
		return
	}

	httpsURL := "https://" + full
	if ok, _ := tryRequest(httpsURL, true); ok {
		found <- SubdomainResult{
			Domain:    domain,
			Subdomain: full,
			SSL:       true,
		}
		return
	}

	httpURL := "http://" + full
	if ok, _ := tryRequest(httpURL, false); ok {
		found <- SubdomainResult{
			Domain:    domain,
			Subdomain: full,
			SSL:       false,
		}
	}
}

func AnalyzeSingleURL(client *http.Client, targetURL string) tasks.SiteAnalysis {
	return tasks.AnalyzeSingleURL(client, targetURL)
}

func extractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	return host, nil
}

func ScanHandler(c echo.Context) error {
	userID := users.GetUserID(c)
	userIDUint, err := strconv.ParseUint(userID, 10, 64)

	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 		"success": false,
	// 		"error":   "Unathorized: Invalid user ID",
	// 	})
	// }

	var req ScanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if len(req.URLs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLs are required"})
	}

	if req.WordlistURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wordlist_url is required"})
	}

	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}

	var siteAnalyses []tasks.SiteAnalysis
	for _, url := range req.URLs {
		siteAnalyses = append(siteAnalyses, AnalyzeSingleURL(client, url))
	}

	domains := req.Domains
	if len(domains) == 0 {
		domainSet := make(map[string]struct{})
		for _, u := range req.URLs {
			d, err := extractDomain(u)
			if err == nil && d != "" {
				domainSet[d] = struct{}{}
			}
		}
		for d := range domainSet {
			domains = append(domains, d)
		}
	}

	var subdomainResults []SubdomainResult
	if len(domains) > 0 && req.WordlistURL != "" {
		wordlist, err := fetchRemoteWordlist(req.WordlistURL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch wordlist: " + err.Error()})
		}

		var wg sync.WaitGroup
		found := make(chan SubdomainResult, 1000)
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
			for _, sub := range wordlist {
				sub = strings.TrimSpace(sub)
				if sub == "" {
					continue
				}
				wg.Add(1)
				go scanSubdomain(&wg, sub, domain, found)
				time.Sleep(2 * time.Millisecond)
			}
		}

		wg.Wait()
		close(found)
		collectorWg.Wait()
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

	slug := generator.String(8, 16)
	newScan := models.Scans{
		Slug:       slug,
		Urls:       string(jsonUrlsData),
		CreatedAt:  time.Now(),
		Data:       string(jsonResultsData),
		Subdomains: string(jsonSubdomainsData),
		Wordlist:   req.WordlistURL,
		UserId:     uint(userIDUint),
	}

	result := configs.DB.Create(&newScan)
	if result.Error != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success":        true,
		"id":             slug,
		"results":        siteAnalyses,
		"execution_time": executionTime,
		"subdomains":     subdomainResults,
		"wordlist_url":   req.WordlistURL,
		"html_page":      utils.GetScanPage(slug),
		"message":        "Scan created successfully",
	})
}
