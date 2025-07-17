package main

import (
	"log"
	"os"

	"httpshield/configs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"httpshield/controllers/scans"
	"httpshield/controllers/scans/get"
)

func main() {
	e := echo.New()
	configs.InitDB()

	// Enable CORS with credentials support
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.HEAD, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Cookie", "Set-Cookie"},
		ExposeHeaders:    []string{"Set-Cookie", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Add middleware to set common headers
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			return next(c)
		}
	})

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, trying to use system environment variables")
	}

	e.POST("/scan", scans.AnalyzeHeaders)
	e.GET("/scan/:id", get_scan.GetScanDetails)
	e.POST("/export", scans.ExportPDF)

	e.Start(":" + os.Getenv("PORT"))
}
