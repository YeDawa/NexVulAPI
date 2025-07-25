package tasks

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

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
