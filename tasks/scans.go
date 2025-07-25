package tasks

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func FetchRemoteWordlist(url string) ([]string, error) {
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

func TryRequest(url string, useTLS bool) (bool, error) {
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

	if resp.StatusCode >= 200 {
		return true, nil
	}
	return false, nil
}

func ScanSubdomain(wg *sync.WaitGroup, subdomain, domain string, found chan<- SubdomainResult) {
	defer wg.Done()

	full := fmt.Sprintf("%s.%s", subdomain, domain)

	_, err := net.LookupHost(full)
	if err != nil {
		return
	}

	httpsURL := "https://" + full
	if _, err := TryRequest(httpsURL, true); err == nil {
		found <- SubdomainResult{
			Domain:    domain,
			Subdomain: full,
			SSL:       true,
		}
		return
	}

	httpURL := "http://" + full
	if _, err := TryRequest(httpURL, false); err == nil {
		found <- SubdomainResult{
			Domain:    domain,
			Subdomain: full,
			SSL:       false,
		}
	}
}

func AnalyzeURL(client *http.Client, targetURL string) SiteAnalysis {
	return AnalyzeSingleURL(client, targetURL)
}

func ExtractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	return host, nil
}
