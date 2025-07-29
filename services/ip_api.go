package services

import (
	"fmt"
	"net/http"

	"nexvul/utils"
)

func GetIPInfo(ip string) string {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Error fetching IP info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Error fetching IP info: %v", err)
	}

	var ipInfo map[string]interface{}
	if err := utils.DecodeJSON(resp.Body, &ipInfo); err != nil {
		return fmt.Sprintf("Error decoding IP info: %v", err)
	}

	jsonBytes, err := utils.EncodeJSON(ipInfo)
	if err != nil {
		return fmt.Sprintf("Error encoding IP info: %v", err)
	}
	return string(jsonBytes)
}
