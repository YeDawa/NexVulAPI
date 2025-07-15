package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"httpshield/controllers/engine"
)

func main() {
	e := echo.New()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, trying to use system environment variables")
	}

	e.POST("/headers", engine.AnalyzeHeaders)
	e.POST("/export", engine.ExportPDF)

	e.Start(":" + os.Getenv("PORT"))
}
