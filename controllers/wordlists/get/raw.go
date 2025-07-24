package get_wordlist

import (
	"net/http"
	
	"httpshield/configs"
	"httpshield/models"
	"httpshield/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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

	var content string
	err := utils.StreamRemoteFile(wordlist.Url, func(line string) {
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
