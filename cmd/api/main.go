package main

import (
	"database/sql"
	"gapstack-api/internal/routes"
	"gapstack-api/internal/seeder"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite" // Correct driver name
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3" // SQLite driver for database/sql

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sql.Open("sqlite3", "./transactions.db")
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()
	// Run migrations
	m, err := migrate.New(
		"file://internal/migrations",
		"sqlite://transactions.db",
	)
	if err != nil {
		log.Fatalf("Migration setup failed: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	// ✅ Run seeder before starting API
	if err := seeder.SeedTransactions(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db)

	log.Println("✅ Server running on port 8080")
	r.Run(":8080")
}
