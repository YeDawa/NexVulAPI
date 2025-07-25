package wordlists

import (
	"bufio"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetWordlistPreviewContent(c echo.Context) error {
	url := c.QueryParam("url")
	maxLinesParam := c.QueryParam("max_lines")

	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "error": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "error": "Failed to fetch remote file"})
	}

	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().WriteHeader(http.StatusOK)
	scanner := bufio.NewScanner(resp.Body)

	lineCount := 0
	maxLines := 1000

	if maxLinesParam != "" && maxLinesParam != "0" {
		if n, err := strconv.Atoi(maxLinesParam); err == nil {
			maxLines = n
		}
	}
	
	for scanner.Scan() {
		if lineCount >= maxLines {
			break
		}
		line := scanner.Text() + "\n"
		c.Response().Write([]byte(line))
		lineCount++
	}

	return nil
}
