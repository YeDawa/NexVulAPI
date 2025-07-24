package utils

import (
	"fmt"

	"httpshield/configs"

	"github.com/labstack/echo/v4"
)

func GetScanPage(id string) string {
	return fmt.Sprintf("%s/scan/%s", configs.HTMLPageURI, id)
}

func GetWordlistPage(id string) string {
	return fmt.Sprintf("%s/wordlist/%s", configs.HTMLPageURI, id)
}

func GetOwnerProfilePage(user string) string {
	return fmt.Sprintf("%s/user/%s", configs.HTMLPageURI, user)
}

func GetScanApiPage(c echo.Context, id string) string {
	return fmt.Sprintf("%s/scan/%s", configs.GetRootURL(c), id)
}

func GetScanApiReportPage(c echo.Context, id string) string {
	return fmt.Sprintf("%s/scan/%s/export", configs.GetRootURL(c), id)
}

func GetWordlistRawPage(c echo.Context, id string) string {
	return fmt.Sprintf("%s/wordlist/%s/raw", configs.GetRootURL(c), id)
}