package database

import (
	"fmt"
	"github.com/weikunlu/go-api-template/services/database"
)

func RunMigration(direction string, step string) error {
	fmt.Printf("ready to migrate %s to version %s\n", direction, step)

	err := database.NewDatabase()
	if err != nil {
		fmt.Printf("get database connection error: %s", err.Error())
	}
	defer database.GetDb().Close()

	database.PerformMigrations(direction, step)

	return nil
}
