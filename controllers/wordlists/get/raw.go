package get_wordlist

import (
	"bufio"
	"fmt"
	"net/http"
	
	"httpshield/configs"
	"httpshield/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func StreamRemoteFile(url string, lineHandler func(string)) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to perform GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lineHandler(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	return nil
}

func GetWordlistRawContent(c echo.Context) error {
	id := c.Param("id")
	var wordlist models.CustomWordlists
	result := configs.DB.Where("slug = ?", id).First(&wordlist)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Item not found",
			})
		}
		
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   result.Error.Error(),
		})
	}

	err := StreamRemoteFile(wordlist.Url, func(line string) {
		fmt.Println(line)
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to stream file: " + err.Error(),
		})
	}

	var content string
	err = StreamRemoteFile(wordlist.Url, func(line string) {
		content += line + "\n"
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to stream file: " + err.Error(),
		})
	}

	return c.String(http.StatusOK, content)
}
