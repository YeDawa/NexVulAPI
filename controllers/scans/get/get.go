package get_scan

import (
	"encoding/json"
	"net/http"
	"time"

	"nexvul/configs"
	"nexvul/controllers/users"
	"nexvul/models"
	"nexvul/tasks"
	"nexvul/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HeaderResult struct {
	Header string `json:"header"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type ScanData struct {
	Ip              string         `json:"ip,omitempty"`
	URL             string         `json:"url"`
	Server          string         `json:"server"`
	Method          string         `json:"method"`
	ExecutionTime   int64          `json:"execution_time"`
	StatusCode      int            `json:"status_code"`
	ContentType     string         `json:"content_type"`
	Results         []HeaderResult `json:"results"`
	SecurityScore   int            `json:"security_score"`
	Recommendations []string       `json:"recommendations"`
}

type DomainGroup struct {
	Domain     string   `json:"domain,omitempty"`
	Subdomains []string `json:"subdomains,omitempty"`
}

type ScanResponse struct {
	Id         string                 `json:"id"`
	Data       []ScanData             `json:"data"`
	Urls       []string               `json:"urls"`
	CORS       []tasks.CORSScanResult `json:"cors,omitempty"`
	Subdomains []DomainGroup          `json:"subdomains,omitempty"`
	Wordlist   ScanWordlist           `json:"wordlist,omitempty"`
	Robots     []tasks.RobotsData     `json:"robots,omitempty"`
	HtmlPage   string                 `json:"html_page"`
	ReportPage string                 `json:"report_page"`
	ApiPage    string                 `json:"api_page"`
	Public     bool                   `json:"public"`
	Owner      ScanOwner              `json:"owner,omitempty"`
	CreatedAt  string                 `json:"created_at"`
}

type ScanOwner struct {
	Profile  string `json:"html_page,omitempty"`
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Name     string `json:"name,omitempty"`
}

type ScanWordlist struct {
	TotalLines int    `json:"total_lines,omitempty"`
	HtmlPage   string `json:"html_page,omitempty"`
	Name       string `json:"name,omitempty"`
}

type SubdomainInfo struct {
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	SSL       bool   `json:"ssl"`
}

func GetScanDetails(c echo.Context) error {
	id := c.Param("id")

	var scans models.Scans
	result := configs.DB.Where("slug = ?", id).First(&scans)

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

	var scanData []ScanData
	var urls []string
	var domainGroups []DomainGroup
	var robots []tasks.RobotsData
	var cors []tasks.CORSScanResult

	if scans.Subdomains != "" {
		var subdomainList []SubdomainInfo
		if err := json.Unmarshal([]byte(scans.Subdomains), &subdomainList); err == nil {
			domainMap := make(map[string][]string)

			for _, item := range subdomainList {
				domainMap[item.Domain] = append(domainMap[item.Domain], item.Subdomain)
			}

			for domain, subs := range domainMap {
				domainGroups = append(domainGroups, DomainGroup{
					Domain:     domain,
					Subdomains: subs,
				})
			}
		}
	}

	if err := json.Unmarshal([]byte(scans.Data), &scanData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to deserialize 'data' field: " + err.Error(),
		})
	}

	if err := json.Unmarshal([]byte(scans.Urls), &urls); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to deserialize 'urls' field: " + err.Error(),
		})
	}

	if err := json.Unmarshal([]byte(scans.Robots), &robots); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to deserialize 'robots' field: " + err.Error(),
		})
	}

	if err := json.Unmarshal([]byte(scans.CORS), &cors); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to deserialize 'cors' field: " + err.Error(),
		})
	}

	var user models.Users
	if scans.UserId > 0 {
		if err := configs.DB.Where("id = ?", scans.UserId).First(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
	}

	var wordlistData ScanWordlist
	if scans.Wordlist != 0 {
		var wordlist models.CustomWordlists
		if err := configs.DB.Where("id = ?", scans.Wordlist).First(&wordlist).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Failed to retrieve wordlist: " + err.Error(),
			})
		}

		wordlistData = ScanWordlist{
			Name:       wordlist.Name,
			TotalLines: wordlist.TotalLines,
			HtmlPage:   utils.GetWordlistPage(wordlist.Slug),
		}
	}

	var owner ScanOwner
	if scans.UserId != 0 {
		owner = ScanOwner{
			Name:     user.Name,
			Username: user.Username,
			Avatar:   users.GetAvatarByID(user.Id),
			Profile:  utils.GetOwnerProfilePage(user.Username),
		}
	}

	response := ScanResponse{
		Id:         scans.Slug,
		Data:       scanData,
		CORS:       cors,
		Subdomains: domainGroups,
		Robots:     robots,
		Urls:       urls,
		Public:     scans.Public,
		Wordlist:   wordlistData,
		HtmlPage:   utils.GetScanPage(scans.Slug),
		ApiPage:    utils.GetScanApiPage(c, scans.Slug),
		ReportPage: utils.GetScanApiReportPage(c, scans.Slug),
		CreatedAt:  scans.CreatedAt.Format(time.RFC3339),
	}

	if scans.UserId > 0 {
		response.Owner = owner
	}

	return c.JSON(http.StatusOK, response)
}
