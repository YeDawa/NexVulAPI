package main

import (
	"fmt"
	"nexvul/configs"
	"nexvul/models"
)

func main() {
	configs.InitDB()

	// configs.DB.Migrator().DropTable(
	// 	&models.Scans{},
	// 	&models.CustomWordlists{},
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
