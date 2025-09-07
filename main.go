package main

import (
	"os"

	"nexvul/configs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"nexvul/controllers/account"
	"nexvul/controllers/profile"
	"nexvul/controllers/scans"
	"nexvul/controllers/scans/get"
	"nexvul/controllers/tools"
	"nexvul/controllers/users"
	"nexvul/controllers/wordlists"
	"nexvul/controllers/wordlists/get"
)

func main() {
	e := echo.New()
	configs.InitDB()

	// Enable CORS with credentials support
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.HEAD, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Cookie", "Set-Cookie"},
		ExposeHeaders:    []string{"Set-Cookie", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Security Middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
	}))

	// Add middleware to set common headers
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			return next(c)
		}
	})

	// API's Users
	e.POST("/login", users.Login)
	e.POST("/register", users.CreateUser)
	e.DELETE("/users/logoff", users.Logoff)
	e.GET("/users/check/:user", users.CheckUsername)

	// API's User Autenticated Content
	e.GET("/users/me", account.UserLogged)
	e.GET("/users/me/scans", account.ListUserScans)
	e.GET("/users/me/wordlists", account.ListUserWordlists)

	// API's Scans
	e.POST("/scan", scans.ScanHandler)
	e.GET("/scans", scans.ListPublicScans)
	e.GET("/scan/:id", get_scan.GetScanDetails)
	e.GET("/scan/:id/export", get_scan.GenerateReport)

	// API's Wordlists
	e.POST("/wordlist", wordlists.ImportWordlist)
	e.GET("/wordlists/:id", get_wordlist.GetWordlistDetails)
	e.GET("/wordlists/:id/raw", get_wordlist.GetWordlistRawContent)
	e.GET("/wordlists/preview", wordlists.GetWordlistPreviewContent)

	// API's Profile
	e.GET("/user/:user", profile.ProfilePublic)
	e.GET("/user/:user/scans", profile.ListPublicScans)
	e.GET("/user/:user/wordlists", profile.ListPublicWordlists)

	// External API's
	e.GET("/ext/ip/:ip", tools.GetIPInfo)

	e.Start(":" + os.Getenv("PORT"))
}
