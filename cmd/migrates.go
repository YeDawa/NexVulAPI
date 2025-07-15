package main

import (
	"fmt"
	"httpshield/configs"
	"httpshield/models"
)

func main() {
	configs.InitDB()

	configs.DB.AutoMigrate(
		&models.User{},
		&models.Scans{},
		&models.Exports{},
	)

	fmt.Println("Migrate complete!")
}
