package main

import (
	"fmt"
	"httpshield/configs"
	"httpshield/models"
)

func main() {
	configs.InitDB()

	configs.DB.Migrator().DropTable(&models.Users{})

	configs.DB.AutoMigrate(
		&models.Users{},
		&models.Scans{},
		&models.Reports{},
		&models.ScansXSS{},
		&models.Payloads{},
		&models.CustomHeaders{},
	)

	fmt.Println("Migrate complete!")
}
