package main

import (
	"fmt"
	"httpshield/configs"
	"httpshield/models"
)

func main() {
	configs.InitDB()

	// configs.DB.Migrator().DropTable(&models.Users{})

	configs.DB.AutoMigrate(
		&models.Users{},
		&models.Scans{},
		&models.Reports{},
		&models.Profile{},
		&models.ScansXSS{},
		&models.Payloads{},
		&models.CustomHeaders{},
		&models.CustomWordlists{},
		&models.DefaultWordLists{},
	)

	fmt.Println("Migrate complete!")
}
