package utils

import (
	"fmt"

	"httpshield/configs"

	"github.com/labstack/echo/v4"
)

func GetScanPage(id string) string {
	return fmt.Sprintf("%sscan/%s", configs.HTMLPageURI, id)
}

func GetScanApiPage(c echo.Context, id string) string {
	return fmt.Sprintf("%s/scan/%s", configs.GetRootURL(c), id)
}