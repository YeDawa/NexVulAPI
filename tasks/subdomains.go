package tasks

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

func resolveSubdomains(domain string, wordlist []string) []string {
	var found []string
	var wg sync.WaitGroup
	foundChan := make(chan string, len(wordlist))

	for _, sub := range wordlist {
		sub := strings.TrimSpace(sub)
		if sub == "" {
			continue
		}

		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			full := sub + "." + domain
			_, err := net.LookupHost(full)
			if err == nil {
				foundChan <- full
			}
		}(sub)

		time.Sleep(5 * time.Millisecond)
	}

	wg.Wait()
	close(foundChan)

	for item := range foundChan {
		found = append(found, item)
	}

	return found
}

func AnalyzeSubdomainsFromURL(domain string, wordlistURL string) ([]SiteAnalysis, error) {
	wordlist, err := fetchRemoteWordlist(wordlistURL)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	subdomains := resolveSubdomains(domain, wordlist)

	var results []SiteAnalysis
	for _, sub := range subdomains {
		url := "http://" + sub
		analysis := AnalyzeSingleURL(client, url)
		results = append(results, analysis)
	}

	return results, nil
}
