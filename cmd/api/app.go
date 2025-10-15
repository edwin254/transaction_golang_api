package app

import (
	"database/sql"
	"fmt"
	config "gapstack-api/internal/config"
	"gapstack-api/internal/routes"
	"gapstack-api/internal/seeder"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Correct driver name
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

type App struct {
	DB     *sql.DB
	Routes *gin.Engine
}

// LoadDBConfig loads environment variables from a .env file and the system,
// then populates and returns a DBConfig struct.
func LoadDBConfig() (config.DBConfig, error) {
	// 1. Load .env file. gotdotenv.Load() handles if the file is missing or loads it successfully.
	// If you need specific file paths, use godotenv.Load(".env.development", ".env.test", ...)
	if err := godotenv.Load(); err != nil {

		// Log a warning if the file isn't found,
		log.Println("Warning: Could not find .env file.")
	}

	var cfg config.DBConfig

	// 2. Load required string variables and check for existence
	cfg.POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	cfg.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	cfg.POSTGRES_DB = os.Getenv("POSTGRES_DB")

	// Check for required string variables
	if cfg.POSTGRES_HOST == "" || cfg.POSTGRES_USER == "" || cfg.POSTGRES_DB == "" {
		return config.DBConfig{}, fmt.Errorf("required environment variables (HOST, USER, DB) are missing")
	}

	log.Printf("Configuration loaded successfully. Host: %s, DB: %s",
		cfg.POSTGRES_HOST, cfg.POSTGRES_DB)

	return cfg, nil
}

func (a *App) CreateConnection(dbCfg config.DBConfig) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", &dbCfg.POSTGRES_USER, &dbCfg.POSTGRES_PASSWORD, &dbCfg.POSTGRES_HOST, &dbCfg.POSTGRES_DB)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	// Test the connection immediately
	if err = db.Ping(); err != nil {
		db.Close() // Close connection on failure
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Store the successful connection in the App struct
	a.DB = db
	log.Println("Database connection established successfully.")
	return nil
}

func (a *App) CreateRoutes() {
	r := gin.Default()
	routes.RegisterRoutes(r, a.DB)
}

func (a *App) Run() {
	a.Routes.Run(":8080")

	log.Println("Server running on port 8080")
}

func (a *App) RunMigrations() {

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

}

func (a *App) RunSeeder() {
	// Run seeder before starting API
	if err := seeder.SeedTransactions(a.DB); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("✅ Seeding DB Done")
}
