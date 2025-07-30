package tasks

import (
	"fmt"
	"net/http"
	"time"

	"nexvul/configs"
)

func ScanCORS(target string) (*CORSScanResult, error) {
	origin := configs.HTMLPageURI

	req, err := http.NewRequest("OPTIONS", target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", "GET")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode == 404 {
		req, _ = http.NewRequest("GET", target, nil)
		req.Header.Set("Origin", origin)
		resp, err = client.Do(req)

		if err != nil {
			return &CORSScanResult{Error: err.Error()}, nil
		}
	}
	defer resp.Body.Close()

	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	result := &CORSScanResult{
		URL:              target,
		Origin:           origin,
		Status:           resp.StatusCode,
		AllowOrigin:      allowOrigin,
		AllowMethods:     resp.Header.Get("Access-Control-Allow-Methods"),
		AllowHeaders:     resp.Header.Get("Access-Control-Allow-Headers"),
		AllowCredentials: resp.Header.Get("Access-Control-Allow-Credentials"),
		Permissive:       allowOrigin == "*",
		Reflected:        allowOrigin == origin,
	}

	for k, v := range resp.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

	return result, nil
}

