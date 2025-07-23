package main

import (
	"fmt"
	"httpshield/configs"
	"httpshield/models"
)

func main() {
	configs.InitDB()

	// configs.DB.Migrator().DropTable(
	// 	&models.CustomWordlists{},
	// 	&models.DefaultWordLists{},
	// 	&models.CustomHeaders{},
	// )

	configs.DB.AutoMigrate(
		&models.Users{},
		&models.Scans{},
		&models.Reports{},
		&models.Profile{},
		&models.ScansXSS{},
		&models.Payloads{},
		&models.CustomWordlists{},
	)

	fmt.Println("Migrate complete!")
}
