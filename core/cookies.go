package core

import (
	"log"
	"net/http"
	"os"
	"strings"

	"nexvul/configs"

	"github.com/labstack/echo/v4"
)

func SetCookie(c echo.Context, name, value string, hoursDuration int) {
	isDevelopment := strings.Contains(c.Request().Host, "localhost") || strings.Contains(c.Request().Host, "127.0.0.1")

	domain := ""
	if !isDevelopment {
		domain = configs.DomainName
	}

	cookie := &http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Name:     name,
		Value:    value,
		Domain:   domain,
		MaxAge:   hoursDuration * 3600,
		SameSite: http.SameSiteLaxMode,
	}

	if isDevelopment {
		cookie.SameSite = http.SameSiteLaxMode
		if origin := c.Request().Header.Get("Origin"); origin != "" {
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		}
	}

	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	c.SetCookie(cookie)

	if os.Getenv("DEBUG") == "true" {
		log.Printf("Setting cookie: name=%s, value=%s, secure=%t, httpOnly=%t, isDev=%t",
			name, value, cookie.Secure, cookie.HttpOnly, isDevelopment)
	}
}

func GetCookie(c echo.Context, name string) string {
	cookie, err := c.Cookie(name)
	if err != nil {
		if os.Getenv("DEBUG") == "true" {
			log.Printf("Cookie not found: name=%s, error=%v", name, err)
		}
		return ""
	}

	if os.Getenv("DEBUG") == "true" {
		log.Printf("Cookie found: name=%s, value=%s", name, cookie.Value)
	}

	return cookie.Value
}
