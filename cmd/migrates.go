package main

import (
	"fmt"
	"httpshield/configs"
	"httpshield/models"
)

func main() {
	configs.InitDB()

	configs.DB.AutoMigrate(
		&models.Users{},
		&models.Scans{},
		&models.Exports{},
		&models.CustomHeaders{},
	)

	fmt.Println("Migrate complete!")
}
