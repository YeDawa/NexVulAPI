package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, trying to use system environment variables")
	}
}

func InitDB() {
	LoadConfig()

	dsn := os.Getenv("DATABASE_PUBLIC_URL")
	if dsn == "" {
		log.Fatalf("DATABASE_URL environment variable not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}
