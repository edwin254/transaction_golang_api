package main

import (
	app "gapstack-api/cmd/api"
	"gapstack-api/internal/routes"
	"gapstack-api/internal/seeder"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Correct driver name
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gin-gonic/gin"
)

func main() {
	var app app.App
	// Run migrations
	m, err := migrate.New(
		"file://internal/migrations",
		"postgres://postgres:postgres@localhost:5432/example?sslmode=disable",
	)
	if err != nil {
		log.Fatalf("Migration setup failed: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	// Run seeder before starting API
	if err := seeder.SeedTransactions(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db)

	log.Println("✅ Server running on port 8080")
	r.Run(":8080")
}
