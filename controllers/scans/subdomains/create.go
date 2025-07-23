package subdomains

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type ScanRequest struct {
	Domains     []string `json:"domains"`
	WordlistURL string   `json:"wordlist_url"`
}

type SubdomainResult struct {
	Domain     string `json:"domain"`
	Subdomain  string `json:"subdomain"`
	StatusCode int    `json:"status_code"`
	SSL        bool   `json:"ssl"`
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

func tryRequest(url string, useTLS bool) (int, bool) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	if useTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // ignora erro de certificado
			},
		}
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, false
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, false
	}
	defer resp.Body.Close()
	return resp.StatusCode, true
}

func scanSubdomain(wg *sync.WaitGroup, subdomain, domain string, found chan<- SubdomainResult) {
	defer wg.Done()

	full := fmt.Sprintf("%s.%s", subdomain, domain)

	// DNS resolve
	_, err := net.LookupHost(full)
	if err != nil {
		return
	}

	// Tenta HTTPS
	httpsURL := "https://" + full
	if code, ok := tryRequest(httpsURL, true); ok {
		found <- SubdomainResult{
			Domain:     domain,
			Subdomain:  full,
			StatusCode: code,
			SSL:        true,
		}
		return
	}

	// Fallback para HTTP
	httpURL := "http://" + full
	if code, ok := tryRequest(httpURL, false); ok {
		found <- SubdomainResult{
			Domain:     domain,
			Subdomain:  full,
			StatusCode: code,
			SSL:        false,
		}
	}
}

func ScanHandler(c echo.Context) error {
	var req ScanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	if len(req.Domains) == 0 || req.WordlistURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "domains and wordlist_url are required"})
	}

	wordlist, err := fetchRemoteWordlist(req.WordlistURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch wordlist: " + err.Error()})
	}

	var globalWg sync.WaitGroup
	var mutex sync.Mutex
	var results []SubdomainResult

	for _, domain := range req.Domains {
		domain := strings.TrimSpace(domain)
		if domain == "" {
			continue
		}

		globalWg.Add(1)
		go func(domain string) {
			defer globalWg.Done()

			var wg sync.WaitGroup
			found := make(chan SubdomainResult, 100)

			collectorWg := sync.WaitGroup{}
			collectorWg.Add(1)
			go func() {
				defer collectorWg.Done()
				for result := range found {
					mutex.Lock()
					results = append(results, result)
					mutex.Unlock()
				}
			}()

			for _, word := range wordlist {
				word = strings.TrimSpace(word)
				if word == "" {
					continue
				}
				wg.Add(1)
				go scanSubdomain(&wg, word, domain, found)
				time.Sleep(3 * time.Millisecond)
			}

			wg.Wait()
			close(found)
			collectorWg.Wait()
		}(domain)
	}

	globalWg.Wait()
	return c.JSON(http.StatusOK, results)
}
