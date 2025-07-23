package subdomains

import (
	"bufio"
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

type ScanResponse struct {
	Found map[string][]string `json:"found"`
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

func scanSubdomain(wg *sync.WaitGroup, subdomain, domain string, found chan<- string) {
	defer wg.Done()
	full := fmt.Sprintf("%s.%s", subdomain, domain)
	hosts, err := net.LookupHost(full)
	if err == nil && len(hosts) > 0 {
		found <- full
	}
}

func ScanHandler(c echo.Context) error {
	var req ScanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid JSON"})
	}

	if len(req.Domains) == 0 || req.WordlistURL == "" {
		return c.JSON(400, map[string]string{"error": "domains and wordlist_url are required"})
	}

	wordlist, err := fetchRemoteWordlist(req.WordlistURL)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch wordlist: " + err.Error()})
	}

	output := make(map[string][]string)
	var globalWg sync.WaitGroup
	mutex := sync.Mutex{}

	for _, domain := range req.Domains {
		domain := strings.TrimSpace(domain)
		if domain == "" {
			continue
		}

		globalWg.Add(1)
		go func(domain string) {
			defer globalWg.Done()

			var wg sync.WaitGroup
			found := make(chan string, 100)
			var subs []string

			subsMutex := sync.Mutex{}
			go func() {
				for sub := range found {
					subsMutex.Lock()
					subs = append(subs, sub)
					subsMutex.Unlock()
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

			mutex.Lock()
			output[domain] = subs
			mutex.Unlock()
		}(domain)
	}

	globalWg.Wait()
	return c.JSON(200, output)
}
