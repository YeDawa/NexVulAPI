package tasks

import (
	"bufio"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
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

func ExtractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	return host, nil
}
