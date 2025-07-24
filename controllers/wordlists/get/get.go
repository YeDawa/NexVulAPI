package get_wordlist

import (
	"net/http"
	"time"

	"httpshield/configs"
	"httpshield/controllers/users"
	"httpshield/models"
	"httpshield/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type WordlistResponse struct {
	Id         string        `json:"id"`
	RawUrl     string        `json:"raw_url"`
	Name       string        `json:"name"`
	FileName   string        `json:"file_name"`
	TotalLines uint          `json:"total_lines"`
	Owner      WordlistOwner `json:"owner,omitempty"`
	CreatedAt  string        `json:"created_at"`
}

type WordlistOwner struct {
	Profile  string `json:"html_page,omitempty"`
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Name     string `json:"name,omitempty"`
}

func GetWordlistDetails(c echo.Context) error {
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

	var user models.Users
	if wordlist.UserId > 0 {
		if err := configs.DB.Where("id = ?", wordlist.UserId).First(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
	}

	var owner WordlistOwner
	if wordlist.UserId != 0 {
		owner = WordlistOwner{
			Name:     user.Name,
			Username: user.Username,
			Avatar:   users.GetAvatarByID(user.Id),
			Profile:  utils.GetOwnerProfilePage(user.Username),
		}
	}

	response := WordlistResponse{
		Id:         wordlist.Slug,
		RawUrl:     wordlist.Url,
		Name:       wordlist.Name,
		FileName:   wordlist.FileName,
		TotalLines: uint(wordlist.TotalLines),
		CreatedAt:  wordlist.CreatedAt.Format(time.RFC3339),
	}

	if wordlist.UserId > 0 {
		response.Owner = owner
	}

	return c.JSON(http.StatusOK, response)
}
