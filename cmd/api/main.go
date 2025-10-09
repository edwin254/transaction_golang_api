package main

import (
	"gapstack-api/internal/config"
	"gapstack-api/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db)

	log.Println("✅ Server running on port 8080")
	r.Run(":8080")
}
