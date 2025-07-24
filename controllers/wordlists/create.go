package wordlists

import (
	"fmt"
	"time"

	"nexvul/configs"
	"nexvul/generator"
	"nexvul/models"
	"nexvul/utils"
)

type CreateWordlistRequest struct {
	WordlistName string `json:"wordlist_name"`
	WordlistURL  string `json:"wordlist_url"`
}

type CreateWordlistResponse struct {
	Id       uint                   `json:"id"`
	Slug     string                 `json:"slug"`
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Wordlist models.CustomWordlists `json:"wordlist"`
}

func CreateWordlist(wordlistURL string, UserId uint) (CreateWordlistResponse, error) {
	slug := generator.String(8, 16)
	lastPart := utils.GetLastPartOfURL(wordlistURL)

	totalLines, err := utils.CountRemoteFileLines(wordlistURL)
	if err != nil {
		return CreateWordlistResponse{}, fmt.Errorf("failed to count lines in wordlist: %v", err)
	}

	customWordlist := models.CustomWordlists{
		Slug:       slug,
		Name:       lastPart,
		Url:        wordlistURL,
		UserId:     UserId,
		FileName:   lastPart,
		CreatedAt:  time.Now(),
		TotalLines: totalLines,
	}

	if err := configs.DB.Create(&customWordlist).Error; err != nil {
		return CreateWordlistResponse{}, fmt.Errorf("failed to create wordlist: %v", err)
	}

	response := CreateWordlistResponse{
		Id:       customWordlist.Id,
		Slug:     slug,
		Success:  true,
		Message:  "Wordlist created successfully",
		Wordlist: customWordlist,
	}

	return response, nil
}
