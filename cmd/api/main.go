package main

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Correct driver name
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	// 1. Check for the .env file in the current directory (good for local testing from root)
	err := godotenv.Load("/app/app.env")

	// If it still fails, log the error but don't halt,
	// as environment variables might be set directly (e.g., in a Docker environment)
	if err != nil {
		log.Printf("Warning: .env file not found, loading system environment variables only. Error: %v", err)
	} else {
		fmt.Println(".env file loaded successfully.")
	}
}

func main() {

	var app App
	cfg, _ := LoadDBConfig()
	app.CreateConnection(cfg)
	app.RunMigrations(cfg)
	app.RunSeeder()
	app.CreateRoutes()
	app.Run()

}
