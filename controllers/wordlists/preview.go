package wordlists

import (
	"net/http"
	"strings"

	"nexvul/utils"

	"github.com/labstack/echo/v4"
)

func GetWordlistPreviewContent(c echo.Context) error {
	url := c.QueryParam("url")

	var builder strings.Builder
	err := utils.StreamRemoteFile(url, func(line string) {
		builder.WriteString(line)
		builder.WriteString("\n")
	})

	content := builder.String()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to stream file: " + err.Error(),
		})
	}

	return c.String(http.StatusOK, content)
}
