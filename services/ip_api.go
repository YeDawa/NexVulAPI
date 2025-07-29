package services

import (
	"fmt"
	"net/http"

	"nexvul/utils"

	"github.com/labstack/echo/v4"
)

func GetIPInfo(c echo.Context) error {
	ip := c.QueryParam("ip")
	if ip == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "IP address is required"})
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch IP information"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "Failed to fetch IP information"})
	}

	var ipInfo map[string]interface{}
	if err := utils.DecodeJSON(resp.Body, &ipInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode IP information"})
	}

	return c.JSON(http.StatusOK, ipInfo)
}
